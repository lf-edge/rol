package validators

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/dtos"
)

func validatePOEType(value interface{}) error {
	s, _ := value.(string)
	switch s {
	case "poe+":
	case "poe":
	case "passive24":
	case "none":
	default:
		return fmt.Errorf("wrong poe power type")

	}
	return nil
}

//ValidateEthernetSwitchPortCreateDto validates ethernet switch port create dto
//	Return
//	error - if an error occurs, otherwise nil
func ValidateEthernetSwitchPortCreateDto(dto dtos.EthernetSwitchPortCreateDto) error {
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
