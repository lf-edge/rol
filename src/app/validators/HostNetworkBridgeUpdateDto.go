package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/dtos"
)

//ValidateHostNetworkBridgeUpdateDto validates host network bridge update dto
//	Return
//	error - if an error occurs, otherwise nil
func ValidateHostNetworkBridgeUpdateDto(dto dtos.HostNetworkBridgeUpdateDto) error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.Slaves, validation.Each(
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
		)),
		validation.Field(&dto.Addresses, []validation.Rule{
			validation.By(sliceOfCidrStringsValidation),
		}...))
	return convertOzzoErrorToValidationError(err)
}
