// Package services stores business logic for each entity
package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/domain"
	"rol/dtos"
	"strings"
)

//ProjectService is an object for create all project dependencies such as traffic rule and bridge
type ProjectService struct {
	repository         interfaces.IGenericRepository[uuid.UUID, domain.Project]
	config             *domain.AppConfig
	dhcpService        *DHCP4ServerService
	tftpService        *TFTPServerService
	hostNetworkService *HostNetworkService
	logger             *logrus.Logger
}

//NewProjectService ProjectService constructor
func NewProjectService(repo interfaces.IGenericRepository[uuid.UUID, domain.Project], DHCPService *DHCP4ServerService,
	TFTPService *TFTPServerService, networkService *HostNetworkService, config *domain.AppConfig, log *logrus.Logger) *ProjectService {
	return &ProjectService{
		repository:         repo,
		config:             config,
		dhcpService:        DHCPService,
		tftpService:        TFTPService,
		hostNetworkService: networkService,
		logger:             log,
	}
}

func (s ProjectService) createBridgeDto(createDto dtos.ProjectCreateDto) dtos.HostNetworkBridgeCreateDto {
	bridgeIP := createDto.Subnet[:len(createDto.Subnet)-1] + "1/24"
	return dtos.HostNetworkBridgeCreateDto{
		Name: createDto.Name,
		HostNetworkBridgeBaseDto: dtos.HostNetworkBridgeBaseDto{
			Addresses: []string{bridgeIP},
			Slaves:    []string{s.config.Network.Interface},
		},
	}
}

func (s ProjectService) createDHCPDto(bridgeIP, bridgeName string) dtos.DHCP4ServerCreateDto {
	ipLastDot := strings.LastIndex(bridgeIP, ".")
	ipPart := bridgeIP[:ipLastDot]
	return dtos.DHCP4ServerCreateDto{
		Range:     ipPart + ".2-" + ipPart + ".254",
		Mask:      "255.255.255.0",
		ServerID:  bridgeIP,
		Interface: bridgeName,
		Gateway:   bridgeIP,
		DNS:       "8.8.8.8;" + bridgeIP,
		NTP:       bridgeIP,
		Enabled:   true,
		Port:      67,
		LeaseTime: 3600,
	}
}

func (s ProjectService) createTFTPDto(bridgeIP string) dtos.TFTPServerCreateDto {
	return dtos.TFTPServerCreateDto{TFTPServerBaseDto: dtos.TFTPServerBaseDto{
		Address: bridgeIP,
		Port:    "69",
		Enabled: true,
	}}
}

func (s ProjectService) createTrafficRuleDto(createDto dtos.ProjectCreateDto) dtos.HostNetworkTrafficRuleCreateDto {
	return dtos.HostNetworkTrafficRuleCreateDto{HostNetworkTrafficRuleBaseDto: dtos.HostNetworkTrafficRuleBaseDto{
		Chain:       "POSTROUTING",
		Action:      "MASQUERADE",
		Source:      createDto.Subnet + "/24",
		Destination: "",
	}}
}

//GetByID gets project by id
//
//Params:
//	ctx - context is used only for logging
//	id - project ID
//Return:
//	dtos.ProjectDto - created project dto
//	error - if an error occurs, otherwise nil
func (s ProjectService) GetByID(ctx context.Context, id uuid.UUID) (dtos.ProjectDto, error) {
	return GetByID[dtos.ProjectDto](ctx, s.repository, id, nil)
}

//GetList gets list of projects with search and pagination
//
//Params:
//	ctx - context is used only for logging
//	search - string for search in entity string fields
//	orderBy - order by entity field name
//	orderDirection - ascending or descending order
//	page - page number
//	pageSize - page size
//Return
//	dtos.PaginatedItemsDto[dtos.ProjectDto] - paginated list of projects
//	error - if an error occurs, otherwise nil
func (s ProjectService) GetList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.ProjectDto], error) {
	return GetList[dtos.ProjectDto](ctx, s.repository, search, orderBy, orderDirection, page, pageSize)
}

