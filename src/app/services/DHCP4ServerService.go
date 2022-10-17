package services

import (
	"context"
	"github.com/google/uuid"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"
)

//DHCP4ServerService service structure for managing DHCP servers
type DHCP4ServerService struct {
	configsRepo interfaces.IGenericRepository[domain.DHCP4Config]
	leasesRepo  interfaces.IGenericRepository[domain.DHCP4Lease]
	factory     interfaces.IDHCP4ServerFactory
	servers     map[uuid.UUID]interfaces.IDHCP4Server
}

//NewDHCP4ServerService constructor for DHCPServerService service
//
//Params:
//	servers - repository with domain.DHCPServer entity
//	leases - repository with domain.DHCPLease entity
//Return:
//	*DHCPServerService - New DHCP servers service
func NewDHCP4ServerService(
	configs interfaces.IGenericRepository[domain.DHCP4Config],
	leases interfaces.IGenericRepository[domain.DHCP4Lease],
	dhcp4factory interfaces.IDHCP4ServerFactory,
) *DHCP4ServerService {
	return &DHCP4ServerService{
		configsRepo: configs,
		leasesRepo:  leases,
		servers:     map[uuid.UUID]interfaces.IDHCP4Server{},
		factory:     dhcp4factory,
	}
}

func (s *DHCP4ServerService) getServerState(configID uuid.UUID, enabled bool) domain.DHCPServerState {
	if server, ok := s.servers[configID]; ok {
		return server.GetState()
	} else if enabled {
		return domain.DHCPStateError
	} else {
		return domain.DHCPStateStopped
	}
}

func (s *DHCP4ServerService) serverIsExist(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := s.configsRepo.GetByID(ctx, id)
	if err != nil {
		if errors.As(err, errors.NotFound) {
			return false, nil
		}
		return false, errors.Internal.Wrap(err, "repository failed to get ethernet switch")
	}
	return true, nil
}

func (s *DHCP4ServerService) serverExistenceCheck(ctx context.Context, serverID uuid.UUID) error {
	serverExist, err := s.serverIsExist(ctx, serverID)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to check existence of the server")
	}
	if !serverExist {
		return errors.NotFound.New("server with this ID is not found")
	}
	return nil
}

//GetServerList Get list of DHCP servers with search and pagination
//
//Params:
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	dtos.PaginatedItemsDto[dtos.DHCP4ServerDto] - paginated list of DHCP v4 servers
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) GetServerList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.DHCP4ServerDto], error) {
	paginatedItems, err := GetList[dtos.DHCP4ServerDto](ctx, s.configsRepo, search, orderBy, orderDirection, page, pageSize)
	if err != nil {
		return paginatedItems, err
	}
	for i, dtoItem := range paginatedItems.Items {
		(&paginatedItems.Items[i]).State = s.getServerState(dtoItem.ID, dtoItem.Enabled).String()
	}
	return paginatedItems, nil
}

//DHCP4ServerServiceInit starts all enabled DHCP servers
func DHCP4ServerServiceInit(s *DHCP4ServerService) error {
	ctx := context.Background()
	queryBuilder := s.configsRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("Enabled", "==", true)
	enabledServersCount, err := s.configsRepo.Count(ctx, queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "Failed to count enabled dhcp v4 servers")
	}
	serversConfigs, err := s.configsRepo.GetList(ctx, "", "", 1, int(enabledServersCount), queryBuilder)
	if err != nil {
		return errors.Internal.Wrap(err, "Failed to get configs for enabled dhcp v4 servers")
	}
	for _, config := range serversConfigs {
		server, err := s.factory.Create(config)
		if err != nil {
			return errors.Internal.Wrapf(err, "failed to create dhcp v4 server with id: %s", config.ID.String())
		}
		s.servers[config.ID] = server
		err = server.Start()
		if err != nil {
			return errors.Internal.Wrapf(err, "failed to start dhcp v4 server with id: %s", config.ID.String())
		}
	}
	return nil
}

//GetServerByID Get DHCP v4 server by ID
//Params
//	ctx - context is used only for logging
//	id - DHCP v4 server ID
//Return
//	dtos.DHCP4ServerDto - DHCP v4 server dto
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) GetServerByID(ctx context.Context, id uuid.UUID) (dtos.DHCP4ServerDto, error) {
	dto, err := GetByID[dtos.DHCP4ServerDto](ctx, s.configsRepo, id, nil)
	if err != nil {
		return dto, err
	}
	dto.State = s.getServerState(dto.ID, dto.Enabled).String()
	return dto, nil
}

