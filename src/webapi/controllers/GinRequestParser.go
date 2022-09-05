package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reflect"
	"rol/app/errors"
	"strconv"
)

type paginatedRequestStructForParsing struct {
	Page           int    `query:"page"`
	PageSize       int    `query:"pageSize"`
	OrderBy        string `query:"orderBy"`
	OrderDirection string `query:"orderDirection"`
	Search         string `query:"search"`
}

//newPaginatedRequestStructForParsing create new paginatedRequestStructForParsing with default values
func newPaginatedRequestStructForParsing(page, size int, orderBy, orderDirection, search string) paginatedRequestStructForParsing {
	return paginatedRequestStructForParsing{
		Page:           page,
		PageSize:       size,
		OrderBy:        orderBy,
		OrderDirection: orderDirection,
		Search:         search,
	}
}

func hasTag(v any, tag string) bool {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	return typeHasTag(t, tag)
}

func typeHasTag(t reflect.Type, tag string) bool {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.PkgPath != "" {
			continue // skip unexported fields
		}

		if _, ok := f.Tag.Lookup(tag); ok {
			return true
		}

		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct {
			if typeHasTag(ft, tag) {
				return true
			}
		}
	}

	return false
}

//Convert string to chosen type
//
//Params:
//	string - string for convert
//Return:
//	T - converted value
//	error - if an error occurs, otherwise nil
func Convert[T any](s string) (T, error) {
	z := new(T)
	rt := reflect.TypeOf(*z)

	reflectValue := reflect.ValueOf(z).Elem()
	switch rt.Kind() {
	case reflect.Float32:
		t, err := strconv.ParseFloat(s, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(t))
		return *z, nil
	case reflect.Float64:
		t, err := strconv.ParseFloat(s, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(t))
		return *z, nil
	case reflect.Int:
		t, err := strconv.ParseInt(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(int(t)))
		return *z, nil
	case reflect.Int8:
		t, err := strconv.ParseInt(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(int8(t)))
		return *z, nil
	case reflect.Int16:
		t, err := strconv.ParseInt(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(int16(t)))
		return *z, nil
	case reflect.Int32:
		t, err := strconv.ParseInt(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(int32(t)))
		return *z, nil
	case reflect.Int64:
		t, err := strconv.ParseInt(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(t))
		return *z, nil
	case reflect.Uint:
		t, err := strconv.ParseUint(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(uint(t)))
		return *z, nil
	case reflect.Uint8:
		t, err := strconv.ParseUint(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(uint8(t)))
		return *z, nil
	case reflect.Uint16:
		t, err := strconv.ParseUint(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(uint16(t)))
		return *z, nil
	case reflect.Uint32:
		t, err := strconv.ParseUint(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(uint32(t)))
		return *z, nil
	case reflect.Uint64:
		t, err := strconv.ParseUint(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(t))
		return *z, nil
	case reflect.Uintptr:
		t, err := strconv.ParseUint(s, 10, rt.Bits())
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(uintptr(t)))
		return *z, nil
	case reflect.TypeOf(uuid.UUID{}).Kind():
		t, err := uuid.Parse(s)
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(t))
		return *z, nil
	case reflect.Bool:
		t, err := strconv.ParseBool(s)
		if err != nil {
			return *z, err
		}
		reflectValue.Set(reflect.ValueOf(t))
		return *z, nil
	case reflect.String:
		reflectValue.Set(reflect.ValueOf(s))
		return *z, nil
	default:
		return *z, errors.Internal.New("not supported conversion")
	}
}

//GetValueByKeyFromGin get value by key from gin context (header, param, query)
//
//Params:
//  *gin.Context - gin context
//	string - source for parsing (header, param or query)
//	string - key for parsing
//  T - default value, if value is not exist in context
//Return:
//	T - value or default value
//	error - if an error occurs, otherwise nil
func GetValueByKeyFromGin[T any](ctx *gin.Context, source, key string, def T) (T, error) {
	switch source {
	case "query":
		strVal := ctx.DefaultQuery(key, "")
		if strVal != "" {
			return Convert[T](strVal)
		}
		return def, nil
	case "param":
		strVal := ctx.Param(key)
		if strVal != "" {
			return Convert[T](strVal)
		}
		return def, nil
	case "header":
		strVal := ctx.GetHeader(key)
		if strVal != "" {
			return Convert[T](strVal)
		}
		return def, nil
	default:
		return def, errors.Internal.New("not correct parsing source")
	}
}

