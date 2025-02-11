package request

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	fileType  = reflect.TypeOf((*os.File)(nil))
	filesType = reflect.TypeOf(([]*os.File)(nil))
)

const (
	//form        = "form"
	//query       = "query"
	requiredTag = "required"
	optionalTag = "optional"
)

// data is a struct for payload reflect.Type, reflect.Value and tagType(example 'form' or 'json')
type data struct {
	val reflect.Value
	typ reflect.Type

	tagType     string
	requestData *requestData
}

// requestData is a struct for values from request
type requestData struct {
	//Buff *bytes.Buffer
	Json string
	//Buff   bytes.Buffer
	Values url.Values
	Files  files
}

// files is a map of file names to a slice of os.File
type files map[string][]*os.File

// field is a struct for data structField
type field struct {
	typ reflect.StructField
	val reflect.Value

	tagType string
	tags    fieldTags

	requestData *requestData
}

// fieldTags is a struct with tags for field
type fieldTags struct {
	// Name - field tag name
	Name string
	// Tag - field tag value(required/optional)
	Tag string
}

// RemoveFiles removes any temporary files.
// Log delete file errors.
func (f files) RemoveFiles() {
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

// Add adds a file os.File by key fileName to the [Files] map.
func (f files) Add(fileName string, file *os.File) {
	f[fileName] = append(f[fileName], file)
}

// todo optimize add method in payload struct to get all field name
// fieldTags returns a slice of Data fieldTagName it need to validate request form fields names.
// Return only Error level error.
func (d *data) fieldTags() ([]string, error) {
	tags := make([]string, 0)

	for f := range d.typ.NumField() {
		if d.typ.Field(f).Type.Kind() == reflect.Struct {
			nData := data{val: d.val.Field(f), typ: d.typ.Field(f).Type, tagType: d.tagType}
			embeddedTags, err := nData.fieldTags()
			if err != nil {
				return nil, err
			}
			tags = append(tags, embeddedTags...)

			continue
		}
		fd := field{typ: d.typ.Field(f), val: d.val.Field(f), tagType: d.tagType}
		if fd.typ.IsExported() {
			// todo mb solo function to get only tag name
			err := fd.getFieldTags()
			if err != nil {
				return nil, err
			}
			tags = append(tags, fd.tags.Name)
		}
	}
	return tags, nil
}

// newData validate payload pointer to the struct and return data.
// tagType is tag name example 'form' or 'json' or 'json' it depends on what request format is being read.
// Return only Error level error.
func newData(payload any, tagType string) (*data, error) {
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

	return &data{val: val.Elem(), typ: val.Elem().Type(), tagType: tagType}, nil
}

// setField set payload value into field value.
// Return type convert error.
func (f *field) setField(payload string) error {
	switch f.val.Kind() {
	case reflect.String:
		f.val.SetString(payload)
	case reflect.Int:
		//if data == "" {
		//	data = "0"
		//}
		digit, err := strconv.Atoi(payload)
		if err != nil {
			return fmt.Errorf("the field '%s' must be a digit; ", f.tags.Name)
		}
		f.val.SetInt(int64(digit))
	case reflect.Float64:
		//if data == "" {
		//	data = "0"
		//}
		parsF, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			return fmt.Errorf("the field '%s' must be a float digit; ", f.tags.Name)
		}
		f.val.SetFloat(parsF)
	default:
		return fmt.Errorf("unknown field type: %s", f.tags.Name)
	}
	return nil
}

// newField set data structField value into field value.
// Return only Error level error.
func (d *data) newField(fieldNum int) (*field, error) {
	field := &field{
		typ:         d.typ.Field(fieldNum),
		val:         d.val.Field(fieldNum),
		tagType:     d.tagType,
		requestData: d.requestData,
	}
	err := field.getFieldTags()
	if err != nil {
		return nil, err
	}
	return field, nil
}

