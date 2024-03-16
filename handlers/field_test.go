package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_field_handlers(t *testing.T) {
	var tests = []test{
		{"/api/v1/fields", "GET", "", "", getAllFields, 200},
		{"/api/v1/fields/ff0b708b-c2dd-46db-ae96-1e9cddc485c7", "GET", "", "", getFieldById, 200},
	}

	for _, tst := range tests {
		var req *http.Request

		if tst.json == "" {
			req, _ = http.NewRequest(tst.method, tst.name, nil)
		} else {
			req, _ = http.NewRequest(tst.method, tst.name, strings.NewReader(tst.json))
		}

		if tst.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("user_id", tst.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tst.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != tst.expectedStatus {
			t.Errorf("%s: wrong status returned; expected %d but got %d", tst.name, tst.expectedStatus, rr.Code)
		}

	}

}
