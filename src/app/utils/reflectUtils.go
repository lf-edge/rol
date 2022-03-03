package utils

import (
	"reflect"
)

func GetStringFieldsNames(value interface{}, stringFieldNames *[]string) {
	valueOf := reflect.ValueOf(value)
	typeOf := reflect.TypeOf(value)
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	for i := 0; i < valueOf.NumField(); i++ {
		fieldValue := valueOf.Field(i)
		if fieldValue.Kind() == reflect.  {
			GetStringFieldsNames(fieldValue, stringFieldNames)
			continue
		}
		if fieldValue.Kind() == reflect.String {
			*stringFieldNames = append(*stringFieldNames, typeOf.Field(i).Name)
		}
	}
}
