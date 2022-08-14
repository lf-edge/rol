package tests

import (
	"rol/app/errors"
	"rol/app/validators"
	"rol/dtos"
	"testing"
)

var (
	blankErr      = "cannot be blank"
	whitespaceErr = "field cannot start or end with spaces"
)

func Test_EthernetSwitchCreateDto(t *testing.T) {

	dto := dtos.EthernetSwitchCreateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        " Validate ",
			Serial:      "",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.123",
			Username:    "AutoUser",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	err := validators.ValidateEthernetSwitchCreateDto(dto)
	errContext := errors.GetErrorContext(err)
	for fieldName, errDescription := range errContext {
		if errDescription != blankErr && errDescription != whitespaceErr {
			t.Errorf("%s: %s", fieldName, errDescription)
		}
	}
}

func Test_EthernetSwitchUpdateDto(t *testing.T) {
	dto := dtos.EthernetSwitchUpdateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        " Validate ",
			Serial:      "",
			SwitchModel: "unifi_switch_us-24-250w",
			Address:     "123.123.123.123",
			Username:    "AutoUser",
		},
		//  pragma: allowlist nextline secret
		Password: "AutoPass",
	}
	err := validators.ValidateEthernetSwitchUpdateDto(dto)
	errContext := errors.GetErrorContext(err)
	for fieldName, errDescription := range errContext {
		if errDescription != blankErr && errDescription != whitespaceErr {
			t.Errorf("%s: %s", fieldName, errDescription)
		}
	}
}
