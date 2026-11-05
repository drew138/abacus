package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

func mapImportResponse(r domain.ImportResult) ImportResponse {
	return ImportResponse{Added: r.Added, Updated: r.Updated, Kept: r.Kept}
}

var domainErrors = []struct {
	err    error
	status int
	code   string
	field  string
}{
	{domain.ErrUnknownCategory, http.StatusBadRequest, "unknown_category", "category"},
}

var errorKinds = []struct {
	err    error
	status int
	typ    string
	title  string
}{
	{domain.ErrNotFound, http.StatusNotFound, "not_found", "The resource does not exist."},
	{domain.ErrConflict, http.StatusConflict, "resource_conflict", "The request conflicts with existing data."},
	{domain.ErrInvalid, http.StatusBadRequest, "invalid_argument", "The request is invalid."},
	{domain.ErrUnavailable, http.StatusServiceUnavailable, "service_unavailable", "The service is temporarily unavailable."},
	{context.DeadlineExceeded, http.StatusServiceUnavailable, "service_unavailable", "The service is temporarily unavailable."},
}

var tagCodes = map[string]string{
	"required":         "value_required",
	"required_if":      "value_required",
	"required_unless":  "value_required",
	"required_with":    "value_required",
	"required_without": "value_required",
	"oneof":            "invalid_option",
	"enum":             "invalid_option",
	"datetime":         "invalid_date",
}

var lengthCodes = map[string]string{
	"max": "max_length",
	"lte": "max_length",
	"min": "min_length",
	"gte": "min_length",
	"len": "length_mismatch",
}

var numberCodes = map[string]string{
	"max": "less_than_or_equal",
	"lte": "less_than_or_equal",
	"lt":  "less_than",
	"min": "greater_than_or_equal",
	"gte": "greater_than_or_equal",
	"gt":  "greater_than",
}

func mapErrorResponse(err error) ErrorResponse {
	var tooLarge *http.MaxBytesError
	var wrongType *json.UnmarshalTypeError
	var failed validator.ValidationErrors
	switch {
	case errors.As(err, &tooLarge):
		return ErrorResponse{Type: "payload_too_large", Title: "The request body is too large.", Status: http.StatusRequestEntityTooLarge}
	case errors.As(err, &wrongType):
		return mapInvalidParamResponse(wrongType.Field, jsonType(wrongType.Type))
	case errors.Is(err, errMalformedBody):
		return ErrorResponse{Type: "malformed_request", Title: "The request body is not valid JSON.", Status: http.StatusBadRequest}
	case errors.As(err, &failed):
		return ErrorResponse{Type: "invalid_argument", Title: "The request has invalid fields.", Status: http.StatusBadRequest, Errors: mapErrorDetails(failed)}
	case errors.Is(err, errBadImportToken):
		return ErrorResponse{Type: "unauthenticated", Title: "A valid import token is required.", Status: http.StatusUnauthorized}
	}

	var details []ErrorDetail
	status := 0
	for _, d := range domainErrors {
		if errors.Is(err, d.err) {
			if status == 0 {
				status = d.status
			}
			details = append(details, ErrorDetail{Code: d.code, Field: d.field})
		}
	}
	for _, k := range errorKinds {
		if errors.Is(err, k.err) {
			if status == 0 {
				status = k.status
			}
			return ErrorResponse{Type: k.typ, Title: k.title, Status: status, Errors: details}
		}
	}
	return ErrorResponse{Type: "internal_error", Title: "Something went wrong.", Status: http.StatusInternalServerError}
}

func mapInvalidParamResponse(field, expected string) ErrorResponse {
	return ErrorResponse{Type: "invalid_argument", Title: "The request has invalid fields.", Status: http.StatusBadRequest,
		Errors: []ErrorDetail{{Code: "invalid_type", Field: field, Param: expected}}}
}

func mapErrorDetails(failed validator.ValidationErrors) []ErrorDetail {
	details := make([]ErrorDetail, 0, len(failed))
	for _, fe := range failed {
		field := fe.Field()
		if _, path, ok := strings.Cut(fe.Namespace(), "."); ok {
			field = path
		}
		details = append(details, mapErrorDetail(fe, field))
	}
	return details
}

func mapErrorDetail(fe validator.FieldError, field string) ErrorDetail {
	switch fe.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		if code, ok := lengthCodes[fe.Tag()]; ok {
			return ErrorDetail{Code: code, Field: field, Param: fe.Param()}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		if code, ok := numberCodes[fe.Tag()]; ok {
			return ErrorDetail{Code: code, Field: field, Param: fe.Param()}
		}
	}
	if code, ok := tagCodes[fe.Tag()]; ok {
		return ErrorDetail{Code: code, Field: field}
	}
	return ErrorDetail{Code: "invalid_value", Field: field}
}
