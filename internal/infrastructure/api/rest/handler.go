package rest

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

const statusClientClosedRequest = 499

var (
	errMalformedBody  = errors.New("malformed body")
	errBadImportToken = errors.New("bad import token")
)

type Application interface {
	ImportMovements(ctx context.Context, movements []domain.Movement, categories []domain.Category) (domain.ImportResult, error)
}

type Handler struct {
	app          Application
	validate     *validator.Validate
	importToken  string
	maxBodyBytes int64
}

func New(router chi.Router, app Application, importToken string, maxBodyBytes int64) *Handler {
	h := &Handler{app: app, validate: newValidator(), importToken: importToken, maxBodyBytes: maxBodyBytes}
	router.Post("/api/import", h.logged(h.importer(h.importMovements)))
	return h
}

var enums = map[string]func(string) bool{
	"method": func(v string) bool { return domain.Method(v).Valid() },
}

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
	_ = v.RegisterValidation("enum", func(fl validator.FieldLevel) bool {
		valid, ok := enums[fl.Param()]
		return ok && valid(fl.Field().String())
	})
	return v
}

func (h *Handler) validateRequest(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, h.maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		var tooLarge *http.MaxBytesError
		var wrongType *json.UnmarshalTypeError
		if errors.As(err, &tooLarge) || errors.As(err, &wrongType) {
			return err
		}
		return fmt.Errorf("%w: %w", errMalformedBody, err)
	}
	return h.validate.Struct(dst)
}

func (h *Handler) importer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || subtle.ConstantTimeCompare([]byte(got), []byte(h.importToken)) != 1 {
			h.writeError(w, errBadImportToken)
			return
		}
		next(w, r)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	err    error
}

func (rr *responseRecorder) WriteHeader(status int) {
	rr.status = status
	rr.ResponseWriter.WriteHeader(status)
}

func (h *Handler) logged(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rr, r)

		attrs := []any{"method", r.Method, "route", chi.RouteContext(r.Context()).RoutePattern(), "status", rr.status, "duration_ms", time.Since(start).Milliseconds()}
		switch {
		case rr.status == statusClientClosedRequest:
			slog.InfoContext(r.Context(), "request cancelled", attrs...)
			return
		case rr.status >= http.StatusInternalServerError:
			slog.ErrorContext(r.Context(), "request failed", append(attrs, "error", rr.err)...)
			return
		}
		slog.InfoContext(r.Context(), "request handled", attrs...)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	rr, recorded := w.(*responseRecorder)
	if errors.Is(err, context.Canceled) {
		if recorded {
			rr.status = statusClientClosedRequest
		}
		return
	}
	if recorded {
		rr.err = err
	}
	body := mapErrorResponse(err)
	h.writeJSON(w, body.Status, body)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func jsonType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	}
	return "string"
}
