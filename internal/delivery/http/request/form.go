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

type formParser struct {
	val reflect.Value
	typ reflect.Type
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

	files, err := HTTPBodyParse(&r)
	if err != nil {
		return err
	}

	err = setFormDataValue(r, val, files)
	if err != nil {
		tempFilesRemove(files)
		return err
	}

	return nil
}

//if fieldType.Type == reflect.TypeOf((*os.File)(nil)) {
//	//defer func() {
//	//	if mErr.err == nil {
//	//		(*files)[fieldName] = nil
//	//	}
//	//}()
//	defer clearFiles(fieldTagName, mErr, files)
//
//	file := (*files)[fieldTagName]
//
//	if len(file) != 1 {
//		mErr.err = append(mErr.err, fmt.Errorf("the required field '%s' must be a single file; ", fieldTagName))
//		continue
//	}
//	fieldValue.Set(reflect.ValueOf(file[0]))
//
//	continue
//}
//if fieldType.Type == reflect.TypeOf(([]*os.File)(nil)) {
//	//defer func() {
//	//	if mErr.err == nil {
//	//		(*files)[fieldName] = nil
//	//	}
//	//}()
//	defer clearFiles(fieldTagName, mErr, files)
//
//	file := (*files)[fieldTagName]
//
//	if len(file) < 1 {
//		// todo or нету файла вообще
//		mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single or more files; ", fieldTagName))
//		continue
//	}
//	fieldValue.Set(reflect.ValueOf(file))
//
//	continue
//}

//ok := r.Form.Has(fieldTagName)
//if !ok {
//	mErr.err = append(mErr.err, ErrFieldRequired{field: fieldTagName})
//	continue
//}
//formValue := r.Form.Get(fieldTagName)
//
//err = setFieldData(fieldValue, formValue, fieldType.Name)
//if err != nil {
//	mErr.err = append(mErr.err, err)
//	continue
//}

func setFormDataValue(r *http.Request, val reflect.Value, files *map[string][]*os.File) error {

	typ := val.Type()
	var mErr = &MultiError{}
	for fieldNum := range typ.NumField() {

		fieldType := typ.Field(fieldNum)
		fieldValue := val.Field(fieldNum)

		//fieldTag := fieldType.Tag.Get("form")
		fieldTagName, _, err := getFieldTags(&fieldType, FORM)
		if err != nil {
			return err
		}

		switch fieldType.Type {
		case reflect.TypeOf((*os.File)(nil)):
			defer clear(fieldTagName, mErr, files)

			file := (*files)[fieldTagName]

			if len(file) != 1 {
				// delete file
				mErr.err = append(mErr.err, fmt.Errorf("the required field '%s' must be a single file; ", fieldTagName))
				continue
			}
			// do not delete file
			fieldValue.Set(reflect.ValueOf(file[0]))
		case reflect.TypeOf(([]*os.File)(nil)):
			defer clear(fieldTagName, mErr, files)

			file := (*files)[fieldTagName]

			if len(file) < 1 {
				// todo or нету файла вообще
				// delete file
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single or more files; ", fieldTagName))
				continue
			}
			// do not delete file
			fieldValue.Set(reflect.ValueOf(file))
		default:
			ok := r.Form.Has(fieldTagName)
			if !ok {
				mErr.err = append(mErr.err, ErrFieldRequired{field: fieldTagName})
				continue
			}

			err = setFieldData(fieldValue, r.Form.Get(fieldTagName), fieldType.Name)
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

func clear(fieldTagName string, mErr *MultiError, files *map[string][]*os.File) {
	if mErr.err == nil {
		(*files)[fieldTagName] = nil
	}
}

func HTTPBodyParse(r **http.Request) (*map[string][]*os.File, error) {
	reader, err := (*r).MultipartReader()
	if err != nil {
		return nil, err
	}

	(*r).Form = make(url.Values)
	var files = make(map[string][]*os.File)
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if part.FileName() == "" {
			key, value, err := parsePart(part)
			if err != nil {
				return nil, err
			}
			(*r).Form.Add(key, value)
		} else {
			file, fileField, err := parseFile(part)
			if err != nil {
				tempFilesRemove(&files)
				return nil, err
			}
			files[fileField] = append(files[fileField], file)
		}
	}
	return &files, nil
}

func tempFilesRemove(files *map[string][]*os.File) {
	//todo optimize and if files пришел пустой.
	for _, v := range *files {
		if v != nil {
			for i := range v {
				err := os.Remove(v[i].Name())
				if err != nil {
					log.Printf("request: file remove error: %s", err)
				}
			}
		}
	}
}

// todo if content-type != file type
func parseFile(part *multipart.Part) (file *os.File, fileName string, err error) {
	contentType := part.Header.Get("Content-Type")
	if contentType != EPUB && contentType != PDF {
		return nil, "", ErrInvalidContentType
	}
	if !strings.Contains(part.FormName(), "file") {
		return nil, "", ErrInvalidFieldType
	}
	if len(part.FileName()) > 100 {
		return nil, "", ErrFileNameTooLong
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
		return nil, "", err
	}

	buf := make([]byte, 4096)
	for {
		n, readErr := part.Read(buf)
		if n > 0 {
			if _, writeErr := tempFile.Write(buf[:n]); writeErr != nil {
				return nil, "", fmt.Errorf("failed to write to temp file: %w", writeErr)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			var maxBytesError *http.MaxBytesError
			if errors.As(readErr, &maxBytesError) {
				return nil, "", ErrContentToLarge
			}
			return nil, "", fmt.Errorf("failed to read part: %w", readErr)
		}
	}

	return tempFile, part.FormName(), nil
}

func parsePart(part *multipart.Part) (string, string, error) {
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
	//if err != nil && err != io.EOF {
	//	var maxBytesError *http.MaxBytesError
	//	if errors.As(err, &maxBytesError) {
	//		return "", "", ErrContentToLarge
	//	}
	//	return "", "", err
	//}
	return "", "", ErrFieldLength
}
