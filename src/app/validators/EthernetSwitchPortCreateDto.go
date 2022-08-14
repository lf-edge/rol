package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/app/errors"
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
		return errors.Internal.New("wrong poe power type")

	}
	return nil
}

//ValidateEthernetSwitchPortCreateDto validates ethernet switch port create dto
//	Return
//	error - if an error occurs, otherwise nil
func ValidateEthernetSwitchPortCreateDto(dto dtos.EthernetSwitchPortCreateDto) error {
	var err error
	validationErr := validation.ValidateStruct(&dto,
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
	if validationErr != nil {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		for key, value := range validationErr.(validation.Errors) {
			err = errors.AddErrorContext(err, key, value.Error())
		}
	}
	return err
}
