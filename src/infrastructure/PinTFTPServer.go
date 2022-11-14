// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"fmt"
	"github.com/pin/tftp/v3"
	"io"
	"os"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
)

//PinTFTPServer TFTP server implementation for ITFTPServer interface
type PinTFTPServer struct {
	runtime *tftp.Server
	config  domain.TFTPConfig
	paths   *[]domain.TFTPPathRatio
	state   domain.TFTPServerState
}

//NewPinTFTPServer creates new pin tftp server
func NewPinTFTPServer(config domain.TFTPConfig) (interfaces.ITFTPServer, error) {
	server := &PinTFTPServer{
		runtime: nil,
		config:  config,
		paths:   &[]domain.TFTPPathRatio{},
		state:   domain.TFTPStateStopped,
	}
	server.runtime = tftp.NewServer(
		func(filename string, rf io.ReaderFrom) error {
			actualPath := ""
			for _, ratio := range *server.paths {
				if ratio.VirtualPath == filename {
					actualPath = ratio.ActualPath
				}
			}
			if actualPath == "" {
				return errors.NotFound.New("actual file path is empty")
			}
			if _, err := os.Stat(actualPath); err != nil {
				return errors.NotFound.Wrapf(err, "file %s not found", actualPath)
			}
			file, err := os.Open(actualPath)
			if err != nil {
				return errors.Internal.Wrapf(err, "filed to open file: %s", actualPath)
			}
			_, err = rf.ReadFrom(file)
			if err != nil {
				return errors.Internal.Wrapf(err, "failed to read file %s", actualPath)
			}
			return nil
		}, nil,
	)
	return server, nil
}

//ReloadConfig for TFTP server
func (s *PinTFTPServer) ReloadConfig(config domain.TFTPConfig) error {
	s.config = config
	return nil
}

//ReloadPaths for TFTP paths server
func (s *PinTFTPServer) ReloadPaths(paths []domain.TFTPPathRatio) error {
	s.paths = &paths
	return nil
}

//Start TFTP server
func (s *PinTFTPServer) Start() error {
	go func() {
		s.state = domain.TFTPStateLaunched
		err := s.runtime.ListenAndServe(fmt.Sprintf("%s:%s", s.config.Address, s.config.Port))
		if err != nil {
			s.state = domain.TFTPStateError
		}
	}()
	return nil
}

//Stop TFTP server
func (s *PinTFTPServer) Stop() {
	if s.runtime != nil {
		s.runtime.Shutdown()
	}
	s.state = domain.TFTPStateStopped
}

//GetState from TFTP server
func (s *PinTFTPServer) GetState() domain.TFTPServerState {
	return s.state
}