//CreateServer create DHCP v4 server
//Params
//	ctx - context is used only for logging
//	createDto - dto for creating DHCP v4 server
//Return
//	dtos.DHCP4ServerDto - DHCP v4 server dto
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) CreateServer(ctx context.Context, createDto dtos.DHCP4ServerCreateDto) (dtos.DHCP4ServerDto, error) {
	dto := dtos.DHCP4ServerDto{}
	//check create DTO
	err := validators.ValidateDHCP4ServerCreateDto(createDto)
	if err != nil {
		return dto, err
	}

	//Save dhcp v4 configuration
	dto, err = Create[dtos.DHCP4ServerDto](ctx, s.configsRepo, createDto)
	if err != nil {
		return dto, err
	}
	config, err := s.configsRepo.GetByID(ctx, dto.ID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "failed to get config for dhcp v4 server")
	}

	//create runtime server
	server, err := s.factory.Create(config)
	if err != nil {
		return dto, errors.Wrap(err, "failed to create dhcp v4 server")
	}
	s.servers[config.ID] = server

	// Start runtime server and set state to dto
	dto.State = domain.DHCPStateStopped.String()
	if createDto.Enabled {
		err = server.Start()
		if err != nil {
			//TODO: logging error
			dto.State = domain.DHCPStateError.String()
		} else {
			dto.State = server.GetState().String()
		}
	}
	return dto, nil
}

//UpdateServer create DHCP v4 server
//Params
//	ctx - context is used only for logging
//	createDto - dto for creating DHCP v4 server
//Return
//	dtos.DHCP4ServerDto - DHCP v4 server dto
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) UpdateServer(ctx context.Context, id uuid.UUID, updateDto dtos.DHCP4ServerUpdateDto) (dtos.DHCP4ServerDto, error) {
	dto := dtos.DHCP4ServerDto{}
	//check update DTO
	err := validators.ValidateDHCP4ServerUpdateDto(updateDto)
	if err != nil {
		return dto, err
	}

	//Update configuration
	dto, err = Update[dtos.DHCP4ServerDto](ctx, s.configsRepo, updateDto, id, nil)
	if err != nil {
		return dto, err
	}
	config, err := s.configsRepo.GetByID(ctx, dto.ID)
	if err != nil {
		return dto, errors.Internal.Wrap(err, "failed to get config for dhcp v4 server")
	}

	// Reload configuration in runtime server
	var server interfaces.IDHCP4Server
	if server, ok := s.servers[config.ID]; !ok {
		//create new runtime server
		server, err = s.factory.Create(config)
		if err != nil {
			return dto, errors.Wrap(err, "failed to create dhcp v4 server")
		}
		s.servers[config.ID] = server
	} else {
		//Update runtime configuration on existed DHCP v4 server
		server.Stop()
		err = server.ReloadConfiguration(config)
		if err != nil {
			return dto, errors.Wrap(err, "failed to reload DHCP v4 server configuration")
		}
	}
	if config.Enabled {
		err = server.Start()
		if err != nil {
			//TODO: logging error
			dto.State = domain.DHCPStateError.String()
		} else {
			dto.State = server.GetState().String()
		}
	} else {
		dto.State = domain.DHCPStateStopped.String()
	}
	return dto, nil
}

//DeleteServer delete DHCP v4 server from server pool
//Params
//	ctx - context is used only for logging
//	id - ID for DHCP v4 server
//Return
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) DeleteServer(ctx context.Context, id uuid.UUID) error {
	// Stop and delete runtime server
	if server, ok := s.servers[id]; ok {
		server.Stop()
		delete(s.servers, id)
	}

	//Delete all leases
	queryBuilder := s.leasesRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("DHCP4ConfigID", "==", id)
	leasesCount, err := s.leasesRepo.Count(ctx, queryBuilder)
	if err != nil {
		return errors.Wrap(err, "failed to count all leases to remove")
	}
	leases, err := s.leasesRepo.GetList(ctx, "", "", 1, int(leasesCount), queryBuilder)
	if err != nil {
		return errors.Wrap(err, "failed to get all leases to remove")
	}
	for _, lease := range leases {
		err = s.leasesRepo.Delete(ctx, lease.ID)
		if err != nil {
			return errors.Wrap(err, "failed to remove leases")
		}
	}

	// Delete config
	err = s.configsRepo.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to remove dhcp server configuration")
	}
	return nil
}
