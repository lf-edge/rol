package tests

import (
	validation "github.com/go-ozzo/ozzo-validation"
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
			SwitchModel: 0,
			Address:     "123.123.123.123",
			Username:    "AutoUser",
		},
		Password: "AutoPass",
	}
	err := validators.ValidateEthernetSwitchCreateDto(dto)
	for fieldName, errDescription := range err.(validation.Errors) {
		if errDescription.Error() != blankErr && errDescription.Error() != whitespaceErr {
			t.Errorf("%s: %s", fieldName, errDescription)
		}
	}
}

func Test_EthernetSwitchUpdateDto(t *testing.T) {
	dto := dtos.EthernetSwitchUpdateDto{
		EthernetSwitchBaseDto: dtos.EthernetSwitchBaseDto{
			Name:        " Validate ",
			Serial:      "",
			SwitchModel: 0,
			Address:     "123.123.123.123",
			Username:    "AutoUser",
		},
		Password: "AutoPass",
	}
	err := validators.ValidateEthernetSwitchUpdateDto(dto)
	for fieldName, errDescription := range err.(validation.Errors) {
		if errDescription.Error() != blankErr && errDescription.Error() != whitespaceErr {
			t.Errorf("%s: %s", fieldName, errDescription)
		}
	}
}
