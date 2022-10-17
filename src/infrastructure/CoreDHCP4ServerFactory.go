package infrastructure

import (
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
)

//CoreDHCP4ServerFactory fabric for creating dhcp v4 servers
type CoreDHCP4ServerFactory struct {
	leasesRepo interfaces.IGenericRepository[domain.DHCP4Lease]
}

//NewCoreDHCP4ServerFactory constructor for CoreDHCP v4 servers manager
func NewCoreDHCP4ServerFactory(
	leasesRepo interfaces.IGenericRepository[domain.DHCP4Lease],
) interfaces.IDHCP4ServerFactory {
	return &CoreDHCP4ServerFactory{
		leasesRepo: leasesRepo,
	}
}

//Create new coreDHCP v4 server with config
//
//Params:
//	config - dhcp v4 config
//Return:
//	error - if an error occurred, otherwise nil
func (m *CoreDHCP4ServerFactory) Create(config domain.DHCP4Config) (interfaces.IDHCP4Server, error) {
	server, err := NewCoreDHCP4Server(config, m.leasesRepo)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to create dhcp v4 server")
	}
	return server, nil
}