//Create creates project and bridge with postrouting traffic rule for it
//
//Params:
//	project - project entity
//	projectSubnet - subnet for project
//Return:
//	domain.Project - created project
//	error - if an error occurs, otherwise nil
func (s ProjectService) Create(ctx context.Context, createDto dtos.ProjectCreateDto) (dtos.ProjectDto, error) {
	bridgeDto := s.createBridgeDto(createDto)
	bridge, err := s.hostNetworkService.CreateBridge(bridgeDto)
	if err != nil {
		return dtos.ProjectDto{}, errors.Internal.Wrap(err, "network service failed to create bridge")
	}
	bridgeIP := bridge.Addresses[0][:len(bridge.Addresses[0])-3]
	dhcpDto := s.createDHCPDto(bridgeIP, bridge.Name)
	dhcp, err := s.dhcpService.CreateServer(ctx, dhcpDto)
	if err != nil {
		return dtos.ProjectDto{}, errors.Internal.Wrap(err, "dhcp service failed to create server")
	}
	tftpDto := s.createTFTPDto(bridgeIP)
	tftp, err := s.tftpService.CreateServer(ctx, tftpDto)
	if err != nil {
		return dtos.ProjectDto{}, errors.Internal.Wrap(err, "tftp service failed to create server")
	}
	trafficRuleDto := s.createTrafficRuleDto(createDto)
	_, err = s.hostNetworkService.CreateTrafficRule("nat", trafficRuleDto)
	if err != nil {
		return dtos.ProjectDto{}, errors.Internal.Wrap(err, "network service failed to create traffic rule")
	}
	projectEntity := domain.Project{
		Name:         createDto.Name,
		BridgeName:   bridge.Name,
		Subnet:       createDto.Subnet,
		DHCPServerID: dhcp.ID,
		TFTPServerID: tftp.ID,
	}
	newProject, err := s.repository.Insert(ctx, projectEntity)
	if err != nil {
		return dtos.ProjectDto{}, errors.Internal.Wrap(err, "repository failed to create project")
	}
	projectOutDto := &dtos.ProjectDto{}
	err = mappers.MapEntityToDto(newProject, projectOutDto)
	if err != nil {
		return dtos.ProjectDto{}, errors.Internal.Wrap(err, "error map entity to dto")
	}
	err = s.hostNetworkService.Ping()
	if err != nil {
		return dtos.ProjectDto{}, errors.Internal.Wrap(err, "network service failed to save current configuration")
	}
	return *projectOutDto, nil
}

//Delete project and its bridge with postrouting traffic rule
//
//Params:
//	project - project entity
//Return:
//	error - if an error occurs, otherwise nil
func (s ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	project, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = s.dhcpService.DeleteServer(ctx, project.DHCPServerID)
	if err != nil && !errors.As(err, errors.NotFound) {
		return errors.Internal.Wrap(err, "dhcp service failed to delete server")
	}
	err = s.tftpService.DeleteServer(ctx, project.TFTPServerID)
	if err != nil && !errors.As(err, errors.NotFound) {
		return errors.Internal.Wrap(err, "tftp service failed to delete server")
	}
	trafficRule := dtos.HostNetworkTrafficRuleDeleteDto{HostNetworkTrafficRuleBaseDto: dtos.HostNetworkTrafficRuleBaseDto{
		Chain:       "POSTROUTING",
		Action:      "MASQUERADE",
		Source:      project.Subnet + "/24",
		Destination: "",
	}}
	err = s.hostNetworkService.DeleteTrafficRule("nat", trafficRule)
	if err != nil && !errors.As(err, errors.NotFound) {
		return errors.Internal.Wrap(err, "network service failed to delete traffic rule")
	}
	err = s.hostNetworkService.DeleteBridge(project.BridgeName)
	if err != nil && !errors.As(err, errors.NotFound) {
		return errors.Internal.Wrap(err, "network service failed to delete bridge")
	}
	err = s.hostNetworkService.Ping()
	if err != nil {
		return errors.Internal.Wrap(err, "network service failed to save current configuration")
	}

	return s.repository.Delete(ctx, id)
}
