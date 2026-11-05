package htmx

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/drew138/abacus/internal/domain"
)

//go:embed web
var web embed.FS

const (
	sessionCookie             = "abacus_session"
	statusClientClosedRequest = 499
	maxFormBytes              = 64 << 10
)

var (
	errCrossSite      = errors.New("cross-site request")
	errInvalidMonth   = errors.New("invalid month")
	errInvalidID      = errors.New("invalid movement id")
	errMalformedForm  = errors.New("malformed form")
	enums             = map[string]func(string) bool{"method": func(v string) bool { return domain.Method(v).Valid() }}
	formFieldMessages = map[string]string{
		"merchant": "Enter who the money went to, up to 120 characters.",
		"amount":   "Enter the amount in pesos, for example 25.000.",
		"date":     "Pick the date it happened.",
		"category": "Pick a category from the list.",
		"method":   "Pick how it was paid.",
		"note":     "Keep the note under 200 characters.",
	}
)

type Application interface {
	ListMonths(ctx context.Context) ([]domain.MonthSummary, error)
	MonthReport(ctx context.Context, month string, filter domain.ExpenseFilter) (domain.MonthReport, error)
	Trends(ctx context.Context) (domain.Trends, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
	GetMovement(ctx context.Context, id int64) (domain.Movement, error)
	CreateMovement(ctx context.Context, movement domain.Movement) (domain.Movement, error)
	UpdateMovement(ctx context.Context, movement domain.Movement, applyToMerchant bool) (domain.Movement, error)
	RecategoriseMerchant(ctx context.Context, merchant, category string) (int, error)
	DeleteMovement(ctx context.Context, id int64) error
	RestoreMovement(ctx context.Context, id int64) error
	LogIn(ctx context.Context, username, password string) (string, domain.Session, error)
	Authenticate(ctx context.Context, token string) (domain.User, error)
	LogOut(ctx context.Context, token string) error
}

type Handler struct {
	app          Application
	validate     *validator.Validate
	templates    *template.Template
	static       http.Handler
	secureCookie bool
	now          func() time.Time
	zone         *time.Location
}

func New(router chi.Router, app Application, secureCookie bool) (*Handler, error) {
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
	h := &Handler{app: app, validate: newValidator(), templates: templates, static: http.StripPrefix("/static/", http.FileServer(http.FS(static))), secureCookie: secureCookie, now: time.Now, zone: time.FixedZone("America/Bogota", -5*60*60)}

	router.Get("/static/*", h.logged(h.getStatic))
	router.Get("/logo.svg", h.logged(h.getLogo))
	router.Get("/login", h.logged(h.getLogin))
	router.Post("/login", h.logged(h.sameOrigin(h.postLogin)))
	router.Post("/logout", h.logged(h.sameOrigin(h.postLogout)))

	router.Get("/", h.logged(h.authed(h.getIndex)))
	router.Get("/months/{month}", h.logged(h.authed(h.getMonth)))
	router.Get("/months/{month}/expenses", h.logged(h.authed(h.getExpenses)))
	router.Get("/trends", h.logged(h.authed(h.getTrends)))
	router.Get("/movements/new", h.logged(h.authed(h.newMovement)))
	router.Get("/movements/{id}/edit", h.logged(h.authed(h.editMovement)))
	router.Get("/movements/{id}/row", h.logged(h.authed(h.getRow)))
	router.Get("/movements/{id}/delete", h.logged(h.authed(h.confirmDelete)))
	router.Get("/movements/{id}/categories", h.logged(h.authed(h.categoryMenu)))
	router.Post("/movements", h.logged(h.sameOrigin(h.authed(h.createMovement))))
	router.Put("/movements/{id}", h.logged(h.sameOrigin(h.authed(h.updateMovement))))
	router.Delete("/movements/{id}", h.logged(h.sameOrigin(h.authed(h.deleteMovement))))
	router.Post("/movements/{id}/restore", h.logged(h.sameOrigin(h.authed(h.restoreMovement))))
	router.Post("/merchant-rules", h.logged(h.sameOrigin(h.authed(h.recategoriseMerchant))))
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
		name, _, _ := strings.Cut(f.Tag.Get("form"), ",")
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

func (h *Handler) parseForm(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("%w: %w", errMalformedForm, err)
	}
	return nil
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *Handler) authed(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := h.app.Authenticate(r.Context(), sessionToken(r))
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
	return h.now().In(h.zone)
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

func (h *Handler) renderPage(w http.ResponseWriter, r *http.Request, page pageView) {
	if isHTMX(r) {
		h.render(w, r, http.StatusOK, "swap", page)
		return
	}
	h.render(w, r, http.StatusOK, "layout", page)
}

func (h *Handler) renderAfterChange(w http.ResponseWriter, r *http.Request, month string, toast toastView) {
	page, err := h.monthPage(r.Context(), month, domain.ExpenseFilter{})
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	page.Toast = &toast
	page.CloseModal = true
	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("HX-Push-Url", "/months/"+month)
	h.render(w, r, http.StatusOK, "swap", page)
}

func (h *Handler) monthPage(ctx context.Context, month string, filter domain.ExpenseFilter) (pageView, error) {
	months, err := h.app.ListMonths(ctx)
	if err != nil {
		return pageView{}, err
	}
	report, err := h.app.MonthReport(ctx, month, filter)
	if err != nil {
		return pageView{}, err
	}
	view := mapMonthView(report)
	return pageView{Title: view.Title, Sidebar: mapSidebarView(months, month), Month: &view}, nil
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

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
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
	status, message := mapErrorResponse(err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if isHTMX(r) {
		w.Header().Set("HX-Reswap", "none")
		var buf bytes.Buffer
		_ = h.templates.ExecuteTemplate(&buf, "toast-oob", toastView{Message: message})
		w.WriteHeader(status)
		_, _ = w.Write(buf.Bytes())
		return
	}
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "<!doctype html><title>Abacus</title><p>%s</p><p><a href=\"/\">Back to Abacus</a></p>", template.HTMLEscapeString(message))
}
