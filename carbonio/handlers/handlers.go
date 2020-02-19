/*
Package handlers provides common routines for handling requests.

All templates are loaded in the `init()` of this file.

TODO(2020-02-19) Provide HTML error pages.
*/
package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/templates"
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

// tmpls holds all parsed templates.
var tmpls = map[string]*template.Template{}

// hndlrs holds all known handlers.
var hndlrs = []Handler{}

func init() {
	// Load and parse HTML templates.
	as, err := templates.AssetDir("html")
	if err != nil {
		helpers.Exit(fmt.Sprintf("unable to locate html templates"))
	}
	for _, a := range as {
		if err := loadAndParse("html/" + a); err != nil {
			helpers.Exit(fmt.Sprintf("%s", err))
		}
	}
}

// register a handler.
//
// The act of registering validates that the Handler interface is met.
func register(h Handler) {
	hndlrs = append(hndlrs, h)
}

// loadAndParse template from the bindata asset and store in the package global
// `tmpls` map.
func loadAndParse(name string) error {
	data, err := templates.Asset(name)
	if err != nil {
		return fmt.Errorf("failure loading %s asset", name)
	}
	tmpl, err := template.New("status").Parse(string(data))
	if err != nil {
		return fmt.Errorf("error parsing %s template; %s", name, err)
	}
	tmpls[name] = tmpl
	return nil
}

var assetNames map[string]bool

// mustTemplate checks that a template name is known, or returns an error.
func mustTemplate(name string) error {
	if assetNames == nil {
		// Cache the asset names.
		assetNames = map[string]bool{}
		for _, an := range templates.AssetNames() {
			assetNames[an] = true
		}
	}
	if _, ok := assetNames[name]; !ok {
		return fmt.Errorf("asset name %s unknown", name)
	}
	return nil
}
