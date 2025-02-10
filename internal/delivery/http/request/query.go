package request

import (
	"net/http"
)

// 'required' - обязательное поле, должно быть значение
// '' - не обязательное поле, если придут невалидные данные ошибки не будет
// optional - не обязательное поле со значением по умолчанию, если придут невалидные данные будет ошибка

// queryParser is a struct for parsing requestQuery
type queryParser struct {
	data *data
}

// QueryParse is a function for parsing request query into pointer struct.
// Possible errors:
// MultiError;
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

// queryParse parsing request query, and write values in requestData.Values
func (parser *queryParser) queryParse(r *http.Request) {
	parser.data.requestData = &requestData{
		Values: r.URL.Query(),
	}
}
