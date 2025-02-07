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
	"slices"
	"strings"
)

// todo refactor comments

var (
	FileType  = reflect.TypeOf((*os.File)(nil))
	FilesType = reflect.TypeOf(([]*os.File)(nil))
)

const (
	PDF  = "application/pdf"
	EPUB = "application/epub+zip"
)

// SupportedContentType is a slice of supported file format
type SupportedContentType []string

// FormParser is a struct for parsing request body [*http.Request] in multipart/form-data format
type FormParser struct {
	MaxBodySize       int
	supportFileFormat SupportedContentType

	data           *Data
	dataFieldNames []string
}

func (parser *FormParser) HasField(key string) bool {
	fmt.Println(key)
	fmt.Println(parser.dataFieldNames)
	if !slices.Contains(parser.dataFieldNames, strings.ToLower(key)) {
		return false
	}
	return true
}

// FormParse is a function for parsing request body [*http.Request] in multipart/form-data format to pointer struct.
// Return error - ErrUnknownContentType, ErrFieldLength, ErrContentToLarge, ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge,  MultiError.
func FormParse(r *http.Request, data any) error {
	err := requestValidate(r)
	if err != nil {
		return err
	}

	d, err := dataCreate(data)
	if err != nil {
		return err
	}

	parser := &FormParser{
		supportFileFormat: SupportedContentType{EPUB, PDF},
		data:              d,
		dataFieldNames:    d.fieldNames(),
	}

	// todo http data to [Data]
	httpData, err := parser.HTTPBodyParse(r)
	if err != nil {
		return err
	}

	defer httpData.Files.RemoveFiles()

	//d, err := dataCreate(data)
	//if err != nil {
	//	return err
	//}
	//
	//d.typ.FieldByName("")

	// todo important error
	var mErr = &MultiError{}
	err = d.setDataValue(mErr, httpData)
	if err != nil {
		return err
	}

	if mErr.err != nil {
		return mErr
	}

	return nil
}

// requestValidate validate request content type and return error - ErrUnknownContentType
func requestValidate(r *http.Request) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return ErrUnknownContentType
	}
	return nil
}

// setDataValue writes data to fields of a data structure.
// Return error - MultiError.
func (data *Data) setDataValue(mErr *MultiError, httpData *HttpData) error {
	for fieldNum := range data.typ.NumField() {
		fieldData, err := fieldData(data, fieldNum, FORM)
		if err != nil {
			return err
		}
		//fmt.Println(fieldData.fieldType.Name)

		//httpData.Values.Get(fieldData.fieldTagName)
		switch fieldData.fieldType.Type {
		case FileType:
			defer httpData.Files.SetNil(mErr, fieldData.fieldTagName)

			file := httpData.Files[fieldData.fieldTagName]

			if fieldData.fieldTag == RequiredTag && len(file) < 1 || len(file) > 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single file; ", fieldData.fieldTagName))
				continue
			}
			if fieldData.fieldTag == OptionalTag && len(file) > 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single file; ", fieldData.fieldTagName))
				continue
			}

			if file != nil {
				fieldData.fieldValue.Set(reflect.ValueOf(file[0]))
			}
		case FilesType:
			defer httpData.Files.SetNil(mErr, fieldData.fieldTagName)

			files := httpData.Files[fieldData.fieldTagName]
			if fieldData.fieldTag == RequiredTag && len(files) < 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single or more files; ", fieldData.fieldTagName))
				continue
			}
			if files != nil {
				fieldData.fieldValue.Set(reflect.ValueOf(files))
			}
		default:
			err := fieldData.setDefaultValue(httpData, mErr)
			if err != nil {
				mErr.err = append(mErr.err, err)
				continue
			}
		}
	}
	return nil
}

func (fieldData *FieldData) setDefaultValue(httpData *HttpData, mErr *MultiError) error {
	if fieldData == nil {
		log.Fatalf("request: fieldData cannot be nil")
	}
	if fieldData.fieldValue.Kind() == reflect.Struct {
		// todo add tags there?
		d := &Data{
			val: fieldData.fieldValue,
			typ: fieldData.fieldType.Type,
		}

		err := d.setDataValue(mErr, httpData)
		if err != nil {
			return err
		}
		fieldData.fieldValue.Set(d.val)
		return nil
	}

	ok := httpData.Values.Has(fieldData.fieldTagName)
	if !ok && fieldData.fieldTag == RequiredTag {
		return ErrFieldRequired{field: fieldData.fieldTagName}
	}

	if fieldData.fieldTag == OptionalTag && httpData.Values.Get(fieldData.fieldTagName) == "" {
		return nil
	}

	err := setField(fieldData.fieldValue, httpData.Values.Get(fieldData.fieldTagName), fieldData.fieldType.Name)
	if err != nil {
		if fieldData.fieldTag == OptionalTag || fieldData.fieldTag == RequiredTag {
			return err
		}
		return nil
	}
	return nil
}

