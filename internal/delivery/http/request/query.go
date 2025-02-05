package request

import (
	"net/http"
	"net/url"
	"reflect"
)

// 'required' - обязательное поле, должно быть значение
// '' - не обязательное поле, если придут невалидные данные ошибки не будет
// optional - не обязательное поле со значением по умолчанию, если придут невалидные данные будет ошибка

// todo поменять местами ” и 'optional'

type QueryParser struct {
	val  reflect.Value
	typ  reflect.Type
	data any

	Values url.Values
}

func QueryParse(r *http.Request, data any) error {
	val, err := dataValidate(data)
	if err != nil {
		return err
	}

	parser := QueryParser{
		val:    val,
		typ:    val.Type(),
		Values: make(url.Values),
	}

	parser.QueryParse(r)

	err = parser.setQueryDataValue()
	if err != nil {
		return err
	}

	return nil
}

func (p *QueryParser) QueryParse(r *http.Request) {
	p.Values = r.URL.Query()
}

func (p *QueryParser) setQueryDataValue() error {
	var mErr = &MultiError{}

	for i := range p.typ.NumField() {
		fieldValue := p.val.Field(i)
		fieldType := p.typ.Field(i)

		fieldName, fieldTag, err := getFieldTags(&fieldType, QUERY)
		if err != nil {
			return err
		}

		queryValue := p.Values.Get(fieldName)

		if fieldTag == RequiredTag {
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
			if fieldTag == OptionalTag {
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
