/*
Package handlers provides common routines for handling requests.

All templates are loaded in the `init()` of this file.

TODO(2020-02-19) Provide HTML error pages.
*/
package handlers

import (
	"fmt"
	"html/template"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/templates"
)

type Handlers struct {
	opts   *options
	device devices.Device
}

func NewHandlers(device devices.Device, opts ...func(*options) error) (*Handlers, error) {
	if device == nil {
		return nil, fmt.Errorf("device is uninitialized")
	}

	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, fmt.Errorf("invalid option; %s", err)
		}
	}
	if err := o.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate options; %s", err)
	}

	return &Handlers{
		opts:   o,
		device: device,
	}, nil
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
		if err := loadAndParse("html/" + a); err != nil {
			helpers.Exit(fmt.Sprintf("%s", err))
		}
	}
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