// HTTPBodyParse parsing request body [*http.Request] in multipart/form-data format, form value write in
// [*HttpData.Values] and files write in [*HttpData.Files].
// Return error - ErrFieldLength, ErrContentToLarge, ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge
func (parser *FormParser) HTTPBodyParse(r *http.Request) (*HttpData, error) {
	var result = &HttpData{
		Values: make(url.Values),
		Files:  make(Files),
	}

	// todo add custom reader
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if part.FileName() == "" {
			key, value, err := parser.parseFormValue(part)
			if err != nil {
				return nil, err
			}
			result.Values.Add(key, value)
		} else {
			// todo что если произойдет ошибка то удалится только файл у которого была ошибка, а остальные без ошибки не удалятся
			fileName, file, err := parser.parseFormFile(part)
			if err != nil {
				//result.Files.RemoveFiles()
				return nil, err
			}
			result.Files.Add(fileName, file)
		}
	}
	return result, nil
}

// todo if content-type != file type

// parseFormFile reads a file from a multipart.Part and writes it to a temporary file.
// Return error - ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge.
// Log file remove error.
func (parser *FormParser) parseFormFile(part *multipart.Part) (fileName string, file *os.File, err error) {
	// todo test this
	//if _, ok := parser.data.typ.FieldByName(part.FormName()); !ok {
	//	return "", nil, ErrFieldName
	//}
	if !parser.HasField(part.FormName()) {
		return "", nil, fmt.Errorf("request: invalid field name '%s'", part.FormName())
	}

	err = parser.validateFilePart(part)
	if err != nil {
		return "", nil, err
	}

	tempFile, err := os.CreateTemp("", "upload-*_"+part.FileName())
	if err != nil {
		return "", nil, err
	}

	defer func() {
		if err != nil {
			err := os.Remove(tempFile.Name())
			if err != nil {
				log.Printf("request: parseFile: file remove error: %s", err)
			}
		}
	}()
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			log.Printf("request: parseFile: file close error: %s", err)
			return
		}
	}(tempFile)

	err = parser.readFile(part, tempFile)
	if err != nil {
		return "", nil, err
	}

	return strings.ToLower(part.FormName()), tempFile, nil
}

// validateFilePart validate multipart.Part - content type(only EPUB and PDF), form name, file name
// Return error - ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong.
func (parser *FormParser) validateFilePart(part *multipart.Part) error {
	contentType := part.Header.Get("Content-Type")
	// todo refactor
	ok := false
	for _, format := range parser.supportFileFormat {
		if format == contentType {
			ok = true
			break
		}
	}
	if !ok {
		return ErrInvalidContentType
	}

	if !strings.Contains(strings.ToLower(part.FormName()), "file") {
		return ErrInvalidFieldType
	}
	if len(part.FileName()) > 100 {
		return ErrFileNameTooLong
	}
	return nil
}

// readFile reads a file from a multipart.Part and writes it to a temporary file.
// Return error - ErrContentToLarge.
func (parser *FormParser) readFile(part *multipart.Part, file *os.File) error {
	buf := make([]byte, 4096)
	for {
		n, readErr := part.Read(buf)
		if n > 0 {
			if _, writeErr := file.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to temp file: %w", writeErr)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			var maxBytesError *http.MaxBytesError
			if errors.As(readErr, &maxBytesError) {
				return ErrContentToLarge
			}
			return fmt.Errorf("failed to read part: %w", readErr)
		}
	}
	return nil
}

// parseFormValue reads multipart.Part and returns formValueName and Value.
// Return error - ErrFieldLength, ErrContentToLarge.
func (parser *FormParser) parseFormValue(part *multipart.Part) (key string, value string, err error) {
	// todo test this
	if !parser.HasField(part.FormName()) {
		return "", "", fmt.Errorf("request: invalid field name '%s'", part.FormName())
	}

	buf := make([]byte, 512)
	n, err := part.Read(buf)
	if err != nil {
		if err == io.EOF && n >= 0 {
			return strings.ToLower(part.FormName()), string(buf[:n]), nil
		}
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return "", "", ErrContentToLarge
		}
		return "", "", err
	}
	return "", "", ErrFieldLength
}
