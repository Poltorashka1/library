package response

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	StatusCode() int
	Content() any
	// todo mb delete content type?
	ContentType() string
	Headers() http.Header
}

type errorResponse struct {
	ResponseError       string `json:"error"`
	ResponseStatusCode  int    `json:"status_code"`
	responseContentType string
}

func (r *errorResponse) Content() any {
	return r.ResponseError
}
func (r *errorResponse) StatusCode() int {
	return r.ResponseStatusCode
}
func (r *errorResponse) ContentType() string {
	return r.responseContentType
}
func (r *errorResponse) Headers() http.Header {
	return nil
}

// Error sends a JSON response with the error text and status code
func Error(w http.ResponseWriter, err error, statusCode int) {
	response := &errorResponse{
		ResponseError:       err.Error(),
		ResponseStatusCode:  statusCode,
		responseContentType: "application/json", // default value
	}

	Write(w, response)
}

type successResponse struct {
	ResponseContent     any `json:"content"`
	ResponseStatusCode  int `json:"status_code"`
	responseContentType string
	responseHeaders     http.Header
}

// todo check how it work now, and headers to http.Header

// Success send response with content on PDF/JSON format (default JSON) and status code
func Success(w http.ResponseWriter, content any, contentType string, headers ...map[string]string) {
	var httpHeaders = http.Header{}
	if headers != nil {
		for k, v := range headers[0] {
			httpHeaders.Set(k, v)
		}
	}
	response := &successResponse{
		ResponseStatusCode:  http.StatusOK,
		ResponseContent:     content,
		responseContentType: contentType,
		responseHeaders:     httpHeaders,
	}

	Write(w, response)
}

func (r *successResponse) Headers() http.Header {
	return r.responseHeaders
}

//func File(w http.ResponseWriter, content any, contentType string) {
//
//}

func (r *successResponse) Content() any {
	return r.ResponseContent
}
func (r *successResponse) StatusCode() int {
	return r.ResponseStatusCode
}
func (r *successResponse) ContentType() string {
	return r.responseContentType
}

func Write(w http.ResponseWriter, response Response) {
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
	// todo check if pdf and content not file format
	case response.ContentType() == "application/pdf":
		WritePDF(w, response)
	case response.ContentType() == "application/json":
		WriteJSON(w, response)
	case response.ContentType() == "text/html":
		WriteHTML(w, response)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}
}

// todo delete this

func WriteHTML(w http.ResponseWriter, response Response) {
	_, err := w.Write(response.Content().([]byte))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}
}

func WritePDF(w http.ResponseWriter, response Response) {
	_, err := w.Write(response.Content().([]byte))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}
}

func WriteJSON(w http.ResponseWriter, response Response) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Server error"))
	}
}

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}
