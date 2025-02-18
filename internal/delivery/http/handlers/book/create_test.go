package bookhandlers

//
//import (
//	"book/internal/config"
//	loggerMocks "book/internal/logger/mocks"
//	"bytes"
//	"net/http"
//	"net/http/httptest"
//	"net/url"
//	"testing"
//)
//
//func TestBookHandlers_CreateBook(t *testing.T) {
//	type args struct {
//		r *http.Request
//	}
//
//	testTable := []struct {
//		name string
//		//fields
//		args     *args
//		wantErr  bool
//		wantCode int
//	}{
//		{
//			name: "success",
//			args: func() *args {
//				formData := url.Values{}
//				formData.Add("title", "TheGreatBook")
//				formData.Add("isbn", "9781234567890")
//				formData.Add("publication-year", "2025")
//				formData.Add("description", "A great book about history.")
//				formData.Add("publisher", "Famous Publisher")
//
//				req := httptest.NewRequest(http.MethodPost, "/add/", bytes.NewBufferString(formData.Encode()))
//				req.Header.Set("Content-Type", "multipart/form-data")
//				return &args{req}
//			}(),
//			wantErr:  false,
//			wantCode: http.StatusOK,
//		},
//	}
//	for _, tt := range testTable {
//		t.Run(tt.name, func(t *testing.T) {
//			useCase := useCaseMocks.NewUseCase(t)
//			logger := loggerMocks.NewLogger(t)
//			//cfg := configMocks.NewHandlersConfig(t)
//			cfg := config.NewHandlerConfig()
//			s := NewBookHandlers(logger, useCase, cfg)
//			rr := httptest.NewRecorder()
//			s.CreateBook(rr, tt.args.r)
//
//			if rr.Code != tt.wantCode {
//				t.Errorf("CreateBook() got = %v, want %v", rr.Code, tt.wantCode)
//			}
//			if tt.wantErr {
//				// Проверка, что тело ответа содержит сообщение об ошибке
//				if rr.Body.String() == "" {
//					t.Error("expected error message, got empty body")
//				}
//			} else {
//				// Проверка успешного ответа
//				if rr.Body.String() == "" {
//					t.Error("expected response body, got empty body")
//				}
//			}
//		})
//	}
//}
