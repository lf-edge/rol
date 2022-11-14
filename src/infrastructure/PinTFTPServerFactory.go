// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
)

//PinTFTPServerFactory is implementation for ITFTPServerFactory interface
type PinTFTPServerFactory struct{}

//NewPinTFTPServerFactory creates new pin/tftp server factory
func NewPinTFTPServerFactory() (interfaces.ITFTPServerFactory, error) {
	return &PinTFTPServerFactory{}, nil
}

//Create pin tftp server
func (f *PinTFTPServerFactory) Create(config domain.TFTPConfig) (interfaces.ITFTPServer, error) {
	server, err := NewPinTFTPServer(config)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to create new pin/TFTP server")
	}
	return server, nil
}
