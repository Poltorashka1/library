package request

import (
	"net/http"
	"net/url"
)

// 'required' - обязательное поле, должно быть значение
// '' - не обязательное поле, если придут невалидные данные ошибки не будет
// optional - не обязательное поле со значением по умолчанию, если придут невалидные данные будет ошибка

// todo поменять местами ” и 'optional'

type QueryParser struct {
	Data   *Data
	Values url.Values
}

// todo refactor
func QueryParse(r *http.Request, data any) error {
	d, err := dataCreate(data)
	if err != nil {
		return err
	}

	parser := QueryParser{
		Data:   d,
		Values: make(url.Values),
	}

	parser.QueryParse(r)

	//err = parser.setQueryDataValue()
	//if err != nil {
	//	return err
	//}

	return nil
}

func (p *QueryParser) QueryParse(r *http.Request) {
	p.Values = r.URL.Query()
}

//func (p *QueryParser) setQueryDataValue() error {
//	var mErr = &MultiError{}
//
//	for i := range p.Data.typ.NumField() {
//		fieldValue := p.Data.val.Field(i)
//		fieldType := p.Data.typ.Field(i)
//
//		fieldTags, err := getFieldTags(&fieldType, QUERY)
//		if err != nil {
//			return err
//		}
//
//		queryValue := p.Values.Get(fieldTags.fieldTagName)
//
//		if fieldTags.fieldTag == RequiredTag {
//			if queryValue == "" {
//				mErr.err = append(mErr.err, ErrQueryRequired{queryKey: fieldTags.fieldTagName})
//				continue
//			}
//		}
//
//		if queryValue == "" {
//			continue
//		}
//
//		err = setField(fieldValue, queryValue, fieldTags.fieldTagName)
//		if err != nil {
//			if fieldTags.fieldTag == OptionalTag {
//				mErr.err = append(mErr.err, err)
//				continue
//			}
//			continue
//		}
//	}
//	if mErr.err != nil {
//		return mErr
//	}
//	return nil
//}
