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
	"strconv"
	"strings"
)

// todo add несколько режимов работы библиотеки.
// todo mb поменять местами '' и 'optional'
// todo add payload validation?
// todo обобщить ошибки
// todo возможность ограничивать поле не по размеру, а по количеству символов

// todo if field not file а приходит файл в поле.
const (
	pdf  = "application/pdf"
	epub = "application/epub+zip"
)

// todo add validation
// parserConfig is a config params for formParser
type formParserConfig struct {
	// max request body size
	maxBodySize int64
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
	cfg            formParserConfig
	data           *data
	dataFieldsTags map[string]string
}

// HasField is a function for checking if the request form field exists in source struct
func (parser *formParser) HasField(key string) bool {
	if key == "" {
		return false
	}

	if _, ok := parser.dataFieldsTags[key]; !ok {
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

	d, err := newData(payload, "form")
	if err != nil {
		return err
	}

	tags, err := d.fieldTags()
	if err != nil {
		return err
	}

	parser := &formParser{
		cfg: formParserConfig{
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
			d.request.Files.RemoveFiles()
		}
	}()

	//tags = nil
	//parser = nil

	err = d.setValue(mErr)
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
		if r.ContentLength > parser.cfg.maxBodySize {
			return &ErrContentToLarge{limit: parser.cfg.maxBodySize}
		}
		r.Body = http.MaxBytesReader(nil, r.Body, parser.cfg.maxBodySize)
	}
	parser.data.request = &requestData{
		Values: make(url.Values),
		Files:  make(files),
		Arrays: make([]array, 0),
	}

	defer func() {
		if err != nil {
			parser.data.request.Files.RemoveFiles()
		}
	}()

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
			err := parser.parseFormValue(part)
			if err != nil {
				//result.Files.RemoveFiles()
				return err
			}
		} else { // parse file
			err := parser.parseFormFile(part)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// parseFormValue reads multipart.Part and returns formFieldName and formFieldValue.
// Return error - ErrFieldName, ErrContentToLarge, ErrFormValueToLarge.
func (parser *formParser) parseFormValue(part *multipart.Part) error {
	tagKey := part.FormName()

	if strings.HasSuffix(tagKey, "]") {
		idx := strings.Index(tagKey, "[")
		if idx == -1 {
			return &ErrFieldName{fieldName: part.FormName()}
		}
		tagKey = tagKey[:idx]
	}
	if !parser.HasField(tagKey) {
		return &ErrFieldName{fieldName: part.FormName()}
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
				return &ErrContentToLarge{limit: parser.cfg.maxBodySize}
			}
			return fmt.Errorf("request: readFile: failed to read part: %w", readErr)
		}
		if n > 0 {
			if parser.cfg.maxFormValueSize > 0 {
				readSize += n
				if readSize > parser.cfg.maxFormValueSize {
					return &ErrFormValueToLarge{formField: part.FormName(), limit: parser.cfg.maxFormValueSize}
				}
			}
			_, writeErr := result.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("request: parseFormValue: failed to write to buffer: %w", writeErr)
			}
		}

		if readErr == io.EOF {
			break
		}
	}

	trimResult := strings.TrimSpace(result.String())
	trimResult = replaceMultipleSpaces(trimResult)

	// parse form Value if it arrays field
	// todo
	//if tagKey != part.FormName() {
	//	err := parser.parseArrayValue(tagKey, trimResult, part)
	//	if err != nil {
	//		return err
	//	}
	//}

	parser.data.request.Values.Add(part.FormName(), trimResult)
	return nil
}

func replaceMultipleSpaces(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	prevSpace := false

	for _, r := range s {
		if r == ' ' {
			if !prevSpace {
				result.WriteRune(r)
			}
			prevSpace = true
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}

	return result.String()
}

func (parser *formParser) parseArrayValue(tagKey string, formValue string, part *multipart.Part) error {
	var arr array

	// проверка есть ли данные для массива по ключу
	for _, v := range parser.data.request.Arrays {
		if v.fieldName == tagKey {
			arr = v
			break
		}
	}

	// если массив не создан, создаем
	if arr.structArray == nil {
		arr.fieldName = tagKey
		arr.structArray = make(map[int]map[string]string)
		parser.data.request.Arrays = append(parser.data.request.Arrays, arr)
	}

	res := strings.Builder{}
	depth := 0
	index := -1
	// парсинг ключей поля
	for _, v := range part.FormName()[len(tagKey):] {
		switch v {
		case '[':
			depth++
		case ']':
			if depth == 0 {
				return &ErrFieldName{fieldName: part.FormName()}
			}
			if index == -1 {
				newIndex, err := strconv.Atoi(res.String())
				if err != nil {
					return &ErrFieldName{fieldName: part.FormName()}
				}
				// проверяем есть ли данные для массива по индексу, если нет то инициализируем новую хэш таблицу
				if _, ok := arr.structArray[newIndex]; !ok {
					arr.structArray[newIndex] = make(map[string]string)
				}

				index = newIndex
				res.Reset()
				depth--
				continue
			}
			arr.structArray[index][res.String()] = formValue
		default:
			if depth <= 0 {
				return &ErrFieldName{fieldName: part.FormName()}
			}
			res.WriteRune(v)
		}
	}
	return nil
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
func (parser *formParser) parseFormFile(part *multipart.Part) (err error) {
	if !parser.HasField(part.FormName()) {
		return &ErrFieldName{fieldName: part.FormName()}
	}

	err = parser.validateFilePart(part)
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp("", "upload-*_"+part.FileName())
	if err != nil {
		return err
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
		return err
	}

	parser.data.request.Files.Add(strings.ToLower(part.FormName()), tempFile)
	return nil
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
