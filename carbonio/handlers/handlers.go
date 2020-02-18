/*
Package handlers provides common routines for handling requests.

All templates are loaded in the `init()` of this file.
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

func init() {
	// Load and parse HTML templates.
	as, err := templates.AssetDir("html")
	if err != nil {
		helpers.Exit(fmt.Sprintf("unable to locate html templates"))
	}
	for _, a := range as {
		if err := loadTemplate("html/" + a); err != nil {
			helpers.Exit(fmt.Sprintf("%s", err))
		}
	}
}

// loadAndParse template from the bindata asset and store in the package global
// `tmpls` map.
func loadTemplate(name string) error {
	fmt.Printf("loading %s template\n", name)
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
