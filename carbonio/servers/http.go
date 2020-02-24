package servers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/handlers"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/static"
)

var (
	htmlTmpls map[string]*template.Template = map[string]*template.Template{}
	device    devices.Device
)

func HttpServer(port int, device devices.Device) {
	if device == nil {
		log.Printf("device is unitialized")
		return
	}

	h, err := handlers.NewHandlers(device,
		handlers.Port(port))
	if err != nil {
		log.Printf("error instantiating handlers; %s", err)
		return
	}

	ip := device.IP()
	var host string
	if ip.DefaultMask() != nil {
		host = ip.String()
	} else {
		host = fmt.Sprintf("[%s]", ip)
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("carbonio server starting on http://%s\n", addr)

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	// TODO(2020-02-24) Add logging for /static requests.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(static.AssetFile())))

	r.HandleFunc("/list", h.ListHandler)
	r.HandleFunc("/list_query", h.ListQueryHandler)
	r.HandleFunc("/status", h.StatusHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("server started")

	log.SetFlags(0)
	log.Fatal(srv.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("Hello, world!\n"))
	buf.WriteString(fmt.Sprintf("r.URL.Path: %s\n", r.URL.Path))

	status := http.StatusOK
	l, err := w.Write(buf.Bytes())
	if err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
	}
	helpers.CommonLogFormat(r, status, l)
}
