package request

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// todo сделать пакет пере используемым для этот избавиться от зависимости с chi пакетом
// Todo Если в структуре встраивание двух структур и каждая с типом файла, то при ошибке удалится только те файла которые относятся к данной структуре, файлы других структур не будут удалены

const (
	FORM        = "form"
	QUERY       = "query"
	RequiredTag = "required"
	OptionalTag = "optional"
)

// Data is a struct for data
type Data struct {
	val reflect.Value
	typ reflect.Type
}

// HttpData is a struct for values from [*http.Request]
type HttpData struct {
	Values url.Values
	Files  Files
}

// Files is a map of file names to a slice of files
type Files map[string][]*os.File

type FieldData struct {
	fieldType  reflect.StructField
	fieldValue reflect.Value
	tags       FieldTags
}

type FieldTags struct {
	fieldName string
	fieldTag  string
}

// RemoveFiles removes any temporary files associated with a [Files].
// Log delete file errors.
func (f Files) RemoveFiles() {
	fmt.Println("test remove files")
	fmt.Println(len(f))
	fmt.Println(f == nil)

	for _, v := range f {
		if v != nil {
			for i := range v {
				err := os.Remove(v[i].Name())
				if err != nil {
					log.Printf("request: RemoveFiles: file remove error: %s", err)
					return
				}
			}
		}
	}
}

// Add adds a file [*os.File] by key fileName to the [Files] map.
func (f Files) Add(fileName string, file *os.File) {
	f[fileName] = append(f[fileName], file)
}

// SetNil adds a nil file by key fileName to the [Files] map.
// To avoid deleting the necessary files, they must be deleted from [Files].
func (f Files) SetNil(mErr *MultiError, fileName string) {
	if mErr.err == nil {
		f[fileName] = nil
	}
}

func (data *Data) fieldNames() []string {
	var result []string

	for field := range data.typ.NumField() {
		if data.typ.Field(field).Type.Kind() == reflect.Struct {
			newData := Data{val: data.val.Field(field), typ: data.typ.Field(field).Type}
			names := newData.fieldNames()
			// todo if new data has no fields
			result = append(result, names...)
			continue
		}
		result = append(result, strings.ToLower(data.typ.Field(field).Name))
	}
	return result
}

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

// dataCreate validate data and return [*Data].
// Return only debug error.
func dataCreate(data any) (*Data, error) {
	if data == nil {
		return nil, fmt.Errorf("request: data cannot be nil")
	}

	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("request: data need to be pointer to the structure")
	}

	if val.IsNil() {
		return nil, fmt.Errorf("request: data must be a non-nil pointer")
	}

	if val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("request: transmitted pointer must point to the structure")
	}

	return &Data{val: val.Elem(), typ: val.Elem().Type()}, nil
}

// setField set payload to field [reflect.Value] of data structure.
// Return type convert error.
func setField(val reflect.Value, payload string, fieldName string) error {
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

func fieldData(data *Data, fieldNum int, tagName string) (FieldData, error) {
	fd := FieldData{
		fieldType:  data.typ.Field(fieldNum),
		fieldValue: data.val.Field(fieldNum),
	}
	err := fd.getFieldTags(tagName)
	if err != nil {
		return fd, err
	}
	return fd, err
}

// todo add trimSpaces
// getFieldTags get field tags by tag name 'form' or 'query'.
// Return only debug error.
func (fieldData *FieldData) getFieldTags(tagName string) (err error) {
	tags := strings.Split(fieldData.fieldType.Tag.Get(tagName), ",")
	if tags[0] == "" {
		return fmt.Errorf("request: getFieldTags: field '%s' has no tags for key '%s'", fieldData.fieldType.Name, tagName)
	}
	switch len(tags) {
	case 1:
		fieldData.tags.fieldName = tags[0]
		return nil
		//return &FieldTags{tags[0], ""}, nil
	case 2:
		switch tags[1] {
		case "":
			fieldData.tags.fieldName = tags[0]
			return nil
			//return &FieldTags{tags[0], ""}, nil
			//return "", "", fmt.Errorf("request: getFieldTags: invalid tags for field '%s' has no tags for key '%s'", typStruct.Name, tagName)
		case RequiredTag, OptionalTag:
			fieldData.tags.fieldName = tags[0]
			fieldData.tags.fieldTag = tags[1]
			return nil
			//return &FieldTags{tags[0], tags[1]}, nil
		default:
			return fmt.Errorf("request: getFieldTags: unsuported tag '%s' for field '%s' for key '%s'", tags[1], fieldData.fieldType.Name, tagName)
		}
	default:
		if len(tags) > 2 {
			return fmt.Errorf("request: getFieldTags: field '%s' has too many tags for key '%s', only 2 allowed", fieldData.fieldType.Name, tagName)
		}
		return fmt.Errorf("request: getFieldTags: unknown error: field: '%s', key: '%s'", fieldData.fieldType.Name, tagName)
	}
}
