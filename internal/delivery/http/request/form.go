package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
)

// todo add несколько режимов работы библиотеки.
// todo поменять местами ” и 'optional'
// todo add payload validation?
// todo обобщить ошибки
// todo возможность ограничивать поле не по размеру, а по количеству символов

const (
	pdf  = "application/pdf"
	epub = "application/epub+zip"
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

// supportedFileType is a slice of supported file format
type supportedFileType []string

// formParser is a struct for parsing request body in multipart/form-data format
type formParser struct {
	cfg            parserConfig
	data           *data
	dataFieldsTags []string
}

// HasField is a function for checking if the request form field exists in source struct
func (parser *formParser) HasField(key string) bool {
	if !slices.Contains(parser.dataFieldsTags, strings.ToLower(key)) {
		return false
	}
	return true
}

// FormParse is a function for parsing request body in multipart/form-data format into pointer struct.
// Possible errors:
// ErrUnknownContentType;
// ErrContentToLarge;
// ErrFieldName;
// ErrFormValueToLarge;
// ErrInvalidFileType;
// ErrFileNameTooLong;
// MultiError;
func FormParse(r *http.Request, payload any) error {
	err := formRequestValidate(r)
	if err != nil {
		return err
	}

	d, err := newData(payload, form)
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
			supportFileFormat: supportedFileType{pdf, epub},
		},
		data:           d,
		dataFieldsTags: tags,
	}

	err = parser.httpBodyParse(r)
	if err != nil {
		return err
	}

	var mErr = &MultiError{}
	defer func() {
		if err != nil || mErr.err != nil {
			d.requestData.Files.RemoveFiles()
		}
	}()

	tags = nil
	parser = nil

	err = d.setDataValue(mErr)
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

// httpBodyParse parsing request body in multipart/form-data format, and write values in
// requestData.Values and files write in requestData.Files if return error automatically remove all files.
// Possible errors:
// - ErrContentToLarge
// - ErrFieldName
// - ErrFormValueToLarge
// - ErrFieldName
// - ErrInvalidFileType
// - ErrFileNameTooLong
// - ErrFormValueToLarge
func (parser *formParser) httpBodyParse(r *http.Request) (err error) {
	if parser.cfg.maxBodySize > 0 {
		if r.ContentLength > int64(parser.cfg.maxBodySize) {
			return &ErrContentToLarge{limit: parser.cfg.maxBodySize}
		}
		r.Body = http.MaxBytesReader(nil, r.Body, int64(parser.cfg.maxBodySize))
	}

	var result = &requestData{
		Values: make(url.Values),
		Files:  make(files),
	}

	defer func() {
		if err != nil {
			result.Files.RemoveFiles()
		}
	}()

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
			key, value, err := parser.parseFormValue(part)
			if err != nil {
				//result.Files.RemoveFiles()
				return err
			}
			if key != "" {
				result.Values.Add(key, value)
			}
		} else {
			fileName, file, err := parser.parseFormFile(part)
			if err != nil {
				//result.Files.RemoveFiles()
				return err
			}
			result.Files.Add(fileName, file)
		}
	}
	parser.data.requestData = result
	return nil
}

// parseFormValue reads multipart.Part and returns formFieldName and formFieldValue.
// Return error - ErrFieldName, ErrContentToLarge, ErrFormValueToLarge.
func (parser *formParser) parseFormValue(part *multipart.Part) (formFieldName string, formFieldValue string, err error) {
	if part.FormName() == "" {
		return "", "", errors.New("empty form field name")
	}

	if !parser.HasField(part.FormName()) {
		return "", "", &ErrFieldName{fieldName: part.FormName()}
	}

	var result bytes.Buffer
	if parser.cfg.maxFormValueSize > 0 {
		result.Grow(parser.cfg.maxFormValueSize)
	}

	var readSize int
	buf := make([]byte, 4096)
	for {
		n, readErr := part.Read(buf)
		if readErr != nil && readErr != io.EOF {
			var maxBytesError *http.MaxBytesError
			if errors.As(readErr, &maxBytesError) {
				return "", "", &ErrContentToLarge{limit: parser.cfg.maxBodySize}
			}
			return "", "", fmt.Errorf("request: readFile: failed to read part: %w", readErr)
		}
		if n > 0 {
			if parser.cfg.maxFormValueSize > 0 {
				readSize += n
				if readSize > parser.cfg.maxFormValueSize {
					return "", "", &ErrFormValueToLarge{formField: part.FormName(), limit: parser.cfg.maxFormValueSize}
				}
			}
			_, writeErr := result.Write(buf[:n])
			if writeErr != nil {
				return "", "", fmt.Errorf("request: parseFormValue: failed to write to buffer: %w", writeErr)
			}
		}

		if readErr == io.EOF {
			break
		}
	}
	if result.Len() <= 0 {
		return "", "", nil
	}
	return part.FormName(), result.String(), nil
}

// todo if content-type != file type
// parseFormFile reads a file from a multipart.Part and writes it to a temporary file.
// Possible errors:
// - ErrFieldName
// - ErrInvalidFileType
// - ErrFileNameTooLong
// - ErrContentToLarge
// - ErrFormValueToLarge
// Log file remove error.
func (parser *formParser) parseFormFile(part *multipart.Part) (fileName string, file *os.File, err error) {
	if part.FormName() == "" {
		return "", nil, errors.New("empty form field name")
	}

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

// validateFilePart validate multipart.Part - content type, form name, file name
// Return error - ErrInvalidFileType, ErrFileNameTooLong.
func (parser *formParser) validateFilePart(part *multipart.Part) error {
	if !slices.Contains(parser.cfg.supportFileFormat, part.Header.Get("Content-Type")) {
		return &ErrInvalidFileType{strings.Join(parser.cfg.supportFileFormat, ", ")}
	}

	if len(part.FileName()) > 100 {
		return ErrFileNameTooLong
	}
	return nil
}

// readFile reads a file from a multipart.Part and writes it to a temporary file.
// Return error - ErrContentToLarge, ErrFormValueToLarge.
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
