package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/dtos"
)

//ValidateHostNetworkVlanCreateDto validates host network create dto
//	Return
//	error - if an error occurs, otherwise nil
func ValidateHostNetworkVlanCreateDto(dto dtos.HostNetworkVlanCreateDto) error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.VlanID, []validation.Rule{
			validation.Required,
			validation.Max(4098),
			validation.Min(1),
		}...),
		validation.Field(&dto.Parent, []validation.Rule{
			validation.Required,
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
		}...),
		validation.Field(&dto.Addresses, []validation.Rule{
			validation.By(sliceOfCidrStringsValidation),
		}...))
	return convertOzzoErrorToValidationError(err)
}
