package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"rol/app/errors"
	"rol/dtos"
)

//ValidateEthernetSwitchPortUpdateDto validates ethernet switch port update dto
//	Return
//	error - if an error occurs, otherwise nil
func ValidateEthernetSwitchPortUpdateDto(dto dtos.EthernetSwitchPortUpdateDto) error {
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
