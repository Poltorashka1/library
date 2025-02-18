package repo

//
//import (
//	"database/sql"
//	"errors"
//	"fmt"
//	"github.com/georgysavva/scany/pgxscan"
//	"reflect"
//	"strings"
//)
//
//const (
//	//form        = "form"
//	//query       = "query"
//	requiredTag = "required"
//	optionalTag = "optional"
//)
//
//type data struct {
//	val reflect.Value
//	typ reflect.Type
//}
//
//func newData(payload any) (*data, error) {
//	if payload == nil {
//		return nil, fmt.Errorf("request: payload cannot be nil")
//	}
//
//	val := reflect.ValueOf(payload)
//
//	if val.Kind() != reflect.Pointer {
//		return nil, fmt.Errorf("request: payload need to be pointer to the structure")
//	}
//
//	if val.IsNil() {
//		return nil, fmt.Errorf("request: payload must be a non-nil pointer")
//	}
//
//	if val.Elem().Kind() != reflect.Struct {
//		return nil, fmt.Errorf("request: transmitted pointer must point to the structure")
//	}
//
//	return &data{val: val.Elem(), typ: val.Elem().Type()}, nil
//}
//
//type field struct {
//	typ reflect.StructField
//	val reflect.Value
//
//	tagType string
//	tags    fieldTags
//}
//
//type fieldTags struct {
//	// Name - field tag name
//	Name string
//	// Tag - field tag value(required/optional)
//	Tag string
//}
//
//func (d *data) newField(fieldNum int) (*field, error) {
//	field := &field{
//		typ: d.typ.Field(fieldNum),
//		val: d.val.Field(fieldNum),
//	}
//	err := field.getFieldTags()
//	if err != nil {
//		return nil, err
//	}
//	return field, nil
//}
//
//// getFieldTags get fieldTags by tag name 'form' or 'query'.
//// Return only Error level error.
//func (f *field) getFieldTags() (err error) {
//	if f.tagType == "" {
//		return fmt.Errorf("request: getFieldTags: tagType cannot be empty")
//	}
//
//	tags := strings.Split(f.typ.Tag.Get(f.tagType), ",")
//	if tags[0] == "" {
//		if f.typ.Type.Kind() == reflect.Struct {
//			return nil
//		}
//		return fmt.Errorf("request: getFieldTags: field '%s' has no tags for key '%s'", f.typ.Name, f.tagType)
//	}
//	switch len(tags) {
//	case 1:
//		f.tags.Name = tags[0]
//		return nil
//		//return &FieldTags{tags[0], ""}, nil
//	case 2:
//		switch tags[1] {
//		case "":
//			f.tags.Name = tags[0]
//			return nil
//			//return &FieldTags{tags[0], ""}, nil
//			//return "", "", fmt.Errorf("request: getFieldTags: invalid tags for field '%s' has no tags for key '%s'", typStruct.Name, tagName)
//		case requiredTag, optionalTag:
//			f.tags.Name = tags[0]
//			f.tags.Tag = tags[1]
//			return nil
//			//return &FieldTags{tags[0], tags[1]}, nil
//		default:
//			return fmt.Errorf("request: getFieldTags: unsuported tag '%s' for field '%s' for key '%s'", tags[1], f.typ.Name, f.tagType)
//		}
//	default:
//		if len(tags) > 2 {
//			return fmt.Errorf("request: getFieldTags: field '%s' has too many tags for key '%s', only 2 allowed", f.typ.Name, f.tagType)
//		}
//		return fmt.Errorf("request: getFieldTags: unknown error: field: '%s', key: '%s'", f.typ.Name, f.tagType)
//	}
//}
//
//func ScanOne(rows *sql.Rows, payload any) error {
//	data, err := newData(payload)
//	if err != nil {
//		return err
//	}
//
//	for i := range data.typ.NumField() {
//		//field, err := data.newField(i)
//		//if err != nil {
//		//	return err
//		//}
//
//	}
//
//	if !rows.Next() {
//		return errors.New("no rows in result set")
//	}
//	pgxscan.ScanOne()
//	for {
//		rows.Scan()
//		if !rows.Next() {
//			break
//		}
//	}
//
//}
