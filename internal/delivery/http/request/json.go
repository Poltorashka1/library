package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// todo validate data IMPORTNT
func JsonParse(r *http.Request, payload any) error {
	// todo if error is EOF
	_, err := newData(payload, "json")
	if err != nil {
		return err
	}

	err = json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		var eUnmarshal *json.UnmarshalTypeError
		if ok := errors.As(err, &eUnmarshal); ok {
			return &ErrFieldType{
				field: eUnmarshal.Field,
				typ:   eUnmarshal.Type,
			}
		}
		var eSyntax *json.SyntaxError
		if ok := errors.As(err, &eSyntax); ok {
			return ErrInvalidJsonSyntax
		}
		//var e *json.UnsupportedValueError
		//if ok := errors.As(err, &e); ok {
		//	return fmt.Errorf("invalid JSON value: %s", e.Error())
		//}
		return err
	}
	return nil
}

type jsonData map[string]any

type jsonParser struct {
	cfg jsonParserConfig

	json string
	data *data
}

type jsonParserConfig struct {
	maxBodySize int64
}

func JsonParseV2(r *http.Request, payload any) error {
	fmt.Println("JsonParseV2")
	err := jsonRequestValidate(r)
	if err != nil {
		return err
	}

	d, err := newData(payload, "json")
	if err != nil {
		return err
	}

	parser := &jsonParser{
		data: d,
	}

	err = parser.jsonParse(r)
	if err != nil {
		return err
	}

	var mErr = &MultiError{}
	err = d.decodeJson(mErr)
	if err != nil {
		return err
	}
	//var mErr = &MultiError{}
	//err = d.setDataValue(mErr)
	//if err != nil {
	//	return err
	//}
	//
	if mErr.err != nil {
		return mErr
	}
	return nil
}

// jsonRequestValidate validate request content type.
// Return error - ErrUnknownContentType.
func jsonRequestValidate(r *http.Request) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return ErrUnknownContentType
	}
	return nil
}

func (parser *jsonParser) jsonParse(r *http.Request) error {
	if parser.cfg.maxBodySize > 0 {
		if r.ContentLength > parser.cfg.maxBodySize {
			return &ErrContentToLarge{limit: parser.cfg.maxBodySize}
		}
		r.Body = http.MaxBytesReader(nil, r.Body, parser.cfg.maxBodySize)
	}

	err := parser.readBody(r)
	if err != nil {
		return err
	}

	//err = parser.parse()
	//if err != nil {
	//	return err
	//}

	return nil
}

//func (parser *jsonParser) parse() error {
//
//}

func (parser *jsonParser) readBody(r *http.Request) error {
	var result bytes.Buffer

	if parser.cfg.maxBodySize > 0 {
		result.Grow(int(parser.cfg.maxBodySize))
	}

	buf := make([]byte, 4096)
	for {
		n, readErr := r.Body.Read(buf)
		if readErr != nil && readErr != io.EOF {
			var maxBytesErr *http.MaxBytesError
			if errors.As(readErr, &maxBytesErr) {
				return &ErrContentToLarge{limit: parser.cfg.maxBodySize}
			}
			return readErr
		}

		if n > 0 {
			_, writeErr := result.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
		}

		if readErr == io.EOF {
			break
		}
	}

	parser.data.requestData = &requestData{
		Json: result.String(),
	}
	return nil
}

//func (parser *jsonParser) setValue(mErr *MultiError) error {
//	if mErr == nil {
//		return errors.New("request: json: setValue: mErr cannot be nil")
//	}
//	err := parser.decodeJson()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (d *data) decodeJson(mErr *MultiError) error {
	if mErr == nil {
		return errors.New("request: json: setValue: mErr cannot be nil")
	}

	//inObject := false

	for _, v := range d.requestData.Json {
		switch string(v) {
		case `{`:
			//inObject = true
			// todo check buf[:-1] == `}`
			//if d.requestData.Json[len(d.requestData.Json)-1] != '}' {
			//	return errors.New("invalid json syntax")
			//}
			//d.requestData.Json = strings.TrimSpace(d.requestData.Json[1 : len(d.requestData.Json)-1])
			err := d.getObject(mErr)
			if err != nil {
				return err
			}
			return nil
		//case `[`:
		//	err := d.getArray(json, data)
		//	if err != nil {
		//		return err
		//	}
		//	return nil
		default:
			return fmt.Errorf("invalid imput")
		}
	}
	return nil
}

//func (data *data) getDataValue(data any) (reflect.Value, reflect.Type) {
//	if val, ok := data.(reflect.Value); ok {
//		typ := val.Elem().Type()
//		return val, typ
//	}
//	val := reflect.ValueOf(data)
//	typ := val.Elem().Type()
//	return val, typ
//
//}

func (d *data) getObject(mErr *MultiError) error {

	var depth int
	inString := false
	expectComma := false
	var start int

	for i, v := range d.requestData.Json {
		z := string(v)
		_ = z
		switch string(v) {
		case "{", "[":
			if inString == true {
				continue
			}
			expectComma = true
			depth++
		case "}", "]":
			if inString == true {
				continue
			}
			if expectComma == true {
				return errors.New("invalid json syntax")
			}
			depth--
		case "\"":
			inString = !inString
		case ":":
			continue
			// todo check in string
		case ",":
			expectComma = false
			if depth == 1 && !inString {
				//d.requestData.Json = d.requestData.Json[i+1:]
				err := d.setFieldJSON(strings.TrimSpace(d.requestData.Json[start:i+1]), mErr)
				if err != nil {
					return err
				}
				start = i + 1
			} else {
				continue
			}
		case " ", "\r", "\n":
			continue
		default:
			if inString != true {
				return fmt.Errorf("invalid json syntax")
			}
			continue
		}
	}
	if start != len(d.requestData.Json) {
		err := d.setFieldJSON(strings.TrimSpace(d.requestData.Json[start:]), mErr)
		if err != nil {
			return err
		}
	}
	return nil
}

