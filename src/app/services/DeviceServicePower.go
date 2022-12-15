package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/errors"
	"rol/domain"
	"rol/dtos"
)

func (s *DeviceService) DevicePowerCommand(ctx context.Context, deviceID uuid.UUID, dto dtos.DevicePowerCommandDto) error {
	err := s.CheckDeviceExistence(ctx, deviceID)
	if err != nil {
		return err
	}
	dev, err := s.devRepo.GetByID(ctx, deviceID)
	if err != nil {
		return errors.Internal.Wrapf(err, "failed to get device with id %s", deviceID.String())
	}
	powerManagerDriver, err := s.powerDriverProvider.GetPowerManagerDriver(dev)
	if err != nil {
		return errors.Internal.Wrapf(err, "get power manager error for device %s", deviceID.String())
	}
	switch dto.Command {
	case "On":
		err = powerManagerDriver.PowerOn(ctx, dev)
	case "Off":
		err = powerManagerDriver.PowerOff(ctx, dev)
	default:
		return errors.Internal.Newf("unknown power command for device %s", deviceID.String())
	}
	return nil
}

func (s *DeviceService) getLastPowerState(ctx context.Context, deviceID uuid.UUID) (domain.DevPowerState, error) {
	err := s.CheckDeviceExistence(ctx, deviceID)
	if err != nil {
		return domain.DevInPowerUnknownState, err
	}
	queryBuilder := s.getDeviceIDQueryBuilder(ctx, deviceID)
	states, err := s.powerStateHistory.GetList(ctx, "CreatedAt", "desc", 1, 1, queryBuilder)
	if err != nil {
		return domain.DevInPowerUnknownState, errors.Internal.Wrapf(err, "get power state history failed for device %s", deviceID.String())
	}
	if len(states) < 1 {
		return domain.DevInPowerUnknownState, nil
	}
	return states[0].DevPowerState, nil
}

func (s *DeviceService) getLastPowerStateString(ctx context.Context, deviceID uuid.UUID) string {
	state, err := s.getLastPowerState(ctx, deviceID)
	if err != nil {
		logrus.Errorf("get last power state error for device %s: %s", deviceID, err)
	}
	return state.String()
}
