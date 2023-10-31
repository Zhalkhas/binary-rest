package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zhalkhas/binary-rest/internal/repos"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndicesController_getIndexHandler(t *testing.T) {
	tests := []struct {
		name           string
		searchFunc     searchFunc
		value          string
		expectedBody   string
		expectedStatus int
	}{
		{
			name: "non-integer value",
			searchFunc: func(val int64) (int, error) {
				return 0, nil
			},
			value:          "abc",
			expectedBody:   `{"message":"invalid value passed"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "negative integer value",
			searchFunc: func(val int64) (int, error) {
				return 0, nil
			},
			value:          "-1",
			expectedBody:   `{"message":"invalid value passed"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "from 0 to 1000000 if value is found it is returned",
			searchFunc: func(val int64) (int, error) {
				if val >= 0 && val <= 1000000 {
					return 0, nil
				}
				return -1, repos.ErrIndexNotFound
			},
			value:          "1000",
			expectedBody:   `{"index":0,"value":1000}`,
			expectedStatus: http.StatusOK,
		},
		{
			name: "value not found in file",
			searchFunc: func(val int64) (int, error) {
				return -1, repos.ErrIndexNotFound
			},
			value:          "1000",
			expectedBody:   `{"message":"index not found"}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "unexpected error returned during value search",
			searchFunc: func(val int64) (int, error) {
				return -1, fmt.Errorf("unexpected error")
			},
			value:          "1000",
			expectedBody:   `{"message":"unexpected error happened"}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewIndicesController(mockIndicesRepo{tt.searchFunc})

			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/endpoint/"+endpointValuePathParam, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add(endpointValuePathParam, tt.value)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			i.getIndexHandler(rw, req)

			if rw.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rw.Code)
			}

			if strings.TrimSpace(rw.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, rw.Body.String())
			}
		})
	}
}

type searchFunc func(val int64) (int, error)

type mockIndicesRepo struct {
	searchFunc
}

func (m mockIndicesRepo) Search(_ context.Context, val int64) (int, error) {
	return m.searchFunc(val)
}
