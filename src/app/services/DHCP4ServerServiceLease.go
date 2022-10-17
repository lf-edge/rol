package services

import (
	"context"
	"github.com/google/uuid"
	"rol/app/errors"
	"rol/app/mappers"
	"rol/app/validators"
	"rol/domain"
	"rol/dtos"
)

//GetLeaseList Get list of DHCP server leases with search and pagination
//
//Params:
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	dtos.PaginatedItemsDto[dtos.DHCP4LeaseDto] - paginated list of DHCP v4 servers
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) GetLeaseList(ctx context.Context, serverID uuid.UUID, search, orderBy, orderDirection string, page, pageSize int) (
	dtos.PaginatedItemsDto[dtos.DHCP4LeaseDto],
	error,
) {
	paginatedItemsDto := dtos.NewEmptyPaginatedItemsDto[dtos.DHCP4LeaseDto]()
	err := s.serverExistenceCheck(ctx, serverID)
	if err != nil {
		return paginatedItemsDto, err
	}
	queryBuilder := s.leasesRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("DHCP4ConfigID", "==", serverID)
	if len(search) > 3 {
		AddSearchInAllFields(search, s.leasesRepo, queryBuilder)
	}
	return GetListExtended[dtos.DHCP4LeaseDto](ctx, s.leasesRepo, queryBuilder, orderBy, orderDirection, page, pageSize)
}

//GetLeaseByID Get DHCP v4 server lease by ID
//Params
//	ctx - context is used only for logging
//	serverID - DHCP v4 server ID
//	leaseID - DHCP v4 lease ID
//Return
//	dtos.DHCP4LeaseDto - DHCP v4 server dto
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) GetLeaseByID(ctx context.Context, serverID, leaseID uuid.UUID) (
	dtos.DHCP4LeaseDto,
	error,
) {
	dto := dtos.DHCP4LeaseDto{}
	err := s.serverExistenceCheck(ctx, serverID)
	if err != nil {
		return dto, err
	}
	queryBuilder := s.leasesRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("DHCP4ConfigID", "==", serverID)
	return GetByID[dtos.DHCP4LeaseDto](ctx, s.leasesRepo, leaseID, queryBuilder)
}

//CreateLease create DHCP v4 server lease
//Params
//	ctx - context is used only for logging
//	serverID - DHCP v4 server ID
//	createDto - dto for creating DHCP v4 lease
//Return
//	dtos.DHCP4LeaseDto - DHCP v4 lease dto
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) CreateLease(ctx context.Context, serverID uuid.UUID, createDto dtos.DHCP4LeaseCreateDto) (
	dtos.DHCP4LeaseDto,
	error,
) {
	outDto := dtos.DHCP4LeaseDto{}
	err := validators.ValidateDHCP4LeaseCreateDto(createDto)
	if err != nil {
		return outDto, err
	}
	err = s.serverExistenceCheck(ctx, serverID)
	if err != nil {
		return outDto, err
	}

	//Create entity
	entity := new(domain.DHCP4Lease)
	err = mappers.MapDtoToEntity(createDto, entity)
	if err != nil {
		return outDto, errors.Internal.Wrap(err, "error map entity to dto")
	}
	//Set config id
	entity.DHCP4ConfigID = serverID
	newEntity, err := s.leasesRepo.Insert(ctx, *entity)
	if err != nil {
		return outDto, errors.Internal.Wrap(err, "create entity error")
	}

	err = mappers.MapEntityToDto(newEntity, &outDto)
	if err != nil {
		return outDto, errors.Internal.Wrap(err, "error map dto to entity")
	}
	return outDto, nil
}

//UpdateLease create DHCP v4 server lease
//Params
//	ctx - context is used only for logging
//	serverID - DHCP v4 server ID
//	leaseID - DHCP v4 lease ID
//	createDto - dto for creating DHCP v4 lease
//Return
//	dtos.DHCP4ServerDto - DHCP v4 server dto
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) UpdateLease(ctx context.Context, serverID, leaseID uuid.UUID, updateDto dtos.DHCP4LeaseUpdateDto) (dtos.DHCP4LeaseDto, error) {
	dto := dtos.DHCP4LeaseDto{}
	err := validators.ValidateDHCP4LeaseUpdateDto(updateDto)
	if err != nil {
		return dto, nil
	}
	err = s.serverExistenceCheck(ctx, serverID)
	if err != nil {
		return dto, err
	}
	queryBuilder := s.leasesRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("DHCP4ConfigID", "==", serverID)
	return Update[dtos.DHCP4LeaseDto](ctx, s.leasesRepo, updateDto, leaseID, queryBuilder)
}

//DeleteLease delete DHCP v4 server lease
//Params
//	ctx - context is used only for logging
//	serverID - ID for DHCP v4 server
//	leaseID - ID for DHCP v4 lease
//Return
//	error - if an error occurs, otherwise nil
func (s *DHCP4ServerService) DeleteLease(ctx context.Context, serverID, leaseID uuid.UUID) error {
	err := s.serverExistenceCheck(ctx, serverID)
	if err != nil {
		return err
	}
	queryBuilder := s.leasesRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("DHCP4ConfigID", "==", serverID)

	lease, err := s.leasesRepo.GetByIDExtended(ctx, leaseID, queryBuilder)
	if err != nil {
		return errors.Wrap(err, "can't found lease")
	}
	err = s.leasesRepo.Delete(ctx, lease.ID)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to remove lease by id")
	}
	return nil
}
