package config

import (
	"os"
)

type HandlersConfig interface {
	BooksURL() string
	StaticsURL() string
	TemplatesURL() string
	BooksPageNumber() string
	BooksLimit() string
	NotFoundURL() string
	// JSON return string with text: 'application/json'
	JSON() string
	// PDF return string with text: 'application/pdf'
	PDF() string
	// HTML return string with text: 'text/html'
	HTML() string

	// FormatJSON return string with text: 'json'
	FormatJSON() string
	// FormatPDF return string with text: 'pdf'
	FormatPDF() string
	// FormatHTML return string with text: 'html'
	FormatHTML() string
}

type handlerConfig struct {
	// default value for url
	booksURL       string
	staticFilesURL string
	templatesPath  string
	notFoundURL    string

	// default value for pagination
	defaultBooksPageNumber string
	defaultBooksLimit      string

	// default value for header "Content-Type"
	headerJSON string
	headerPDF  string
	headerHTML string

	// default value for file format
	formatJSON string
	formatPDF  string
	formatHTML string
}

func NewHandlerConfig() HandlersConfig {
	return &handlerConfig{
		booksURL:               os.Getenv("BOOKS_URL"),
		staticFilesURL:         os.Getenv("STATIC_URL"),
		templatesPath:          os.Getenv("TEMPLATES_URL"),
		defaultBooksLimit:      os.Getenv("Default_Books_Limit"),
		defaultBooksPageNumber: os.Getenv("Default_Books_Page_Number"),
		notFoundURL:            os.Getenv("NOT_FOUND_URL"),
		headerJSON:             "application/json",
		headerPDF:              "application/pdf",
		headerHTML:             "text/html",
		formatJSON:             "json",
		formatPDF:              "pdf",
		formatHTML:             "html",
	}
}

// BooksURL return books url from env
func (cfg *handlerConfig) BooksURL() string {
	return cfg.booksURL
}

// StaticsURL return static url from env
func (cfg *handlerConfig) StaticsURL() string {
	return cfg.staticFilesURL
}

// TemplatesURL return template url from env
func (cfg *handlerConfig) TemplatesURL() string {
	return cfg.templatesPath
}

// NotFoundURL return notFound url
func (cfg *handlerConfig) NotFoundURL() string {
	return cfg.notFoundURL
}

func (cfg *handlerConfig) BooksPageNumber() string {
	return cfg.defaultBooksPageNumber
}

func (cfg *handlerConfig) BooksLimit() string {
	return cfg.defaultBooksLimit
}

// JSON return string with text: 'application/json'
func (cfg *handlerConfig) JSON() string {
	return cfg.headerJSON
}

// PDF return string with text: 'application/pdf'
func (cfg *handlerConfig) PDF() string {
	return cfg.headerPDF
}

// HTML return string with text: 'text/html'
func (cfg *handlerConfig) HTML() string {
	return cfg.headerHTML
}

func (cfg *handlerConfig) FormatJSON() string {
	return cfg.formatJSON
}

func (cfg *handlerConfig) FormatPDF() string {
	return cfg.formatPDF
}

func (cfg *handlerConfig) FormatHTML() string {
	return cfg.formatHTML
}