// getFieldTags get fieldTags by tag name 'form' or 'query'.
// Return only Error level error.
func (f *field) getFieldTags() (err error) {
	if f.tagType == "" {
		return fmt.Errorf("request: getFieldTags: tagType cannot be empty")
	}

	tags := strings.Split(f.typ.Tag.Get(f.tagType), ",")
	if tags[0] == "" {
		if f.typ.Type.Kind() == reflect.Struct {
			return nil
		}
		return fmt.Errorf("request: getFieldTags: field '%s' has no tags for key '%s'", f.typ.Name, f.tagType)
	}
	switch len(tags) {
	case 1:
		f.tags.Name = tags[0]
		return nil
		//return &FieldTags{tags[0], ""}, nil
	case 2:
		switch tags[1] {
		case "":
			f.tags.Name = tags[0]
			return nil
			//return &FieldTags{tags[0], ""}, nil
			//return "", "", fmt.Errorf("request: getFieldTags: invalid tags for field '%s' has no tags for key '%s'", typStruct.Name, tagName)
		case requiredTag, optionalTag:
			f.tags.Name = tags[0]
			f.tags.Tag = tags[1]
			return nil
			//return &FieldTags{tags[0], tags[1]}, nil
		default:
			return fmt.Errorf("request: getFieldTags: unsuported tag '%s' for field '%s' for key '%s'", tags[1], f.typ.Name, f.tagType)
		}
	default:
		if len(tags) > 2 {
			return fmt.Errorf("request: getFieldTags: field '%s' has too many tags for key '%s', only 2 allowed", f.typ.Name, f.tagType)
		}
		return fmt.Errorf("request: getFieldTags: unknown error: field: '%s', key: '%s'", f.typ.Name, f.tagType)
	}
}

// setDataValue writes value to fields of a data structure.
// Return error - MultiError.
func (d *data) setDataValue(mErr *MultiError) error {
	if mErr == nil {
		return errors.New("request: setDataValue: mErr cannot be nil")
	}
	if d.requestData == nil {
		return errors.New("request: setDataValue: httpData cannot be nil")
	}

	for fieldNum := range d.typ.NumField() {
		field, err := d.newField(fieldNum)
		if err != nil {
			return err
		}
		if field.typ.IsExported() {
			switch field.typ.Type {
			case fileType:
				err := field.setFileValue(mErr)
				if err != nil {
					return err
				}
			case filesType:
				err := field.setFilesValue(mErr)
				if err != nil {
					mErr.Add(err)
					return err
				}
			default:
				err := field.setDefaultValue(mErr)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// setFileValue writes file os.File to field value.
func (f *field) setFileValue(mErr *MultiError) error {
	file := f.requestData.Files[f.tags.Name]

	if f.tags.Tag == requiredTag && len(file) < 1 || len(file) > 1 {
		mErr.Add(fmt.Errorf("the field '%s' must be a single file; ", f.tags.Name))
		return nil
	}
	if f.tags.Tag == optionalTag && len(file) > 1 {
		mErr.Add(fmt.Errorf("the field '%s' must be a single file; ", f.tags.Name))
		return nil
	}

	if file != nil {
		f.val.Set(reflect.ValueOf(file[0]))
	}

	return nil
}

// setFilesValue writes files os.File to field value.
func (f *field) setFilesValue(mErr *MultiError) error {
	files := f.requestData.Files[f.tags.Name]
	if f.tags.Tag == requiredTag && len(files) < 1 {
		mErr.Add(fmt.Errorf("the field '%s' must be a single or more files; ", f.tags.Name))
		return nil
	}

	if files != nil {
		f.val.Set(reflect.ValueOf(files))
	}
	return nil
}

// setDefaultValue writes value to field.
func (f *field) setDefaultValue(mErr *MultiError) error {
	if f == nil {
		return errors.New("request: setDefaultValue: field cannot be nil")
	}
	if f.val.Kind() == reflect.Struct {
		d := &data{
			val:         f.val,
			typ:         f.typ.Type,
			tagType:     f.tagType,
			requestData: f.requestData,
		}

		err := d.setDataValue(mErr)
		if err != nil {
			return err
		}
		f.val.Set(d.val)
		return nil
	}

	ok := f.requestData.Values.Has(f.tags.Name)
	if !ok && f.tags.Tag == requiredTag {
		//mErr.Add(ErrFieldRequired{field: f.tags.Name})
		mErr.Add(fmt.Errorf("the field '%s' is required; ", f.tags.Name))
		return nil
	}

	if f.tags.Tag == optionalTag && f.requestData.Values.Get(f.tags.Name) == "" {
		return nil
	}

	err := f.setField(f.requestData.Values.Get(f.tags.Name))
	if err != nil {
		if f.tags.Tag == optionalTag || f.tags.Tag == requiredTag {
			mErr.Add(err)
			return nil
		}
		return nil
	}
	return nil
}
