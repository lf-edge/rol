package validators

import (
	"fmt"
	"strings"
)

//regexpUsername from 2 to 20 characters long, it can contain Latin uppercase and lowercase letters, as well as numbers. Must start with a letter
const regexpUsername = `^[a-zA-Z][a-zA-Z0-9-_\.]{1,20}$`
const regexpUsernameDesc = "From 2 to 20 characters long, it can contain Latin uppercase and lowercase letters, as well as numbers. Must start with a letter"

//regexpIPv4 IPv4 validation
const regexpIPv4 = `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`
const regexpIPv4Desc = "wrong IPv4 format"

func trimValidation(value interface{}) error {
	s, _ := value.(string)
	if strings.TrimSpace(s) != s {
		return fmt.Errorf("field cannot start or end with spaces")
	}
	return nil
}

func containsSpacesValidation(value interface{}) error {
	s, _ := value.(string)
	if strings.Contains(s, " ") {
		return fmt.Errorf("field cannot contain spaces")
	}
	return nil
}
