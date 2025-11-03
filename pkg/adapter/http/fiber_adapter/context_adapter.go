package customfiber

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"

	"github.com/gofiber/fiber/v2"
)

// fiberContextAdapter adapts fiber.Ctx to contracts.Context
type fiberContextAdapter struct {
	ctx *fiber.Ctx
}

// NewFiberContextAdapter creates a new Fiber context adapter
func NewFiberContextAdapter(ctx *fiber.Ctx) contracts.Context {
	return &fiberContextAdapter{ctx: ctx}
}

func (f *fiberContextAdapter) Request() *http.Request {
	// Convert fasthttp request to http.Request
	req := new(http.Request)
	req.Method = f.ctx.Method()
	req.URL, _ = url.Parse(f.ctx.OriginalURL())
	req.Proto = f.ctx.Protocol()
	req.Host = f.ctx.Hostname()
	req.RemoteAddr = f.ctx.IP()
	req.RequestURI = string(f.ctx.Request().RequestURI())
	req.Header = make(http.Header)
	f.ctx.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})
	return req
}

func (f *fiberContextAdapter) ResponseWriter() http.ResponseWriter {
	return &fiberResponseWriter{ctx: f.ctx}
}

func (f *fiberContextAdapter) Param(name string) string {
	return f.ctx.Params(name)
}

func (f *fiberContextAdapter) QueryParam(name string) string {
	return f.ctx.Query(name)
}

func (f *fiberContextAdapter) QueryParams() url.Values {
	params := url.Values{}
	f.ctx.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		params.Add(string(key), string(value))
	})
	return params
}

func (f *fiberContextAdapter) FormValue(name string) string {
	return f.ctx.FormValue(name)
}

func (f *fiberContextAdapter) FormFile(name string) (*multipart.FileHeader, error) {
	return f.ctx.FormFile(name)
}

func (f *fiberContextAdapter) MultipartForm() (*multipart.Form, error) {
	return f.ctx.MultipartForm()
}

func (f *fiberContextAdapter) Get(key string) interface{} {
	return f.ctx.Locals(key)
}

func (f *fiberContextAdapter) Set(key string, val interface{}) {
	f.ctx.Locals(key, val)
}

func (f *fiberContextAdapter) Bind(i interface{}) error {
	return f.ctx.BodyParser(i)
}

func (f *fiberContextAdapter) Validate(i interface{}) error {
	return nil
}

func (f *fiberContextAdapter) Body() []byte {
	return f.ctx.Body()
}

func (f *fiberContextAdapter) JSON(code int, i interface{}) error {
	return f.ctx.Status(code).JSON(i)
}

func (f *fiberContextAdapter) JSONBlob(code int, b []byte) error {
	f.ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	f.ctx.Status(code)
	return f.ctx.Send(b)
}

func (f *fiberContextAdapter) XML(code int, i interface{}) error {
	return f.ctx.Status(code).XML(i)
}

func (f *fiberContextAdapter) String(code int, s string) error {
	f.ctx.Status(code)
	return f.ctx.SendString(s)
}

func (f *fiberContextAdapter) HTML(code int, html string) error {
	f.ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	f.ctx.Status(code)
	return f.ctx.SendString(html)
}

func (f *fiberContextAdapter) Blob(code int, contentType string, b []byte) error {
	f.ctx.Set(fiber.HeaderContentType, contentType)
	f.ctx.Status(code)
	return f.ctx.Send(b)
}

func (f *fiberContextAdapter) Stream(code int, contentType string, r io.Reader) error {
	f.ctx.Set(fiber.HeaderContentType, contentType)
	f.ctx.Status(code)
	return f.ctx.SendStream(r)
}

func (f *fiberContextAdapter) NoContent(code int) error {
	return f.ctx.SendStatus(code)
}

func (f *fiberContextAdapter) Redirect(code int, url string) error {
	return f.ctx.Redirect(url, code)
}

func (f *fiberContextAdapter) File(filepath string) error {
	return f.ctx.SendFile(filepath)
}

func (f *fiberContextAdapter) Attachment(filepath, filename string) error {
	if filename != "" {
		f.ctx.Set(fiber.HeaderContentDisposition, `attachment; filename="`+filename+`"`)
	}
	return f.ctx.SendFile(filepath)
}

func (f *fiberContextAdapter) GetHeader(key string) string {
	return f.ctx.Get(key)
}

func (f *fiberContextAdapter) SetHeader(key, value string) {
	f.ctx.Set(key, value)
}

func (f *fiberContextAdapter) Cookie(name string) (*http.Cookie, error) {
	value := f.ctx.Cookies(name)
	if value == "" {
		return nil, http.ErrNoCookie
	}
	return &http.Cookie{
		Name:  name,
		Value: value,
	}, nil
}

func (f *fiberContextAdapter) SetCookie(cookie *http.Cookie) {
	// Convert http.SameSite to fiber.Cookie SameSite
	var sameSite string
	switch cookie.SameSite {
	case http.SameSiteDefaultMode:
		sameSite = "lax"
	case http.SameSiteLaxMode:
		sameSite = "lax"
	case http.SameSiteStrictMode:
		sameSite = "strict"
	case http.SameSiteNoneMode:
		sameSite = "none"
	default:
		sameSite = "lax"
	}

	f.ctx.Cookie(&fiber.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Expires:  cookie.Expires,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HttpOnly,
		SameSite: sameSite,
	})
}

func (f *fiberContextAdapter) Cookies() []*http.Cookie {
	// Fiber doesn't provide direct access to all cookies
	// This is a simplified implementation
	return nil
}

func (f *fiberContextAdapter) Accepts(offers ...string) string {
	return f.ctx.Accepts(offers...)
}

func (f *fiberContextAdapter) ContentType() string {
	return f.ctx.Get(fiber.HeaderContentType)
}

func (f *fiberContextAdapter) Status(code int) contracts.Context {
	f.ctx.Status(code)
	return f
}

func (f *fiberContextAdapter) GetStatus() int {
	return f.ctx.Response().StatusCode()
}

func (f *fiberContextAdapter) Error(err error) {
	_ = f.ctx.App().Config().ErrorHandler(f.ctx, err)
}

func (f *fiberContextAdapter) Handler() interface{} {
	return f.ctx.Route().Handlers
}

func (f *fiberContextAdapter) SetHandler(h interface{}) {
	// Not directly supported in Fiber
}

func (f *fiberContextAdapter) Path() string {
	return f.ctx.Path()
}

func (f *fiberContextAdapter) RealIP() string {
	return f.ctx.IP()
}

func (f *fiberContextAdapter) Scheme() string {
	return f.ctx.Protocol()
}

// fiberResponseWriter wraps fiber.Ctx to implement http.ResponseWriter
type fiberResponseWriter struct {
	ctx *fiber.Ctx
}

func (w *fiberResponseWriter) Header() http.Header {
	header := make(http.Header)
	w.ctx.Response().Header.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})
	return header
}

func (w *fiberResponseWriter) Write(b []byte) (int, error) {
	return w.ctx.Write(b)
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.ctx.Status(statusCode)
}

// ConvertFiberHandler converts contracts.HandlerFunc to fiber.Handler
func ConvertFiberHandler(h contracts.HandlerFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		adapter := NewFiberContextAdapter(c)
		return h(adapter)
	}
}

// ConvertFiberMiddleware converts contracts.MiddlewareFunc to fiber middleware
func ConvertFiberMiddleware(m contracts.MiddlewareFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		adapter := NewFiberContextAdapter(c)
		handler := m(func(ctx contracts.Context) error {
			return c.Next()
		})
		return handler(adapter)
	}
}
