package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
	"rol/dtos"
)

//ValidateDHCP4LeaseUpdateDto validates dhcp v4 lease update dto with ozzo-validation
//	Return
//	error - if an error occurs, otherwise nil
func ValidateDHCP4LeaseUpdateDto(dto dtos.DHCP4LeaseUpdateDto) error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.IP, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpIPv4)).
				Error(regexpIPv4Desc),
		}...),
		validation.Field(&dto.MAC, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpMac)).
				Error(regexpMacDesc),
		}...),
		validation.Field(&dto.Expires, []validation.Rule{
			validation.Required,
		}...),
	)
	return convertOzzoErrorToValidationError(err)
}
