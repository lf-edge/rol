package validators

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
	"rol/dtos"
)

func ValidateEthernetSwitchUpdateDto(dto dtos.EthernetSwitchUpdateDto) error {
	expr, _ := regexp.Compile(`^[\w\d.]*$`)
	fieldRule := []validation.Rule{
		validation.Required,
		validation.Match(expr).Error("field cannot start or end with spaces"),
	}
	return validation.ValidateStruct(&dto,
		validation.Field(&dto.Name, fieldRule...),
		validation.Field(&dto.Address, fieldRule...),
		validation.Field(&dto.Serial, fieldRule...),
		validation.Field(&dto.Username, fieldRule...),
		validation.Field(&dto.Password, fieldRule...),
	)
}
