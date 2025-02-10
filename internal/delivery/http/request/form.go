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
// todo add несколько режимов работы библиотеки.
// todo add payload validation?
// todo обобщить ошибки

var (
	FileType  = reflect.TypeOf((*os.File)(nil))
	FilesType = reflect.TypeOf(([]*os.File)(nil))
)

const (
	PDF  = "application/pdf"
	EPUB = "application/epub+zip"
)

// parserConfig is a config params for formParser
type parserConfig struct {
	// max request body size
	maxBodySize int
	// max formValueField size
	maxFormValueSize int
	// max one file size
	maxFileSize int
	// supported file format
	supportFileFormat supportedFileType
}

// SupportedFileType is a slice of supported file format
type supportedFileType []string

// FormParser is a struct for parsing request body [*http.Request] in multipart/form-data format
type formParser struct {
	cfg  parserConfig
	data *Data
	// todo new type there
	dataFieldTags []string
}

// HasField is a function for checking if the request form field exists in source struct, if not return false.
func (parser *formParser) HasField(key string) bool {
	if !slices.Contains(parser.dataFieldTags, strings.ToLower(key)) {
		return false
	}
	return true
}

// FormParse is a function for parsing request body in multipart/form-data format into pointer struct.
// Return error - ErrUnknownContentType, ErrFieldLength, ErrContentToLarge, ErrInvalidFileType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge,  MultiError.
func FormParse(r *http.Request, payload any) error {
	err := formRequestValidate(r)
	if err != nil {
		return err
	}

	d, err := data(payload, FORM)
	if err != nil {
		return err
	}

	tags, err := d.fieldTags()
	if err != nil {
		return err
	}

	parser := &formParser{
		cfg: parserConfig{
			maxFormValueSize:  512,
			maxBodySize:       104857600,
			maxFileSize:       31457280,
			supportFileFormat: supportedFileType{EPUB, PDF},
		},
		data:          d,
		dataFieldTags: tags,
	}

	// todo http data to [Data] and Field
	httpData, err := parser.HTTPBodyParse(r)
	if err != nil {
		return err
	}

	var mErr = &MultiError{}
	defer func() {
		if err != nil || mErr.err != nil {
			httpData.Files.RemoveFiles()
		}
	}()

	tags = nil
	parser = nil

	err = d.setDataValue(mErr, httpData)
	if err != nil {
		return err
	}

	if mErr.err != nil {
		return mErr
	}

	return nil
}

// formRequestValidate validate request content type.
// Return error - ErrUnknownContentType.
func formRequestValidate(r *http.Request) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return ErrUnknownContentType
	}
	return nil
}

// HTTPBodyParse parsing request body in multipart/form-data format, form value write in
// HttpData.Values and files write in HttpData.Files if return error automatically remove all files.
// Return error - ErrContentToLarge, ErrFieldLength, ErrContentToLarge, ErrInvalidFileType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge
func (parser *formParser) HTTPBodyParse(r *http.Request) (httpData *HttpData, err error) {
	// todo test
	if parser.cfg.maxBodySize != 0 {
		if r.ContentLength > int64(parser.cfg.maxBodySize) {
			return nil, &ErrContentToLarge{limit: parser.cfg.maxBodySize}
		}
		r.Body = http.MaxBytesReader(nil, r.Body, int64(parser.cfg.maxBodySize))
	}

	var result = &HttpData{
		Values: make(url.Values),
		Files:  make(Files),
	}

	defer func() {
		if err != nil {
			result.Files.RemoveFiles()
		}
	}()

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
				//result.Files.RemoveFiles()
				return nil, err
			}
			result.Values.Add(key, value)
		} else {
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

// parseFormValue reads multipart.Part and returns formFieldName and FieldValue.
// Return error - ErrFieldName, ErrContentToLarge, ErrFieldLength.
func (parser *formParser) parseFormValue(part *multipart.Part) (formFieldName string, formFieldValue string, err error) {
	// todo test this
	// todo add custom error in maxFormValue with field name
	if !parser.HasField(part.FormName()) {
		return "", "", &ErrFieldName{fieldName: part.FormName()}
	}

	var result []byte
	var readSize int
	buf := make([]byte, 4096)
	for {
		n, readErr := part.Read(buf)
		if readErr != nil && readErr != io.EOF {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				return "", "", &ErrContentToLarge{limit: parser.cfg.maxBodySize}
			}
			return "", "", fmt.Errorf("request: readFile: failed to read part: %w", readErr)
		}
		if n > 0 {
			if parser.cfg.maxFormValueSize != 0 {
				readSize += n
				if readSize > parser.cfg.maxFormValueSize {
					return "", "", &ErrFormValueToLarge{formField: part.FormName(), limit: parser.cfg.maxFormValueSize}
				}
			}
			// todo optimize
			result = append(result, buf[:n]...)
		}

		if readErr == io.EOF {
			break
		}
	}
	return part.FormName(), string(result), nil
}

