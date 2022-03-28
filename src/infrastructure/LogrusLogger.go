package infrastructure

import (
	"os"
	"rol/domain"

	"github.com/sirupsen/logrus"
)

//NewLogrusLogger creates new instance of the logrus logger with two hooks for end-to-end logging
//Params
//	config - yaml configuration struct *domain.AppConfig
//Return
//	*logrus.Logger - logrus logger instance
//	error - if error occurs return error, otherwise nil
func NewLogrusLogger(config *domain.AppConfig) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&Formatter{})

	var err error
	logger.Level, err = logrus.ParseLevel(config.Logger.Level)
	if err != nil {
		return nil, err
	}
	return logger, nil
}