var i int

func (d *data) setFieldJSON(keyValue string, mErr *MultiError) error {
	trimmedKV := strings.TrimSpace(keyValue)
	fmt.Println(i, trimmedKV)
	i++

	//if key == "" || value == "" {
	//	// todo
	//	return errors.New("invalid ")
	//}
	//
	//for index := range d.typ.NumField() {
	//	field, err := d.newField(index)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if field.tags.Name == key {
	//		if field.typ.Type.Kind() == reflect.Struct {
	//			newData := &data{
	//				val:         field.val,
	//				typ:         field.typ.Type,
	//				tagType:     d.tagType,
	//				requestData: d.requestData.Json,
	//			}
	//			fmt.Println(newData.requestData.Json)
	//
	//			err := newData.decodeJson(mErr)
	//			if err != nil {
	//				return err
	//			}
	//		}
	//
	//	}
	//
	//	// todo return error field not found, bad request
	//}
	return nil
}

//func (parser *jsonParser) getKeyValue() []string {
//	KeyValueList := make([]string, 0)
//	var curToken strings.Builder
//	var depth int
//	inString := false
//
//	for _, v := range json {
//		switch string(v) {
//		case "{", "[":
//			depth++
//			curToken.WriteRune(v)
//		case "}", "]":
//			depth--
//			curToken.WriteRune(v)
//		case "\"":
//			inString = !inString
//			curToken.WriteRune(v)
//		case ":":
//			curToken.WriteRune(v)
//		case ",":
//			if depth == 0 && !inString {
//				str := strings.TrimSpace(curToken.String()) // todo refactor
//				KeyValueList = append(KeyValueList, str)
//				curToken.Reset()
//			} else {
//				curToken.WriteRune(v)
//			}
//		default:
//			curToken.WriteRune(v)
//		}
//	}
//	if curToken.Len() > 0 {
//		str := strings.TrimSpace(curToken.String()) // todo refactor
//		KeyValueList = append(KeyValueList, str)
//	}
//	return KeyValueList
//}

//func getField(typ reflect.Type, val reflect.Value, fieldName string) (reflect.Type, reflect.Value, bool) {
//	for i := range typ.NumField() { // todo get by name
//		typeField := typ.Field(i)
//		if strings.ToLower(typeField.Name) == strings.ToLower(fieldName) { // ломается тут
//			valField := val.Elem().Field(i)
//			return typeField.Type, valField, true
//		}
//	}
//	return nil, reflect.Value{}, false
//}
//
//func (parser *jsonParser) getKey(json string) string {
//	keyValue := strings.SplitN(json, ":", 2)
//	key := strings.TrimSpace(keyValue[0])
//	key = strings.Trim(key, `"`)
//	return key
//}

//func (parser *jsonParser) getValue(json string) string {
//	keyValue := strings.SplitN(json, ":", 2)
//	value := strings.TrimSpace(keyValue[1])
//	value = strings.Trim(value, `"`)
//	return value
//}

//func (parser *jsonParser) getArray(json string, data any) error {
//	json = strings.Trim(json, "[]")
//	json = strings.TrimSpace(json)
//
//	//val, typ := getDataValue(data)
//
//	for _, j := range json {
//		switch string(j) {
//		case "{":
//			kv := d.getKeyValue(json)
//			for _, v := range kv {
//				if val, ok := data.(reflect.Value); ok {
//					arrayValueTyp := val.Type().Elem().Elem() // get slice elem type
//					newData := reflect.New(arrayValueTyp)
//					// todo multi error
//
//					err := d.decodeJson(v, newData)
//					if err != nil {
//						return err
//					}
//					val.Elem().Set(reflect.Append(val.Elem(), newData.Elem()))
//
//				}
//
//			}
//			return nil
//		default:
//			val, ok := data.(reflect.Value)
//			if !ok {
//				return fmt.Errorf("QWE")
//			}
//			inString := false
//			var input strings.Builder
//			for _, v := range json {
//				switch string(v) {
//				case "\"":
//					inString = !inString
//				case ",":
//					if !inString {
//						str := strings.TrimSpace(input.String())
//
//						nv := reflect.New(val.Type().Elem().Elem()).Elem()
//						err := setFieldData(nv, str, "idn")
//						fmt.Println(nv)
//						if err != nil {
//							// todo multi error
//						}
//						input.Reset()
//
//						val.Elem().Set(reflect.Append(val.Elem(), nv))
//					}
//				default:
//					input.WriteRune(v)
//				}
//			}
//			if input.Len() > 0 {
//				str := strings.TrimSpace(input.String()) // todo refactor
//
//				nv := reflect.New(val.Type().Elem().Elem())
//				err := setFieldData(nv.Elem(), str, "idn")
//				if err != nil {
//					// todo multi error
//				}
//				reflect.Append(val.Elem(), nv.Elem())
//			}
//			return nil
//		}
//	}
//	return nil
//}
