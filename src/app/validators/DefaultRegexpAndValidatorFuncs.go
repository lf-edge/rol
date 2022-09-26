package validators

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"net"
	"rol/app/errors"
	"strings"
)

//regexpUsername from 2 to 20 characters long, it can contain Latin uppercase and lowercase letters, as well as numbers. Must start with a letter
const regexpUsername = `^[a-zA-Z][a-zA-Z0-9-_\.]{1,20}$`
const regexpUsernameDesc = "From 2 to 20 characters long, it can contain Latin uppercase and lowercase letters, as well as numbers. Must start with a letter"

//regexpIPv4 IPv4 validation
const regexpIPv4 = `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`
const regexpIPv4Desc = "wrong IPv4 format"

func convertOzzoErrorToValidationError(err error) error {
	var custError error
	if err != nil {
		custError = errors.Validation.New(errors.ValidationErrorMessage)
		for key, value := range err.(validation.Errors) {
			custError = errors.AddErrorContext(custError, key, value.Error())
		}
	}
	return custError
}

func sliceOfCidrStringsValidation(value interface{}) error {
	addresses, _ := value.([]string)
	for _, addressStr := range addresses {
		_, _, err := net.ParseCIDR(addressStr)
		if err != nil {
			return fmt.Errorf("wrong address: %s", addressStr)
		}
	}
	return nil
}

func trimValidation(value interface{}) error {
	s, _ := value.(string)
	if strings.TrimSpace(s) != s {
		return errors.Validation.New("field cannot start or end with spaces")
	}
	return nil
}

func containsSpacesValidation(value interface{}) error {
	s, _ := value.(string)
	if strings.Contains(s, " ") {
		return errors.Validation.New("field cannot contain spaces")
	}
	return nil
}

func uuidSliceElemUniqueness(value interface{}) error {
	s, _ := value.([]uuid.UUID)
	keys := make(map[uuid.UUID]bool)
	for _, entry := range s {
		if _, found := keys[entry]; !found {
			keys[entry] = true
		} else {
			return errors.Validation.New("slice can only contains unique UUIDs")
		}
	}
	return nil
}

func uuidsUniqueWithinSlices(fSlice []uuid.UUID, sSlice []uuid.UUID) error {
	for _, fElem := range fSlice {
		for _, sElem := range sSlice {
			if fElem == sElem {
				return errors.Validation.New("uuid should be unique within both slices")
			}
		}
	}
	return nil
}
