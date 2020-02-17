/*
Package handlers provides common routines for handling requests.
*/
package handlers

import (
	"html/template"
	"net/http"
)

type Handler interface {
	http.Handler

	// Name of the handler.
	Name() string
}

const (
	ifs = " "
	ofs = " "
)

var htmlTmpls = map[string]*template.Template{}
