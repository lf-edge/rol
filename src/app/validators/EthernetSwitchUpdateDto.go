package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
	"rol/app/errors"
	"rol/dtos"
)

//ValidateEthernetSwitchUpdateDto validates switch update dto with ozzo-validation
//	Return
//	error - if an error occurs, otherwise nil
func ValidateEthernetSwitchUpdateDto(dto dtos.EthernetSwitchUpdateDto) error {
	var err error
	validationErr := validation.ValidateStruct(&dto,
		validation.Field(&dto.Name, []validation.Rule{
			validation.Required,
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
		}...),
		validation.Field(&dto.Address, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpIPv4)).
				Error(regexpIPv4Desc),
		}...),
		validation.Field(&dto.Serial, []validation.Rule{
			validation.Required,
			validation.By(trimValidation),
			validation.By(containsSpacesValidation),
		}...),
		validation.Field(&dto.Username, []validation.Rule{
			validation.Required,
			validation.Match(regexp.MustCompile(regexpUsername)).
				Error(regexpUsernameDesc),
		}...),
		validation.Field(&dto.Password, []validation.Rule{
			validation.Required,
			validation.Length(6, 60),
		}...),
		validation.Field(&dto.SwitchModel, []validation.Rule{
			validation.Required,
			validation.By(trimValidation),
		}...))
	if validationErr != nil {
		err = errors.Validation.New(errors.ValidationErrorMessage)
		for key, value := range validationErr.(validation.Errors) {
			err = errors.AddErrorContext(err, key, value.Error())
		}
	}
	return err
}
