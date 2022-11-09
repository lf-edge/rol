package infrastructure

import (
	"context"
	"github.com/google/uuid"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
)

//EthernetSwitchManagerProvider struct for switch manager getter
type EthernetSwitchManagerProvider struct {
	switchRepo interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitch]
	managers   map[uuid.UUID]interfaces.IEthernetSwitchManager
}

//NewEthernetSwitchManagerProvider constructor for EthernetSwitchManagerProvider
func NewEthernetSwitchManagerProvider(switchRepo interfaces.IGenericRepository[uuid.UUID, domain.EthernetSwitch]) interfaces.IEthernetSwitchManagerProvider {
	return &EthernetSwitchManagerProvider{
		managers:   make(map[uuid.UUID]interfaces.IEthernetSwitchManager),
		switchRepo: switchRepo,
	}
}

//Get ethernet switch manager
//
//Params:
//	ethernetSwitch - switch entity
//Return:
//	interfaces.IEthernetSwitchManager - switch manager interface
func (e *EthernetSwitchManagerProvider) Get(ctx context.Context, switchID uuid.UUID) (interfaces.IEthernetSwitchManager, error) {
	if e.managers[switchID] == nil {
		ethSwitch, err := e.switchRepo.GetByID(ctx, switchID)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "failed to get ethernet switch configuration from repository")
		}
		switch ethSwitch.SwitchModel {
		case "tl-sg2210mp":
			e.managers[switchID] = NewTPLinkEthernetSwitchManager(ethSwitch.Address+":23", ethSwitch.Username, ethSwitch.Password)
			return e.managers[switchID], nil
		}
		return nil, nil
	}
	return e.managers[switchID], nil
}
