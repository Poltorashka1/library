package request

import (
	"net/http"
	"reflect"
)

// 'required' - обязательное поле, должно быть значение
// '' - не обязательное поле, если придут невалидные данные ошибки не будет
// optional - не обязательное поле со значением по умолчанию, если придут невалидные данные будет ошибка

// todo поменять местами ” и 'optional'
// todo if in query in UPPER case
func QueryParse(r *http.Request, data any) error {
	val, err := dataValidate(data)
	if err != nil {
		return err
	}
	err = setQueryDataValue(r, val)
	if err != nil {
		return err
	}

	return nil
}

func setQueryDataValue(r *http.Request, val reflect.Value) error {
	var mErr = &MultiError{}
	typ := val.Type()

	for i := range typ.NumField() {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)

		fieldName, fieldTag, err := getFieldTags(&fieldType, QUERY)
		if err != nil {
			return err
		}

		queryValue := r.URL.Query().Get(fieldName)

		if fieldTag == "required" {
			if queryValue == "" {
				mErr.err = append(mErr.err, ErrQueryRequired{queryKey: fieldName})
				continue
			}
		}

		if queryValue == "" {
			continue
		}

		err = setFieldData(fieldValue, queryValue, fieldName)
		if err != nil {
			if fieldTag == "optional" {
				mErr.err = append(mErr.err, err)
				continue
			}
			continue
		}
	}
	if mErr.err != nil {
		return mErr
	}
	return nil
}
