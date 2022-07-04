package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/dtos"
)

//ValidateEthernetSwitchPortUpdateDto validates ethernet switch port update dto
//	Return
//	error - if an error occurs, otherwise nil
func ValidateEthernetSwitchPortUpdateDto(dto dtos.EthernetSwitchPortUpdateDto) error {
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Name, []validation.Rule{
			validation.Required,
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
		}...),
		validation.Field(&dto.POEType, []validation.Rule{
			validation.Required,
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
			validation.By(validatePOEType),
		}...))
}
