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

// Files is a map of file names to a slice of files
type Files map[string][]*os.File

// RemoveFiles removes any temporary files associated with a [Files].
// Log delete file errors.
func (f Files) RemoveFiles() {
	for _, v := range f {
		if v != nil {
			for i := range v {
				err := os.Remove(v[i].Name())
				if err != nil {
					log.Printf("request: RemoveFiles: file remove error: %s", err)
					return
				}
			}
		}
	}
}

// Add adds a file [*os.File] by key fileName to the [Files] map.
func (f Files) Add(fileName string, file *os.File) {
	f[fileName] = append(f[fileName], file)
}

// SetNil adds a nil file by key fileName to the [Files] map.
// To avoid deleting the necessary files, they must be deleted from [Files].
func (f Files) SetNil(mErr *MultiError, fileName string) {
	if mErr.err == nil {
		f[fileName] = nil
	}
}

// SupportedContentType is a slice of supported file format
type SupportedContentType []string

// FormParser is a struct for parsing form data
type FormParser struct {
	Data *Data

	Values url.Values
	Files  Files

	supportFileFormat SupportedContentType // []string of supported file format in request
}

type Parser struct {
	MaxFileSize       int
	supportFileFormat SupportedContentType
}

// FormParse is a function for parsing request body [*http.Request] in multipart/form-data format to pointer struct.
// Return error - ErrUnknownContentType, ErrFieldLength, ErrContentToLarge, ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge,  MultiError.
func FormParse(r *http.Request, data any) error {
	err := requestValidate(r)
	if err != nil {
		return err
	}

	parser := &Parser{
		supportFileFormat: SupportedContentType{EPUB, PDF},
	}

	httpData, err := parser.HTTPBodyParse(r)
	if err != nil {
		return err
	}

	defer httpData.Files.RemoveFiles()

	d, err := dataCreate(data, httpData)
	if err != nil {
		return err
	}

	err = d.setDataValue()
	if err != nil {
		return err
	}

	return nil
}

// requestValidate validate request content type and return error - ErrUnknownContentType
func requestValidate(r *http.Request) error {
	// todo add check request max size
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return ErrUnknownContentType
	}
	return nil
}

// setDataValue writes data to fields of a data structure.
// Return error - MultiError.
func (data *Data) setDataValue() error {
	var mErr = &MultiError{}
	for fieldNum := range data.typ.NumField() {

		fieldType := data.typ.Field(fieldNum)
		fieldValue := data.val.Field(fieldNum)

		// fieldTag doesn't work if it is a struct field, work only for struct fields
		fieldTagName, fieldTag, err := getFieldTags(&fieldType, FORM)
		if err != nil {
			return err
		}

		switch fieldType.Type {
		case reflect.TypeOf((*os.File)(nil)):
			defer data.httpData.Files.SetNil(mErr, fieldTagName)

			file := data.httpData.Files[fieldTagName]

			if fieldTag == RequiredTag && len(file) < 1 || len(file) > 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single file; ", fieldTagName))
				continue
			}
			if fieldTag == OptionalTag && len(file) > 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single file; ", fieldTagName))
				continue
			}

			if file != nil {
				fieldValue.Set(reflect.ValueOf(file[0]))
			}
		case reflect.TypeOf(([]*os.File)(nil)):
			defer data.httpData.Files.SetNil(mErr, fieldTagName)

			files := data.httpData.Files[fieldTagName]
			if fieldTag == RequiredTag && len(files) < 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single or more files; ", fieldTagName))
				continue
			}
			if files != nil {
				fieldValue.Set(reflect.ValueOf(files))
			}
		default:
			if fieldValue.Kind() == reflect.Struct {
				// todo add tags there?
				d := &Data{
					val:      fieldValue,
					typ:      fieldType.Type,
					httpData: data.httpData,
				}

				err = d.setDataValue()
				if err != nil {
					mErr.err = append(mErr.err, err)
					continue
				}
				fieldValue.Set(d.val)
				continue
			}

			ok := data.httpData.Values.Has(fieldTagName)
			if !ok && fieldTag == RequiredTag {
				mErr.err = append(mErr.err, ErrFieldRequired{field: fieldTagName})
				continue
			}

			if fieldTag == OptionalTag && data.httpData.Values.Get(fieldTagName) == "" {
				continue
			}

			err = setField(fieldValue, data.httpData.Values.Get(fieldTagName), fieldType.Name)
			if err != nil {
				if fieldTag == OptionalTag || fieldTag == RequiredTag {
					mErr.err = append(mErr.err, err)
					continue
				}
				continue
			}
		}
	}

	if mErr.err != nil {
		return mErr
	}
	return nil
}

// HTTPBodyParse parsing request body [*http.Request] in multipart/form-data format, form value write in
// [FormParser.Values] and files write in [FormParser.Files].
// Return error - ErrFieldLength, ErrContentToLarge, ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge
func (parser *Parser) HTTPBodyParse(r *http.Request) (*HttpData, error) {
	// todo add custom reader
	var result = &HttpData{
		Values: make(url.Values),
		Files:  make(Files),
	}

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
			fileName, file, err := parser.parseFormFile(part)
			if err != nil {
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
func (parser *Parser) parseFormFile(part *multipart.Part) (fileName string, file *os.File, err error) {
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
func (parser *Parser) validateFilePart(part *multipart.Part) error {
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
func (parser *Parser) readFile(part *multipart.Part, file *os.File) error {
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
func (parser *Parser) parseFormValue(part *multipart.Part) (key string, value string, err error) {
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
