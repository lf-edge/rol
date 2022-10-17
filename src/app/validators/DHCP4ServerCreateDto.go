package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
	"rol/dtos"
)

//ValidateDHCP4ServerCreateDto validates dhcp v4 server create dto with ozzo-validation
//	Return
//	error - if an error occurs, otherwise nil
func ValidateDHCP4ServerCreateDto(dto dtos.DHCP4ServerCreateDto) error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.Interface, []validation.Rule{
			validation.Required,
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
		}...),
		validation.Field(&dto.Enabled, []validation.Rule{
			validation.Required,
		}...),
		validation.Field(&dto.Mask, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpIPv4)).
				Error(regexpIPv4Desc),
		}...),
		validation.Field(&dto.ServerID, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpIPv4)).
				Error(regexpIPv4Desc),
		}...),
		validation.Field(&dto.Gateway, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpIPv4)).
				Error(regexpIPv4Desc),
		}...),
		validation.Field(&dto.NTP, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpIPv4)).
				Error(regexpIPv4Desc),
		}...),
		validation.Field(&dto.Range, []validation.Rule{
			validation.Required,
		}...),
		validation.Field(&dto.DNS, []validation.Rule{
			validation.Required,
		}...),
		validation.Field(&dto.Port, []validation.Rule{
			validation.Required,
		}...),
		validation.Field(&dto.LeaseTime, []validation.Rule{
			validation.Required,
			validation.Min(60),
		}...),
	)
	return convertOzzoErrorToValidationError(err)
}
