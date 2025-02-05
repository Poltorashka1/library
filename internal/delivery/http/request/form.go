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
	val  reflect.Value
	typ  reflect.Type
	data any

	Values url.Values
	Files  Files

	supportFileFormat SupportedContentType // []string of supported file format in request
}

// todo test new feature - tag required optional and default

// FormParse is a function for parsing request body [*http.Request] in multipart/form-data format to pointer struct.
// Return error - ErrUnknownContentType, ErrFieldLength, ErrContentToLarge, ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge,  MultiError.
func FormParse(r *http.Request, data any) error {
	err := requestValidate(r)
	if err != nil {
		return err
	}

	val, err := dataValidate(data)
	if err != nil {
		return err
	}

	parser := &FormParser{
		val:               val,
		typ:               val.Type(),
		data:              data,
		supportFileFormat: SupportedContentType{EPUB, PDF},
		Values:            make(url.Values),
		Files:             make(Files),
	}
	defer parser.Files.RemoveFiles()

	err = parser.HTTPBodyParse(r)
	if err != nil {
		return err
	}

	err = parser.setDataValue()
	if err != nil {
		return err
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
func (p *FormParser) setDataValue() error {
	var mErr = &MultiError{}
	for fieldNum := range p.typ.NumField() {

		fieldType := p.typ.Field(fieldNum)
		fieldValue := p.val.Field(fieldNum)

		fieldTagName, tag, err := getFieldTags(&fieldType, FORM)
		if err != nil {
			return err
		}

		switch fieldType.Type {
		case reflect.TypeOf((*os.File)(nil)):
			defer p.Files.SetNil(mErr, fieldTagName)

			file := p.Files[fieldTagName]

			if tag == "required" && len(file) < 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single file; ", fieldTagName))
				continue
			}
			if tag == "optional" && len(file) > 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single file; ", fieldTagName))
				continue
			}

			if file != nil {
				fieldValue.Set(reflect.ValueOf(file[0]))
			}
		case reflect.TypeOf(([]*os.File)(nil)):
			defer p.Files.SetNil(mErr, fieldTagName)

			files := p.Files[fieldTagName]
			if tag == "required" && len(files) < 1 {
				mErr.err = append(mErr.err, fmt.Errorf("the field '%s' must be a single or more files; ", fieldTagName))
				continue
			}
			if files != nil {
				fieldValue.Set(reflect.ValueOf(files))
			}
		default:
			ok := p.Values.Has(fieldTagName)
			if !ok && tag == "required" {
				mErr.err = append(mErr.err, ErrFieldRequired{field: fieldTagName})
				continue
			}

			if tag == "optional" && p.Values.Get(fieldTagName) == "" {
				continue
			}

			err = setFieldData(fieldValue, p.Values.Get(fieldTagName), fieldType.Name)
			if err != nil {
				if tag == "optional" {
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
func (p *FormParser) HTTPBodyParse(r *http.Request) error {
	// todo add custom reader
	reader, err := r.MultipartReader()
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
			key, value, err := p.parseFormValue(part)
			if err != nil {
				return err
			}
			p.Values.Add(key, value)
		} else {
			fileName, file, err := p.parseFormFile(part)
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

// parseFormFile reads a file from a multipart.Part and writes it to a temporary file.
// Return error - ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong, ErrContentToLarge.
// Log file remove error.
func (p *FormParser) parseFormFile(part *multipart.Part) (fileName string, file *os.File, err error) {
	err = p.validateFilePart(part)
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

	err = readFile(part, tempFile)
	if err != nil {
		return "", nil, err
	}

	return strings.ToLower(part.FormName()), tempFile, nil
}

// validateFilePart validate multipart.Part - content type(only EPUB and PDF), form name, file name
// Return error - ErrInvalidContentType, ErrInvalidFieldType, ErrFileNameTooLong.
func (p *FormParser) validateFilePart(part *multipart.Part) error {
	contentType := part.Header.Get("Content-Type")
	for _, format := range p.supportFileFormat {
		if format == contentType {
			break
		}
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
func readFile(part *multipart.Part, file *os.File) error {
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
func (p *FormParser) parseFormValue(part *multipart.Part) (key string, value string, err error) {
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
