package response

import (
	"book/internal/errors"
	"encoding/json"
	"net/http"
	"reflect"
)

const (
	JSON string = "application/json"
)

type Response interface {
	StatusCode() int
	Content() any
	// todo mb delete content type?
	ContentType() string
	Headers() http.Header
}

type errorResponse struct {
	ErrorDetail errorDetail `json:"error"`
}

type errorDetail struct {
	ResponseStatusCode  int    `json:"code"`
	ResponseMessage     string `json:"message"`
	responseContentType string
}

func (r *errorResponse) Content() any {
	return r.ErrorDetail.ResponseMessage
}
func (r *errorResponse) StatusCode() int {
	return r.ErrorDetail.ResponseStatusCode
}
func (r *errorResponse) ContentType() string {
	return r.ErrorDetail.responseContentType
}
func (r *errorResponse) Headers() http.Header {
	return nil
}

// Error sends a JSON response with the error text and status code
func Error(w http.ResponseWriter, err error, statusCode int) {
	response := &errorResponse{errorDetail{ResponseMessage: err.Error(), ResponseStatusCode: statusCode, responseContentType: JSON}}

	write(w, response)
}

// ServerError sends a JSON response with the "server error" and StatusInternalServerError
func ServerError(w http.ResponseWriter) {
	// добавить логирование тогда в ней есть смысл
	// todo тестовая функция возможно стоит ее удалить так как не особо нужна или модифицировать
	response := &errorResponse{
		errorDetail{
			ResponseMessage:     apperrors.ErrServerError.Error(),
			ResponseStatusCode:  http.StatusInternalServerError,
			responseContentType: JSON,
		},
	}

	write(w, response)
}

type successDetail struct {
	ResponseContent     any `json:"content,omitempty"`
	ResponseStatusCode  int `json:"code"`
	responseContentType string
	responseHeaders     http.Header
}

type successResponse struct {
	SuccessDetail successDetail `json:"success"`
}

// todo c heck how it work now, and headers to http.Header

// Success send response with content on PDF/JSON format (default JSON) and status code
func Success(w http.ResponseWriter, content any, contentType string, headers ...map[string]string) {
	// todo add проверку на пустую структуру
	val := reflect.ValueOf(content)
	if val.Kind() == reflect.Pointer && val.IsNil() {
		content = nil
	}

	httpHeaders := http.Header{}
	if headers != nil {
		for k, v := range headers[0] {
			httpHeaders.Set(k, v)
		}
	}

	response := &successResponse{
		successDetail{
			ResponseContent:     content,
			ResponseStatusCode:  http.StatusOK,
			responseContentType: contentType,
			responseHeaders:     httpHeaders,
		},
	}

	write(w, response)
}

func (r *successResponse) Content() any {
	return r.SuccessDetail.ResponseContent
}
func (r *successResponse) StatusCode() int {
	return r.SuccessDetail.ResponseStatusCode
}
func (r *successResponse) ContentType() string {
	return r.SuccessDetail.responseContentType
}
func (r *successResponse) Headers() http.Header {
	return r.SuccessDetail.responseHeaders
}

func write(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", response.ContentType())
	for k, v := range response.Headers() {
		w.Header().Set(k, v[0])
	}
	w.WriteHeader(response.StatusCode())

	// todo test this
	//err := response.Headers().Write(w)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	_, _ = w.Write([]byte("Server error"))
	//}

	switch {
	// todo !!! unknown content type !!!
	// todo check if pdf and content not file format
	case response.ContentType() == "application/pdf":
		writePDF(w, response)
	case response.ContentType() == "application/json":
		writeJSON(w, response)
	case response.ContentType() == "text/html":
		writeHTML(w, response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}
}

// todo delete this

func writeHTML(w http.ResponseWriter, response Response) {
	_, err := w.Write(response.Content().([]byte))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}
}

func writePDF(w http.ResponseWriter, response Response) {
	_, err := w.Write(response.Content().([]byte))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}
}

func writeJSON(w http.ResponseWriter, response Response) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}
}

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}
