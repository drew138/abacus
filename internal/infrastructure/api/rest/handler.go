package rest

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

//go:embed web
var web embed.FS

const (
	sessionCookie             = "abacus_session"
	statusClientClosedRequest = 499
	categoryOther             = "Other"
	groupFood                 = "Food and eating out"
)

var (
	errMalformedBody  = errors.New("malformed body")
	errMalformedForm  = errors.New("malformed form")
	errBadImportToken = errors.New("bad import token")
	errCrossSite      = errors.New("cross-site request")
	errInvalidMonth   = errors.New("invalid month")
	errInvalidID      = errors.New("invalid movement id")
)

var enums = map[string]func(string) bool{
	"method": func(v string) bool { return domain.Method(v).Valid() },
}

var categoryColors = map[string]string{
	"Restaurants and cafés":             "#e0763a",
	"Food delivery":                     "#f0a35e",
	"Groceries":                         "#c9a227",
	"Snacks and candy":                  "#b7791f",
	"Transfers to people":               "#d9534f",
	"Rides and transport":               "#7c5cdb",
	"Gym and fitness":                   "#1f9d74",
	"Phone, internet and subscriptions": "#2b6fd6",
	"Online shopping":                   "#d14d8b",
	"Pharmacy and health":               "#0e9fb0",
	"Gaming":                            "#5865c9",
	"Entertainment":                     "#e87ba4",
	"Travel":                            "#3d8fb0",
	categoryOther:                       "#8a909c",
	domain.CategorySavings:              "#15803d",
	groupFood:                           "#e0763a",
}

var spareColors = []string{"#a0622d", "#b5179e", "#6f8f00", "#3d5a99", "#d98e04", "#008300", "#e34948", "#0096b4"}

var methodLabels = map[domain.Method]string{
	domain.MethodCredit: "Credit Card",
	domain.MethodDebit:  "Bank Transfer",
}

type Application interface {
	ListMonths(ctx context.Context) ([]domain.MonthSummary, error)
	ListMovements(ctx context.Context, month string) ([]domain.Movement, error)
	ListCategoryTotals(ctx context.Context) (map[string]map[string]int64, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
	GetMovement(ctx context.Context, id int64) (domain.Movement, error)
	CreateMovement(ctx context.Context, movement domain.Movement) (domain.Movement, error)
	UpdateMovement(ctx context.Context, movement domain.Movement, applyToMerchant bool) (domain.Movement, error)
	RecategoriseMerchant(ctx context.Context, merchant, category string) (int, error)
	DeleteMovement(ctx context.Context, id int64) error
	RestoreMovement(ctx context.Context, id int64) error
	RegisterUser(ctx context.Context, username, password string) (domain.User, error)
	CheckRegistration(ctx context.Context) (bool, error)
	LogInUser(ctx context.Context, username, password string) (string, time.Time, error)
	AuthenticateUser(ctx context.Context, token string) (domain.User, error)
	LogOutUser(ctx context.Context, token string) error
	ImportMovements(ctx context.Context, movements []domain.Movement, categories []domain.Category) (added, updated, kept int, err error)
}

type Handler struct {
	app          Application
	validate     *validator.Validate
	templates    *template.Template
	static       http.Handler
	importToken  string
	secureCookie bool
	zone         *time.Location
}

func New(router chi.Router, app Application, importToken string, secureCookie bool) (*Handler, error) {
	static, err := fs.Sub(web, "web/static")
	if err != nil {
		return nil, fmt.Errorf("opening static files: %w", err)
	}
	versions, err := assetVersions(static)
	if err != nil {
		return nil, fmt.Errorf("fingerprinting static files: %w", err)
	}
	templates, err := template.New("abacus").Funcs(template.FuncMap{"asset": func(name string) string { return "/static/" + name + "?v=" + versions[name] }}).ParseFS(web, "web/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}
	h := &Handler{
		app:          app,
		validate:     newValidator(),
		templates:    templates,
		static:       http.StripPrefix("/static/", http.FileServer(http.FS(static))),
		importToken:  importToken,
		secureCookie: secureCookie,
		zone:         time.FixedZone("America/Bogota", -5*60*60),
	}

	router.Get("/static/*", h.logged(h.getStaticFile))
	router.Get("/logo.svg", h.logged(h.getLogo))
	router.Get("/login", h.logged(h.getLoginPage))
	router.Post("/login", h.logged(h.sameOrigin(h.createSession)))
	router.Get("/register", h.logged(h.getRegisterPage))
	router.Post("/register", h.logged(h.sameOrigin(h.createUser)))
	router.Post("/logout", h.logged(h.sameOrigin(h.deleteSession)))

	router.Get("/", h.logged(h.authed(h.getIndex)))
	router.Get("/months/{month}", h.logged(h.authed(h.getMonth)))
	router.Get("/months/{month}/expenses", h.logged(h.authed(h.getExpenses)))
	router.Get("/trends", h.logged(h.authed(h.getTrends)))
	router.Get("/movements/create-form", h.logged(h.authed(h.getCreateMovementForm)))
	router.Get("/movements/{id}/update-form", h.logged(h.authed(h.getUpdateMovementForm)))
	router.Get("/movements/{id}/row", h.logged(h.authed(h.getMovementRow)))
	router.Get("/movements/{id}/delete", h.logged(h.authed(h.getMovementDeletion)))
	router.Get("/movements/{id}/categories", h.logged(h.authed(h.getCategoryMenu)))
	router.Post("/movements", h.logged(h.sameOrigin(h.authed(h.createMovement))))
	router.Put("/movements/{id}", h.logged(h.sameOrigin(h.authed(h.updateMovement))))
	router.Delete("/movements/{id}", h.logged(h.sameOrigin(h.authed(h.deleteMovement))))
	router.Post("/movements/{id}/restore", h.logged(h.sameOrigin(h.authed(h.restoreMovement))))
	router.Post("/merchant-rules", h.logged(h.sameOrigin(h.authed(h.recategoriseMerchant))))

	router.Post("/api/import", h.logged(h.importer(h.importMovements)))
	return h, nil
}

func assetVersions(static fs.FS) (map[string]string, error) {
	entries, err := fs.ReadDir(static, ".")
	if err != nil {
		return nil, err
	}
	versions := make(map[string]string, len(entries))
	for _, e := range entries {
		body, err := fs.ReadFile(static, e.Name())
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(body)
		versions[e.Name()] = hex.EncodeToString(sum[:6])
	}
	return versions, nil
}

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		tag := f.Tag.Get("json")
		if tag == "" {
			tag = f.Tag.Get("form")
		}
		name, _, _ := strings.Cut(tag, ",")
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

func (h *Handler) validateRequest(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		var wrongType *json.UnmarshalTypeError
		if errors.As(err, &wrongType) {
			return err
		}
		return fmt.Errorf("%w: %w", errMalformedBody, err)
	}
	return h.validate.Struct(dst)
}

func (h *Handler) parseForm(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("%w: %w", errMalformedForm, err)
	}
	return nil
}

