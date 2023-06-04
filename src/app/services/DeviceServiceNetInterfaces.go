package services

import (
	"context"
	"github.com/google/uuid"
	"rol/app/interfaces"
	"rol/dtos"
)

func (s *DeviceService) getDeviceIDQueryBuilder(ctx context.Context, deviceID uuid.UUID) interfaces.IQueryBuilder {
	queryBuilder := s.devEthRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("DeviceID", "==", deviceID)
	return queryBuilder
}

func (s *DeviceService) GetNetInterfacesList(ctx context.Context, deviceID uuid.UUID, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.DeviceNetworkInterfaceDto], error) {
	err := s.CheckDeviceExistence(ctx, deviceID)
	if err != nil {
		return dtos.PaginatedItemsDto[dtos.DeviceNetworkInterfaceDto]{}, err
	}
	queryBuilder := s.getDeviceIDQueryBuilder(ctx, deviceID)
	AddSearchInAllFields(search, s.devEthRepo, queryBuilder)
	netInterfaces, err := GetListExtended[dtos.DeviceNetworkInterfaceDto](ctx, s.devEthRepo, queryBuilder, orderBy,
		orderDirection, page, pageSize)
	if err != nil {
		return netInterfaces, err
	}
	return netInterfaces, nil
}

func (s *DeviceService) GetNetInterfaceByID(ctx context.Context, deviceID uuid.UUID, devNetID uuid.UUID) (
	dtos.DeviceNetworkInterfaceDto, error,
) {
	queryBuilder := s.getDeviceIDQueryBuilder(ctx, deviceID)
	netInterface, err := GetByID[dtos.DeviceNetworkInterfaceDto](ctx, s.devEthRepo, devNetID, queryBuilder)
	if err != nil {
		return netInterface, err
	}
	return netInterface, nil
}

func (s *DeviceService) CreateNetInterface(ctx context.Context, deviceID uuid.UUID,
	createDto dtos.DeviceNetworkInterfaceCreateDto) (dtos.DeviceNetworkInterfaceDto, error) {
	err := s.CheckDeviceExistence(ctx, deviceID)
	if err != nil {
		return dtos.DeviceNetworkInterfaceDto{}, err
	}
	netInterface, err := Create[dtos.DeviceNetworkInterfaceDto](ctx, s.devEthRepo, createDto)
	if err != nil {
		return netInterface, err
	}
	return netInterface, nil
}

func (s *DeviceService) UpdateNetInterface(ctx context.Context, deviceID uuid.UUID, devNetID uuid.UUID,
	updateDto dtos.DeviceNetworkInterfaceUpdateDto) (dtos.DeviceNetworkInterfaceDto, error) {
	queryBuilder := s.getDeviceIDQueryBuilder(ctx, deviceID)
	netInterface, err := Update[dtos.DeviceNetworkInterfaceDto](ctx, s.devEthRepo, updateDto, devNetID, queryBuilder)
	if err != nil {
		return netInterface, err
	}
	return netInterface, nil
}

func (s *DeviceService) DeleteNetInterface(ctx context.Context, deviceID uuid.UUID, devNetID uuid.UUID) error {
	err := s.CheckDeviceExistence(ctx, deviceID)
	if err != nil {
		return err
	}
	return s.devEthRepo.Delete(ctx, devNetID)
}
