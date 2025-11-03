// Package contracts
package contracts

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Context represents a framework-agnostic HTTP context
// Implements common functionality across Echo, Fiber, and Gin frameworks
type Context interface {
	// ===== Request Methods =====

	// Request returns the HTTP request
	Request() *http.Request

	// Response returns a response writer interface
	ResponseWriter() http.ResponseWriter

	// ===== Path & Query Parameters =====

	// Param returns path parameter by name
	Param(name string) string

	// QueryParam returns the query parameter by name
	QueryParam(name string) string

	// QueryParams returns the query parameters as url.Values
	QueryParams() url.Values

	// ===== Form Data & File Uploads =====

	// FormValue returns the form field value by name
	FormValue(name string) string

	// FormFile returns the multipart form file for the given key
	FormFile(name string) (*multipart.FileHeader, error)

	// MultipartForm returns the multipart form
	MultipartForm() (*multipart.Form, error)

	// ===== Context Storage =====

	// Get retrieves data from the context
	Get(key string) interface{}

	// Set saves data in the context
	Set(key string, val interface{})

	// ===== Request Body =====

	// Bind binds the request body into provided type
	// Supports JSON, XML, form data based on Content-Type
	Bind(i interface{}) error

	// Validate validates provided struct
	Validate(i interface{}) error

	// Body returns the raw request body
	Body() []byte

	// ===== Response Methods =====

	// JSON sends a JSON response with status code
	JSON(code int, i interface{}) error

	// JSONBlob sends a JSON blob response with status code
	JSONBlob(code int, b []byte) error

	// XML sends an XML response with status code
	XML(code int, i interface{}) error

	// String sends a string response with status code
	String(code int, s string) error

	// HTML sends an HTML response with status code
	HTML(code int, html string) error

	// Blob sends a blob response with status code and content type
	Blob(code int, contentType string, b []byte) error

	// Stream sends a streaming response with status code and content type
	Stream(code int, contentType string, r io.Reader) error

	// NoContent sends a response with no body and status code
	NoContent(code int) error

	// Redirect redirects the request to a provided URL with status code
	Redirect(code int, url string) error

	// ===== File Downloads =====

	// File sends a file as response
	File(filepath string) error

	// Attachment sends a file as an attachment (triggers download)
	Attachment(filepath, filename string) error

	// ===== Headers =====

	// GetHeader returns the request header value for the given key
	GetHeader(key string) string

	// SetHeader sets a response header
	SetHeader(key, value string)

	// ===== Cookies =====

	// Cookie returns the named cookie from the request
	Cookie(name string) (*http.Cookie, error)

	// SetCookie adds a Set-Cookie header to the response
	SetCookie(cookie *http.Cookie)

	// Cookies returns all cookies from the request
	Cookies() []*http.Cookie

	// ===== Content Negotiation =====

	// Accepts returns the best match from the offered content types
	// Returns empty string if no match is found
	Accepts(offers ...string) string

	// ContentType returns the Content-Type header of the request
	ContentType() string

	// ===== Status =====

	// Status sets the HTTP status code (for chaining)
	Status(code int) Context

	// GetStatus returns the current HTTP status code
	GetStatus() int

	// ===== Error Handling =====

	// Error invokes the registered error handler
	Error(err error)

	// ===== Handler & Metadata =====

	// Handler returns the matched handler by router
	Handler() interface{}

	// SetHandler sets the matched handler by router
	SetHandler(h interface{})

	// Path returns the registered path for the handler
	Path() string

	// RealIP returns the client's real IP address
	RealIP() string

	// Scheme returns the HTTP protocol scheme, http or https
	Scheme() string
}

// HandlerFunc defines a function to serve HTTP requests
type HandlerFunc func(Context) error

// MiddlewareFunc defines a function to process middleware
type MiddlewareFunc func(HandlerFunc) HandlerFunc
