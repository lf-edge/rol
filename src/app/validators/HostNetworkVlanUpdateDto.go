package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/dtos"
)

//ValidateHostNetworkVlanUpdateDto validates host network update dto
//	Return
//	error - if an error occurs, otherwise nil
func ValidateHostNetworkVlanUpdateDto(dto dtos.HostNetworkVlanUpdateDto) error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.Addresses, []validation.Rule{
			validation.Required,
			validation.By(sliceOfCidrStringsValidation),
		}...))
	return convertOzzoErrorToValidationError(err)
}
