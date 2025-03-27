// Package services stores business logic for each entity
package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/providers"
	"rol/domain"
	"rol/dtos"
)

//DeviceService service structure for manage device entity and sub-entities
type DeviceService struct {
	devRepo             interfaces.IGenericRepository[uuid.UUID, domain.Device]
	devEthRepo          interfaces.IGenericRepository[uuid.UUID, domain.DeviceNetworkInterface]
	powerStateHistory   interfaces.IGenericRepository[uuid.UUID, domain.DevicePowerState]
	powerDriverProvider *providers.DevicePowerDriverProvider
	log                 *logrus.Logger
}

func NewDeviceService(
	devRepo interfaces.IGenericRepository[uuid.UUID, domain.Device],
	devEthRepo interfaces.IGenericRepository[uuid.UUID, domain.DeviceNetworkInterface],
	powerStateHistory interfaces.IGenericRepository[uuid.UUID, domain.DevicePowerState],
	powerDriverProvider *providers.DevicePowerDriverProvider,
	log *logrus.Logger,
) *DeviceService {
	return &DeviceService{
		devRepo:             devRepo,
		devEthRepo:          devEthRepo,
		powerStateHistory:   powerStateHistory,
		powerDriverProvider: powerDriverProvider,
		log:                 log,
	}
}

func (s *DeviceService) CheckDeviceExistence(ctx context.Context, deviceID uuid.UUID) error {
	exist, err := s.devRepo.IsExist(ctx, deviceID, nil)
	if err != nil {
		return errors.Internal.Wrapf(err, "failed to check existence of the device with id: %s", deviceID)
	}
	if !exist {
		return errors.NotFound.Newf("device with id %s not found", deviceID)
	}
	return nil
}

func (s *DeviceService) GetDevicesList(ctx context.Context, search, orderBy, orderDirection string, page, pageSize int) (dtos.PaginatedItemsDto[dtos.DeviceDto], error) {
	devices, err := GetList[dtos.DeviceDto](ctx, s.devRepo, search, orderBy, orderDirection, page, pageSize)
	if err != nil {
		return devices, errors.Internal.Wrap(err, "failed to get tftp server configs list")
	}
	for i, device := range devices.Items {
		(&devices.Items[i]).PowerState = s.getLastPowerStateString(ctx, device.ID)
	}
	return devices, nil
}

func (s *DeviceService) GetDeviceByID(ctx context.Context, id uuid.UUID) (dtos.DeviceDto, error) {
	dto, err := GetByID[dtos.DeviceDto](ctx, s.devRepo, id, nil)
	if err != nil {
		return dto, err
	}
	state := s.getLastPowerStateString(ctx, id)
	dto.PowerState = state
	return dto, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, id uuid.UUID, updateDto dtos.DeviceUpdateDto) (dtos.DeviceDto, error) {
	dto, err := Update[dtos.DeviceDto](ctx, s.devRepo, updateDto, id, nil)
	if err != nil {
		return dto, err
	}
	dto.PowerState = s.getLastPowerStateString(ctx, id)
	return dto, nil
}

func (s *DeviceService) CreateDevice(ctx context.Context, createDto dtos.DeviceCreateDto) (dtos.DeviceDto, error) {
	dto, err := Create[dtos.DeviceDto](ctx, s.devRepo, createDto)
	if err != nil {
		return dto, err
	}
	dto.PowerState = s.getLastPowerStateString(ctx, dto.ID)
	return dto, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	err := s.CheckDeviceExistence(ctx, id)
	if err != nil {
		return err
	}
	queryBuilder := s.getDeviceIDQueryBuilder(ctx, id)
	err = s.devEthRepo.DeleteAll(ctx, queryBuilder)
	if err != nil {
		return errors.Internal.Wrapf(err, "delete network interfaces failed for device %s", id.String())
	}
	err = s.powerStateHistory.DeleteAll(ctx, queryBuilder)
	if err != nil {
		return errors.Internal.Wrapf(err, "delete power history failed for device %s", id.String())
	}
	return s.devRepo.Delete(ctx, id)
}