// todo if content-type != file type
// parseFormFile reads a file from a multipart.Part and writes it to a temporary file.
// Return error - ErrInvalidFileType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge.
// Log file remove error.
func (parser *formParser) parseFormFile(part *multipart.Part) (fileName string, file *os.File, err error) {
	if !parser.HasField(part.FormName()) {
		return "", nil, &ErrFieldName{fieldName: part.FormName()}
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
// Return error - ErrInvalidFileType, ErrInvalidFieldType, ErrFileNameTooLong.
func (parser *formParser) validateFilePart(part *multipart.Part) error {
	if !slices.Contains(parser.cfg.supportFileFormat, part.Header.Get("Content-Type")) {
		return &ErrInvalidFileType{strings.Join(parser.cfg.supportFileFormat, ", ")}
	}

	if len(part.FileName()) > 100 {
		return ErrFileNameTooLong
	}
	return nil
}

// todo add support max one file size
// readFile reads a file from a multipart.Part and writes it to a temporary file.
// Return error - ErrContentToLarge.
func (parser *formParser) readFile(part *multipart.Part, file *os.File) error {
	buf := make([]byte, 4096)
	var readBytes int
	for {
		n, readErr := part.Read(buf)
		if readErr != nil && readErr != io.EOF {
			var maxBytesError *http.MaxBytesError
			if errors.As(readErr, &maxBytesError) {
				return &ErrContentToLarge{limit: parser.cfg.maxBodySize}
			}
			return fmt.Errorf("request: readFile: failed to read part: %w", readErr)
		}

		if n > 0 {
			if parser.cfg.maxFileSize != 0 {
				readBytes += n
				if readBytes > parser.cfg.maxFileSize {
					return &ErrFormValueToLarge{formField: part.FormName(), limit: parser.cfg.maxFileSize}
				}
			}
			if _, writeErr := file.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to temp file: %w", writeErr)
			}
		}

		if readErr == io.EOF {
			break
		}
	}
	return nil
}

// setDataValue writes data to fields of a data structure.
// Return error - MultiError.
func (data *Data) setDataValue(mErr *MultiError, httpData *HttpData) error {
	for fieldNum := range data.typ.NumField() {
		field, err := field(data, fieldNum)
		if err != nil {
			return err
		}

		switch field.typ.Type {
		case FileType:
			err := field.setFileValue(httpData, mErr)
			if err != nil {
				return err
			}
		case FilesType:
			err := field.setFilesValue(httpData, mErr)
			if err != nil {
				mErr.Add(err)
				return err
			}
		default:
			err := field.setDefaultValue(httpData, mErr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (field *Field) setFileValue(httpData *HttpData, mErr *MultiError) error {
	file := httpData.Files[field.tags.Name]

	if field.tags.Tag == RequiredTag && len(file) < 1 || len(file) > 1 {
		mErr.Add(fmt.Errorf("the field '%s' must be a single file; ", field.tags.Name))
		return nil
	}
	if field.tags.Tag == OptionalTag && len(file) > 1 {
		mErr.Add(fmt.Errorf("the field '%s' must be a single file; ", field.tags.Name))
		return nil
	}

	if file != nil {
		field.val.Set(reflect.ValueOf(file[0]))
	}

	return nil
}

func (field *Field) setFilesValue(httpData *HttpData, mErr *MultiError) error {
	files := httpData.Files[field.tags.Name]
	if field.tags.Tag == RequiredTag && len(files) < 1 {
		mErr.Add(fmt.Errorf("the field '%s' must be a single or more files; ", field.tags.Name))
		return nil
	}

	if files != nil {
		field.val.Set(reflect.ValueOf(files))
	}
	return nil
}

func (field *Field) setDefaultValue(httpData *HttpData, mErr *MultiError) error {
	if field == nil {
		return errors.New("request: setDefaultValue: field cannot be nil")
	}
	if field.val.Kind() == reflect.Struct {
		d := &Data{
			val:     field.val,
			typ:     field.typ.Type,
			tagType: field.tagType,
		}

		err := d.setDataValue(mErr, httpData)
		if err != nil {
			return err
		}
		field.val.Set(d.val)
		return nil
	}

	ok := httpData.Values.Has(field.tags.Name)
	if !ok && field.tags.Tag == RequiredTag {
		mErr.Add(ErrFieldRequired{field: field.tags.Name})
		return nil
	}

	if field.tags.Tag == OptionalTag && httpData.Values.Get(field.tags.Name) == "" {
		return nil
	}

	err := field.setField(httpData.Values.Get(field.tags.Name))
	if err != nil {
		if field.tags.Tag == OptionalTag || field.tags.Tag == RequiredTag {
			mErr.Add(err)
			return nil
		}
		return nil
	}
	return nil
}
