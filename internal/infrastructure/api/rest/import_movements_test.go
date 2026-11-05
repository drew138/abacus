package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drew138/abacus/internal/domain"
)

func TestHandler_importMovements(t *testing.T) {
	t.Parallel()

	valid := `{"categories":[{"name":"Board games"}],"movements":[{"ref":"r1","date":"2026-10-03","merchant":"Rappi","amount":42900,"category":"Food delivery","method":"credit"}]}`
	tests := []struct {
		name       string
		auth       string
		body       string
		setupMocks func(*MockApplication)
		wantStatus int
		wantBody   string
	}{
		{
			name: "imports the movements",
			auth: "Bearer " + testToken,
			body: valid,
			setupMocks: func(app *MockApplication) {
				app.On("ImportMovements", anyCtx,
					[]domain.Movement{{Ref: "r1", Date: time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit}},
					[]domain.Category{{Name: "Board games"}}).
					Return(domain.ImportResult{Added: 1}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"added":1,"updated":0,"kept":0}`,
		},
		{
			name:       "rejects a missing token",
			body:       valid,
			setupMocks: func(*MockApplication) {},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"type":"unauthenticated","title":"A valid import token is required.","status":401}`,
		},
		{
			name:       "rejects a wrong token",
			auth:       "Bearer nope",
			body:       valid,
			setupMocks: func(*MockApplication) {},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"type":"unauthenticated","title":"A valid import token is required.","status":401}`,
		},
		{
			name:       "rejects invalid json",
			auth:       "Bearer " + testToken,
			body:       `{"movements":`,
			setupMocks: func(*MockApplication) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"type":"malformed_request","title":"The request body is not valid JSON.","status":400}`,
		},
		{
			name:       "rejects a wrong json type",
			auth:       "Bearer " + testToken,
			body:       `{"movements":[{"amount":"many"}]}`,
			setupMocks: func(*MockApplication) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"type":"invalid_argument","title":"The request has invalid fields.","status":400,"errors":[{"code":"invalid_type","field":"movements.0.amount","param":"integer"}]}`,
		},
		{
			name:       "reports every invalid field, and refuses an unknown method",
			auth:       "Bearer " + testToken,
			body:       `{"movements":[{"ref":"","date":"2026-13-01","merchant":"","amount":0,"category":"x","method":"cash"}]}`,
			setupMocks: func(*MockApplication) {},
			wantStatus: http.StatusBadRequest,
			wantBody: `{"type":"invalid_argument","title":"The request has invalid fields.","status":400,"errors":[` +
				`{"code":"value_required","field":"movements[0].ref"},{"code":"invalid_date","field":"movements[0].date"},` +
				`{"code":"value_required","field":"movements[0].merchant"},{"code":"value_required","field":"movements[0].amount"},` +
				`{"code":"invalid_option","field":"movements[0].method"}]}`,
		},
		{
			name:       "rejects too many characters",
			auth:       "Bearer " + testToken,
			body:       `{"movements":[{"ref":"r","date":"2026-10-01","merchant":"` + strings.Repeat("a", 121) + `","amount":1,"category":"x","method":"debit"}]}`,
			setupMocks: func(*MockApplication) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"type":"invalid_argument","title":"The request has invalid fields.","status":400,"errors":[{"code":"max_length","field":"movements[0].merchant","param":"120"}]}`,
		},
		{
			name:       "rejects a body over the limit",
			auth:       "Bearer " + testToken,
			body:       `{"movements":[` + strings.Repeat(" ", 2048) + `]}`,
			setupMocks: func(*MockApplication) {},
			wantStatus: http.StatusRequestEntityTooLarge,
			wantBody:   `{"type":"payload_too_large","title":"The request body is too large.","status":413}`,
		},
		{
			name: "answers 503 when the store is down",
			auth: "Bearer " + testToken,
			body: valid,
			setupMocks: func(app *MockApplication) {
				app.On("ImportMovements", anyCtx,
					[]domain.Movement{{Ref: "r1", Date: time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit}},
					[]domain.Category{{Name: "Board games"}}).
					Return(domain.ImportResult{}, fmt.Errorf("importing: %w", domain.ErrUnavailable))
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   `{"type":"service_unavailable","title":"The service is temporarily unavailable.","status":503}`,
		},
		{
			name: "answers 500 for anything unexpected",
			auth: "Bearer " + testToken,
			body: valid,
			setupMocks: func(app *MockApplication) {
				app.On("ImportMovements", anyCtx,
					[]domain.Movement{{Ref: "r1", Date: time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit}},
					[]domain.Category{{Name: "Board games"}}).
					Return(domain.ImportResult{}, fmt.Errorf("importing: %w", domain.ErrInternal))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"type":"internal_error","title":"Something went wrong.","status":500}`,
		},
		{
			name: "answers 400 for an unknown category",
			auth: "Bearer " + testToken,
			body: valid,
			setupMocks: func(app *MockApplication) {
				app.On("ImportMovements", anyCtx,
					[]domain.Movement{{Ref: "r1", Date: time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC), Merchant: "Rappi", Amount: 42900, Category: "Food delivery", Method: domain.MethodCredit}},
					[]domain.Category{{Name: "Board games"}}).
					Return(domain.ImportResult{}, domain.ErrUnknownCategory)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"type":"invalid_argument","title":"The request is invalid.","status":400,"errors":[{"code":"unknown_category","field":"category"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := newMockApplication(t)
			tt.setupMocks(app)
			router := chi.NewRouter()
			New(router, app, testToken, 1024)

			req := httptest.NewRequest(http.MethodPost, "/api/import", strings.NewReader(tt.body))
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			require.True(t, json.Valid(rec.Body.Bytes()))
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