//parseGinRequest parse gin request to request params compilation object.
//
//Parsing query, params, and headers, but no BODY!.
//In the params compilation object, we need to set correct flags, for example:
//	type TestRequestCompilation struct {
//		Name      string    `query:"name"`
//		Header    string    `header:"X-TEST-HEADER"`
//		Page      int       `query:"page"`
//		UUID      uuid.UUID `param:"id"`
//	}
// if you fill in the data in the fields of this structure, then these data will be the default values.
//
//Params:
//  *gin.Context - gin context
//	any - any pointer to object with correct tagged fields
//Return:
//	error - if an error occurs, otherwise nil
func parseGinRequest(ctx *gin.Context, reqParamsCompilation any) error {
	if reqParamsCompilation == nil {
		return nil
	}
	if hasTag(reqParamsCompilation, "query") {
		err := parseKeysValuesToObjectFromGin(ctx, reqParamsCompilation, "query")
		if err != nil {
			return err
		}
	}
	if hasTag(reqParamsCompilation, "param") {
		err := parseKeysValuesToObjectFromGin(ctx, reqParamsCompilation, "param")
		if err != nil {
			return err
		}
	}
	if hasTag(reqParamsCompilation, "header") {
		err := parseKeysValuesToObjectFromGin(ctx, reqParamsCompilation, "header")
		if err != nil {
			return err
		}
	}
	return nil
}

func parseKeysValuesToObjectFromGin(ctx *gin.Context, obj any, parsingSource string) error {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.PkgPath != "" {
			continue // skip unexported fields
		}
		parsingName := ""
		name, ok := f.Tag.Lookup(parsingSource)
		if !ok {
			continue
		}
		if len(name) > 0 {
			parsingName = name
		} else {
			parsingName = f.Name
		}

		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		defaultValue := v.FieldByName(f.Name).Interface()

		var parsedValue interface{}
		var err error

		k := ft.Kind()
		switch k {
		case reflect.Struct:
			err = parseKeysValuesToObjectFromGin(ctx, ft, parsingSource)
			if err != nil {
				continue
			}
		case reflect.String:
			parsedValue, err = GetValueByKeyFromGin[string](ctx, parsingSource, parsingName, defaultValue.(string))
			if err != nil {
				continue
			}
		case reflect.Int:
			parsedValue, err = GetValueByKeyFromGin[int](ctx, parsingSource, parsingName, defaultValue.(int))
			if err != nil {
				continue
			}
		case reflect.Int8:
			parsedValue, err = GetValueByKeyFromGin[int8](ctx, parsingSource, parsingName, defaultValue.(int8))
			if err != nil {
				continue
			}
		case reflect.Int16:
			parsedValue, err = GetValueByKeyFromGin[int16](ctx, parsingSource, parsingName, defaultValue.(int16))
			if err != nil {
				continue
			}
		case reflect.Int32:
			parsedValue, err = GetValueByKeyFromGin[int32](ctx, parsingSource, parsingName, defaultValue.(int32))
			if err != nil {
				continue
			}
		case reflect.Int64:
			parsedValue, err = GetValueByKeyFromGin[int64](ctx, parsingSource, parsingName, defaultValue.(int64))
			if err != nil {
				continue
			}
		case reflect.Uint:
			parsedValue, err = GetValueByKeyFromGin[uint](ctx, parsingSource, parsingName, defaultValue.(uint))
			if err != nil {
				continue
			}
		case reflect.Uint8:
			parsedValue, err = GetValueByKeyFromGin[uint8](ctx, parsingSource, parsingName, defaultValue.(uint8))
			if err != nil {
				continue
			}
		case reflect.Uint16:
			parsedValue, err = GetValueByKeyFromGin[uint16](ctx, parsingSource, parsingName, defaultValue.(uint16))
			if err != nil {
				continue
			}
		case reflect.Uint32:
			parsedValue, err = GetValueByKeyFromGin[uint32](ctx, parsingSource, parsingName, defaultValue.(uint32))
			if err != nil {
				continue
			}
		case reflect.Uint64:
			parsedValue, err = GetValueByKeyFromGin[uint64](ctx, parsingSource, parsingName, defaultValue.(uint64))
			if err != nil {
				continue
			}
		case reflect.Float32:
			parsedValue, err = GetValueByKeyFromGin[float32](ctx, parsingSource, parsingName, defaultValue.(float32))
			if err != nil {
				continue
			}
		case reflect.Float64:
			parsedValue, err = GetValueByKeyFromGin[float64](ctx, parsingSource, parsingName, defaultValue.(float64))
			if err != nil {
				continue
			}
		case reflect.Bool:
			parsedValue, err = GetValueByKeyFromGin[bool](ctx, parsingSource, parsingName, defaultValue.(bool))
			if err != nil {
				continue
			}
		case reflect.TypeOf(uuid.UUID{}).Kind():
			parsedValue, err = GetValueByKeyFromGin[uuid.UUID](ctx, parsingSource, parsingName, defaultValue.(uuid.UUID))
			if err != nil {
				continue
			}
		default:
			return errors.Internal.Newf("not supported type %s", k)
		}
		if parsedValue != nil {
			tt := v.FieldByName(f.Name)
			tt.Set(reflect.ValueOf(parsedValue))
		}
	}
	return nil
}
