package request

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// todo сделать пакет пере используемым для этот избавиться от зависимости с chi пакетом

const (
	PDF   = "application/pdf"
	EPUB  = "application/epub+zip"
	FORM  = "form"
	QUERY = "query"
)

// BodyParse parse request body in Json or Form format into pointer struct
func BodyParse(r *http.Request, data any) error {
	switch r.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		err := FormParse(r, data)
		if err != nil {
			return err
		}
	case "application/json":
		err := JsonParse(r, data)
		if err != nil {
			return err
		}
	default:
		return ErrUnknownContentType
	}

	return nil
}

// dataValidate validate data and return reflect.Value of pointer value.
// Return only debug error.
func dataValidate(data any) (reflect.Value, error) {
	if data == nil {
		return reflect.Value{}, fmt.Errorf("request: data cannot be nil")
	}

	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Pointer {
		return reflect.Value{}, fmt.Errorf("request: data need to be pointer to the structure")
	}

	if val.IsNil() {
		return reflect.Value{}, fmt.Errorf("request: data must be a non-nil pointer")
	}

	if val.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("request: transmitted pointer must point to the structure")
	}

	val = val.Elem()
	return val, nil
}

// setFieldData set payload to field [reflect.Value] of data structure.
// Return type convert error.
func setFieldData(val reflect.Value, payload string, fieldName string) error {
	switch val.Kind() {
	case reflect.String:
		val.SetString(payload)
	case reflect.Int:
		//if data == "" {
		//	data = "0"
		//}
		digit, err := strconv.Atoi(payload)
		if err != nil {
			return fmt.Errorf("the field '%s' must be a digit; ", fieldName)
		}
		val.SetInt(int64(digit))
	case reflect.Float64:
		//if data == "" {
		//	data = "0"
		//}
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			return fmt.Errorf("the field '%s' must be a float digit; ", fieldName)
		}
		val.SetFloat(f)
	default:
		return fmt.Errorf("unknown field type: %s", fieldName)
	}
	return nil
}

// getFieldTags get field tags by tag name 'form' or 'query'.
// Return only debug error.
func getFieldTags(typStruct *reflect.StructField, tagName string) (fieldName string, tag string, err error) {
	tags := strings.Split(typStruct.Tag.Get(tagName), ",")
	if len(tags) == 1 {
		if tags[0] == "" {
			return "", "", fmt.Errorf("request: getFieldTags: field '%s' has no tags for key '%s'", typStruct.Name, tagName)
		}
		return tags[0], "", nil
	}
	if len(tags) >= 2 {
		return tags[0], tags[1], nil
	}
	return "", "", nil
}
