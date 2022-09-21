package interfaces

import (
	"context"
	"github.com/google/uuid"
)

//IEthernetSwitchManagerProvider is the interface is used to get ethernet switch manager
type IEthernetSwitchManagerProvider interface {
	//Get ethernet switch manager
	Get(ctx context.Context, switchID uuid.UUID) (IEthernetSwitchManager, error)
}
