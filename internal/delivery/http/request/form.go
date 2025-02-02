package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

type Files map[string][]*os.File

func (f Files) RemoveFiles() {
	for _, v := range f {
		if v != nil {
			for i := range v {
				err := os.Remove(v[i].Name())
				if err != nil {
					log.Printf("request: tempFilesRemove: file remove error: %s", err)
				}
			}
		}
	}
}

func (f Files) Add(fileName string, file *os.File) {
	f[fileName] = append(f[fileName], file)
}

func (f Files) SetNil(mErr *MultiError, fileName string) {
	if mErr.err == nil {
		f[fileName] = nil
	}

}

type FormParser struct {
	val  reflect.Value
	typ  reflect.Type
	data any

	Values url.Values
	Files  Files
}

func FormParse(r *http.Request, data any) error {
	// todo test this part
	//contentType := r.Header.Get("Content-Type")
	//fmt.Println(contentType)
	//if contentType != "application/x-www-form-urlencoded" && contentType != "multipart/form-data" {
	//	return ErrUnknownContentType
	//}

	val, err := dataValidate(data)
	if err != nil {
		return err
	}

	parser := &FormParser{
		val:    val,
		typ:    val.Type(),
		data:   data,
		Values: make(url.Values),
		Files:  make(Files),
	}
	defer parser.Files.RemoveFiles()

	err = parser.HTTPBodyParse(&r)
	if err != nil {
		return err
	}

	err = parser.setDataValue()
	if err != nil {
		return err
	}

	return nil
}

func (p *FormParser) setDataValue() error {
	var mErr = &MultiError{}
	for fieldNum := range p.typ.NumField() {

		fieldType := p.typ.Field(fieldNum)
		fieldValue := p.val.Field(fieldNum)

		fieldTagName, _, err := getFieldTags(&fieldType, FORM)
		if err != nil {
			return err
		}

		switch fieldType.Type {
		case reflect.TypeOf((*os.File)(nil)):
			defer p.Files.SetNil(mErr, fieldTagName)

			file := p.Files[fieldTagName]

			if len(file) != 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the required field '%s' must be a single file; ", fieldTagName))
				continue
			}
			fieldValue.Set(reflect.ValueOf(file[0]))
		case reflect.TypeOf(([]*os.File)(nil)):
			defer p.Files.SetNil(mErr, fieldTagName)

			files := p.Files[fieldTagName]
			if len(files) < 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single or more files; ", fieldTagName))
				continue
			}
			fieldValue.Set(reflect.ValueOf(files))
		default:
			ok := p.Values.Has(fieldTagName)
			if !ok {
				mErr.err = append(mErr.err, ErrFieldRequired{field: fieldTagName})
				continue
			}

			err = setFieldData(fieldValue, p.Values.Get(fieldTagName), fieldType.Name)
			if err != nil {
				mErr.err = append(mErr.err, err)
				continue
			}
		}
	}

	if mErr.err != nil {
		return mErr
	}
	return nil
}

func (p *FormParser) HTTPBodyParse(r **http.Request) error {
	reader, err := (*r).MultipartReader()
	if err != nil {
		return err
	}

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if part.FileName() == "" {
			key, value, err := p.parsePart(part)
			if err != nil {
				return err
			}
			p.Values.Add(key, value)
		} else {
			fileName, file, err := p.parseFile(part)
			if err != nil {
				//todo use clear func или вообще не стоит использовать тут
				return err
			}
			p.Files.Add(fileName, file)

		}
	}
	return nil
}

// todo if content-type != file type
func (p *FormParser) parseFile(part *multipart.Part) (fileName string, file *os.File, err error) {
	contentType := part.Header.Get("Content-Type")
	if contentType != EPUB && contentType != PDF {
		return "", nil, ErrInvalidContentType
	}
	if !strings.Contains(part.FormName(), "file") {
		return "", nil, ErrInvalidFieldType
	}
	if len(part.FileName()) > 100 {
		return "", nil, ErrFileNameTooLong
	}

	tempFile, err := os.CreateTemp("", "upload-*_"+part.FileName())
	defer func() {
		if err != nil {
			err := os.Remove(tempFile.Name())
			if err != nil {
				log.Println(err)
			}
		}
	}()
	defer tempFile.Close()

	if err != nil {
		return "", nil, err
	}

	buf := make([]byte, 4096)
	for {
		n, readErr := part.Read(buf)
		if n > 0 {
			if _, writeErr := tempFile.Write(buf[:n]); writeErr != nil {
				return "", nil, fmt.Errorf("failed to write to temp file: %w", writeErr)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			var maxBytesError *http.MaxBytesError
			if errors.As(readErr, &maxBytesError) {
				return "", nil, ErrContentToLarge
			}
			return "", nil, fmt.Errorf("failed to read part: %w", readErr)
		}
	}

	return part.FormName(), tempFile, nil
}

func (p *FormParser) parsePart(part *multipart.Part) (key string, value string, err error) {
	buf := make([]byte, 512)
	n, err := part.Read(buf)
	if err != nil {
		if err == io.EOF && n > 0 {
			return part.FormName(), string(buf[:n]), nil
		}
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return "", "", ErrContentToLarge
		}
		return "", "", err
	}
	return "", "", ErrFieldLength
}
