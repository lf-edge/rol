package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/dtos"
)

//ValidateHostNetworkBridgeCreateDto validates host network bridge create dto
//	Return
//	error - if an error occurs, otherwise nil
func ValidateHostNetworkBridgeCreateDto(dto dtos.HostNetworkBridgeCreateDto) error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.Name, []validation.Rule{
			validation.Required,
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
		}...),
		validation.Field(&dto.Slaves, validation.Each(
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
		)),
		validation.Field(&dto.Addresses, []validation.Rule{
			validation.By(sliceOfCidrStringsValidation),
		}...))
	return convertOzzoErrorToValidationError(err)
}
