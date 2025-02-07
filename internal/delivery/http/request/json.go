package request

import (
	"encoding/json"
	"errors"
	"net/http"
)

// todo validate data IMPORTNT
func JsonParse(r *http.Request, data any) error {
	// todo if error is EOF
	_, err := dataCreate(data)
	if err != nil {
		return err
	}

	err = json.NewDecoder(r.Body).Decode(data)
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
