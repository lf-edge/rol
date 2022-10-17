package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
	"rol/dtos"
)

//ValidateDHCP4ServerUpdateDto validates dhcp v4 server update dto with ozzo-validation
//	Return
//	error - if an error occurs, otherwise nil
func ValidateDHCP4ServerUpdateDto(dto dtos.DHCP4ServerUpdateDto) error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.Enabled, []validation.Rule{
			validation.Required,
		}...),
		validation.Field(&dto.NTP, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpIPv4)).
				Error(regexpIPv4Desc),
		}...),
		validation.Field(&dto.Port, []validation.Rule{
			validation.Required,
		}...),
		validation.Field(&dto.DNS, []validation.Rule{
			validation.Required,
		}...),
		validation.Field(&dto.LeaseTime, []validation.Rule{
			validation.Required,
			validation.Min(60),
		}...),
	)
	return convertOzzoErrorToValidationError(err)
}