func parseMonth(r *http.Request) (string, error) {
	month := chi.URLParam(r, "month")
	if _, err := time.Parse("2006-01", month); err != nil || len(month) != 7 {
		return "", errInvalidMonth
	}
	return month, nil
}

func parseID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errInvalidID
	}
	return id, nil
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *Handler) authed(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := h.app.AuthenticateUser(r.Context(), sessionToken(r))
		switch {
		case errors.Is(err, domain.ErrNotAuthenticated):
			h.toLogin(w, r)
		case err != nil:
			h.writeError(w, r, err)
		default:
			next(w, r)
		}
	}
}

func (h *Handler) toLogin(w http.ResponseWriter, r *http.Request) {
	if isHTMX(r) {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) sameOrigin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isSameOrigin(r) {
			h.writeError(w, r, errCrossSite)
			return
		}
		next(w, r)
	}
}

func isSameOrigin(r *http.Request) bool {
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" {
		return site == "same-origin"
	}
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	_, host, _ := strings.Cut(origin, "://")
	return host == r.Host
}

func (h *Handler) importer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || subtle.ConstantTimeCompare([]byte(got), []byte(h.importToken)) != 1 {
			h.writeAPIError(w, errBadImportToken)
			return
		}
		next(w, r)
	}
}

func sessionToken(r *http.Request) string {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return ""
	}
	return c.Value
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, value string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) today() time.Time {
	return time.Now().In(h.zone)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, status int, name string, data any) {
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, name, data); err != nil {
		h.writeError(w, r, fmt.Errorf("rendering %s: %w: %w", name, domain.ErrInternal, err))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func (h *Handler) showMonth(w http.ResponseWriter, month, message string, undoID int64) {
	location, _ := json.Marshal(map[string]string{"path": "/months/" + month, "target": "#main"})
	w.Header().Set("HX-Location", string(location))
	h.toast(w, message, undoID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) toast(w http.ResponseWriter, message string, undoID int64) {
	payload := map[string]any{"message": message}
	if undoID > 0 {
		payload["undo"] = undoID
	}
	body, _ := json.Marshal(map[string]any{"toast": payload})
	w.Header().Set("HX-Trigger", asciiJSON(body))
}

func asciiJSON(body []byte) string {
	var b strings.Builder
	for len(body) > 0 {
		r, size := utf8.DecodeRune(body)
		body = body[size:]
		switch {
		case r < 0x80:
			b.WriteRune(r)
		case r > 0xFFFF:
			r -= 0x10000
			fmt.Fprintf(&b, `\u%04x\u%04x`, 0xD800+(r>>10), 0xDC00+(r&0x3FF))
		default:
			fmt.Fprintf(&b, `\u%04x`, r)
		}
	}
	return b.String()
}

func categoryColor(name string) string {
	if c, ok := categoryColors[name]; ok {
		return c
	}
	var hash uint32
	for _, r := range name {
		hash = hash*31 + uint32(r)
	}
	return spareColors[hash%uint32(len(spareColors))]
}

func categoryTier(name, parent string) int {
	switch {
	case name == domain.CategorySavings:
		return 0
	case name == groupFood || parent == groupFood:
		return 1
	case name == categoryOther:
		return 3
	}
	return 2
}

func monthName(key string) string {
	t, _ := time.Parse("2006-01", key)
	return t.Format("January 2006")
}

func monthShort(key string) string {
	t, _ := time.Parse("2006-01", key)
	return t.Format("Jan")
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

func (h *Handler) record(w http.ResponseWriter, err error) bool {
	rr, recorded := w.(*responseRecorder)
	if errors.Is(err, context.Canceled) {
		if recorded {
			rr.status = statusClientClosedRequest
		}
		return false
	}
	if recorded {
		rr.err = err
	}
	return true
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	if !h.record(w, err) {
		return
	}
	status, message := errorPage(err)
	if isHTMX(r) {
		w.Header().Set("HX-Reswap", "none")
		h.toast(w, message, 0)
		w.WriteHeader(status)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "<!doctype html><title>Abacus</title><p>%s</p><p><a href=\"/\">Back to Abacus</a></p>", template.HTMLEscapeString(message))
}

func errorPage(err error) (int, string) {
	switch {
	case errors.Is(err, errMalformedForm):
		return http.StatusBadRequest, "That form could not be read. Reload the page and try again."
	case errors.Is(err, errCrossSite):
		return http.StatusForbidden, "That request did not come from Abacus. Reload the page and try again."
	case errors.Is(err, errInvalidMonth):
		return http.StatusBadRequest, "That month does not exist. Pick one from the sidebar."
	case errors.Is(err, errInvalidID):
		return http.StatusBadRequest, "That movement does not exist. Reload the month."
	case errors.Is(err, domain.ErrUnknownCategory):
		return http.StatusBadRequest, "Pick a category from the list."
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "That movement no longer exists. Reload the month."
	case errors.Is(err, domain.ErrInvalid):
		return http.StatusBadRequest, "Abacus could not accept that. Check the values and try again."
	case errors.Is(err, domain.ErrUnavailable), errors.Is(err, context.DeadlineExceeded):
		return http.StatusServiceUnavailable, "Abacus is busy right now. Try again in a moment."
	}
	return http.StatusInternalServerError, "Something went wrong. Try again, and reload the page if it keeps happening."
}

type ErrorResponse struct {
	Type   string        `json:"type"`
	Title  string        `json:"title"`
	Status int           `json:"status"`
	Errors []ErrorDetail `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Code  string `json:"code"`
	Field string `json:"field,omitempty"`
	Param string `json:"param,omitempty"`
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
	"required": "value_required",
	"oneof":    "invalid_option",
	"enum":     "invalid_option",
	"datetime": "invalid_date",
}

var lengthCodes = map[string]string{"max": "max_length", "lte": "max_length", "min": "min_length", "gte": "min_length", "len": "length_mismatch"}

var numberCodes = map[string]string{"max": "less_than_or_equal", "lte": "less_than_or_equal", "lt": "less_than", "min": "greater_than_or_equal", "gte": "greater_than_or_equal", "gt": "greater_than"}

func (h *Handler) writeAPIError(w http.ResponseWriter, err error) {
	if !h.record(w, err) {
		return
	}
	body := errorResponse(err)
	h.writeJSON(w, body.Status, body)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func errorResponse(err error) ErrorResponse {
	var wrongType *json.UnmarshalTypeError
	var failed validator.ValidationErrors
	switch {
	case errors.As(err, &wrongType):
		return ErrorResponse{Type: "invalid_argument", Title: "The request has invalid fields.", Status: http.StatusBadRequest,
			Errors: []ErrorDetail{{Code: "invalid_type", Field: wrongType.Field, Param: jsonType(wrongType.Type)}}}
	case errors.Is(err, errMalformedBody):
		return ErrorResponse{Type: "malformed_request", Title: "The request body is not valid JSON.", Status: http.StatusBadRequest}
	case errors.Is(err, errBadImportToken):
		return ErrorResponse{Type: "unauthenticated", Title: "A valid import token is required.", Status: http.StatusUnauthorized}
	case errors.As(err, &failed):
		return ErrorResponse{Type: "invalid_argument", Title: "The request has invalid fields.", Status: http.StatusBadRequest, Errors: validationDetails(failed)}
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

func validationDetails(failed validator.ValidationErrors) []ErrorDetail {
	details := make([]ErrorDetail, 0, len(failed))
	for _, fe := range failed {
		field := fe.Field()
		if _, path, ok := strings.Cut(fe.Namespace(), "."); ok {
			field = path
		}
		code, param := "invalid_value", ""
		switch fe.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			if c, ok := lengthCodes[fe.Tag()]; ok {
				code, param = c, fe.Param()
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			if c, ok := numberCodes[fe.Tag()]; ok {
				code, param = c, fe.Param()
			}
		}
		if code == "invalid_value" {
			if c, ok := tagCodes[fe.Tag()]; ok {
				code = c
			}
		}
		details = append(details, ErrorDetail{Code: code, Field: field, Param: param})
	}
	return details
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
