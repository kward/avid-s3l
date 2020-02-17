/*
Package handlers provides common routines for handling requests.
*/
package handlers

import (
	"html/template"
	"io"
	"net/http"
)

type Handler interface {
	http.Handler

	// ServeCommand handles a command request.
	ServeCommand(w io.Writer)
	// Name of the handler.
	Name() string
}

const (
	ifs = " "
	ofs = " "
)

var htmlTmpls = map[string]*template.Template{}
