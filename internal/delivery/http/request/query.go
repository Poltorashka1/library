package request

import (
	"net/http"
)

// 'required' - обязательное поле, должно быть значение
// '' - не обязательное поле, если придут невалидные данные ошибки не будет
// optional - не обязательное поле со значением по умолчанию, если придут невалидные данные будет ошибка

// todo поменять местами ” и 'optional'

type queryParser struct {
	data *data
}

func QueryParse(r *http.Request, payload any) error {
	d, err := newData(payload, query)
	if err != nil {
		return err
	}

	parser := queryParser{
		data: d,
	}

	parser.queryParse(r)

	var mErr = &MultiError{}
	err = d.setDataValue(mErr)
	if err != nil {
		return err
	}

	if mErr.err != nil {
		return mErr
	}
	return nil
}

func (parser *queryParser) queryParse(r *http.Request) {
	parser.data.requestData = &requestData{
		Values: r.URL.Query(),
	}
}

//func (parser *queryParser) setQueryDataValue() error {
//	var mErr = &MultiError{}
//
//	for i := range parser.data.typ.NumField() {
//		field, err := newField(parser.data, i)
//		if err != nil {
//			return err
//		}
//
//		//if field.typ.Type.Kind() == reflect.Struct {
//		//	nField := &data{
//		//		val:         field.val,
//		//		typ:         field.typ.Type,
//		//		tagType:     parser.data.tagType,
//		//		requestData: parser.data.requestData,
//		//	}
//		//
//		//	nField.setDataValue()
//		//}
//		queryValue := parser.data.requestData.Values.Get(field.tags.Name)
//
//		if field.tags.Tag == requiredTag && queryValue == "" {
//			mErr.Add(&ErrQueryRequired{queryKey: field.tags.Name})
//			return nil
//		}
//	}
//
//	//	if fieldTags.fieldTag == RequiredTag {
//	//		if queryValue == "" {
//	//			mErr.err = append(mErr.err, ErrQueryRequired{queryKey: fieldTags.fieldTagName})
//	//			continue
//	//		}
//	//	}
//	//
//	//	if queryValue == "" {
//	//		continue
//	//	}
//	//
//	//	err = setField(fieldValue, queryValue, fieldTags.fieldTagName)
//	//	if err != nil {
//	//		if fieldTags.fieldTag == OptionalTag {
//	//			mErr.err = append(mErr.err, err)
//	//			continue
//	//		}
//	//		continue
//	//	}
//	//}
//	//if mErr.err != nil {
//	//	return mErr
//	//}
//	//return nil
//}
