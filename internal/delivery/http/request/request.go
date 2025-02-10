package request

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	FORM        = "form"
	QUERY       = "query"
	RequiredTag = "required"
	OptionalTag = "optional"
)

// Data is a struct for payload reflect.Type, reflect.Value and tagType(example 'form' or 'json')
type Data struct {
	val reflect.Value
	typ reflect.Type

	tagType  string
	httpData *HttpData
}

// HttpData is a struct for values from [*http.Request]
type HttpData struct {
	Values url.Values
	Files  Files
}

// Files is a map of file names to a slice of files
type Files map[string][]*os.File

type Field struct {
	typ reflect.StructField
	val reflect.Value

	tagType string
	tags    FieldTags

	httpData *HttpData
}

// FieldTags is a struct for struct field
type FieldTags struct {
	// Name - field tag name
	Name string
	// Tag - field tag value(required/optional)
	Tag string
}

// RemoveFiles removes any temporary files associated with a [Files].
// Log delete file errors.
func (f Files) RemoveFiles() {
	if f == nil {
		log.Printf("request: RemoveFiles: files is nil")
		return
	}

	if len(f) == 0 {
		return
	}

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

// fieldTags returns a slice of Data fieldTagName it need to validate request form fields names.
// Return only Error level error.
func (data *Data) fieldTags() ([]string, error) {
	tags := make([]string, 0)

	for field := range data.typ.NumField() {
		if data.typ.Field(field).Type.Kind() == reflect.Struct {
			newData := Data{val: data.val.Field(field), typ: data.typ.Field(field).Type, tagType: data.tagType}
			embeddedTags, err := newData.fieldTags()
			if err != nil {
				return nil, err
			}
			tags = append(tags, embeddedTags...)

			continue
		}
		fd := Field{typ: data.typ.Field(field), val: data.val.Field(field), tagType: data.tagType}

		// todo mb solo function to get only tag name
		err := fd.getFieldTags()
		if err != nil {
			return nil, err
		}
		tags = append(tags, fd.tags.Name)
	}
	return tags, nil
}

// data validate payload pointer to the struct and return *Data.
// tagType is tag name example 'form' or 'json' it depends on what request format is being read.
// Return only Error level error.
func data(payload any, tagType string) (*Data, error) {
	if payload == nil {
		return nil, fmt.Errorf("request: payload cannot be nil")
	}

	val := reflect.ValueOf(payload)

	if val.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("request: payload need to be pointer to the structure")
	}

	if val.IsNil() {
		return nil, fmt.Errorf("request: payload must be a non-nil pointer")
	}

	if val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("request: transmitted pointer must point to the structure")
	}

	return &Data{val: val.Elem(), typ: val.Elem().Type(), tagType: tagType}, nil
}

// setField set payload value into field value.
// Return type convert error.
func (field *Field) setField(payload string) error {
	switch field.val.Kind() {
	case reflect.String:
		field.val.SetString(payload)
	case reflect.Int:
		//if data == "" {
		//	data = "0"
		//}
		digit, err := strconv.Atoi(payload)
		if err != nil {
			return fmt.Errorf("the field '%s' must be a digit; ", field.tags.Name)
		}
		field.val.SetInt(int64(digit))
	case reflect.Float64:
		//if data == "" {
		//	data = "0"
		//}
		f, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			return fmt.Errorf("the field '%s' must be a float digit; ", field.tags.Name)
		}
		field.val.SetFloat(f)
	default:
		return fmt.Errorf("unknown field type: %s", field.tags.Name)
	}
	return nil
}

func field(data *Data, fieldNum int) (*Field, error) {
	field := &Field{
		typ:      data.typ.Field(fieldNum),
		val:      data.val.Field(fieldNum),
		tagType:  data.tagType,
		httpData: data.httpData,
	}
	err := field.getFieldTags()
	if err != nil {
		return nil, err
	}
	return field, nil
}

// todo add trimSpaces in tags and test it
// getFieldTags get field tags by tag name 'form' or 'query'.
// Return only Error level error.
func (field *Field) getFieldTags() (err error) {
	if field.tagType == "" {
		return fmt.Errorf("request: getFieldTags: tagType cannot be empty")
	}

	tags := strings.Split(field.typ.Tag.Get(field.tagType), ",")
	if tags[0] == "" {
		return fmt.Errorf("request: getFieldTags: field '%s' has no tags for key '%s'", field.typ.Name, field.tagType)
	}
	switch len(tags) {
	case 1:
		field.tags.Name = tags[0]
		return nil
		//return &FieldTags{tags[0], ""}, nil
	case 2:
		switch tags[1] {
		case "":
			field.tags.Name = tags[0]
			return nil
			//return &FieldTags{tags[0], ""}, nil
			//return "", "", fmt.Errorf("request: getFieldTags: invalid tags for field '%s' has no tags for key '%s'", typStruct.Name, tagName)
		case RequiredTag, OptionalTag:
			field.tags.Name = tags[0]
			field.tags.Tag = tags[1]
			return nil
			//return &FieldTags{tags[0], tags[1]}, nil
		default:
			return fmt.Errorf("request: getFieldTags: unsuported tag '%s' for field '%s' for key '%s'", tags[1], field.typ.Name, field.tagType)
		}
	default:
		if len(tags) > 2 {
			return fmt.Errorf("request: getFieldTags: field '%s' has too many tags for key '%s', only 2 allowed", field.typ.Name, field.tagType)
		}
		return fmt.Errorf("request: getFieldTags: unknown error: field: '%s', key: '%s'", field.typ.Name, field.tagType)
	}
}
