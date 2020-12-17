package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/paulcager/go-http-middleware"
	"github.com/paulcager/gosdata"
	"github.com/paulcager/osgridref"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
)

const (
	apiVersion = "v4"
)

var (
	port       string
	dataDir    string
	tileServer *gosdata.TileServer
)

func main() {
	flag.StringVarP(&port, "port", "p", ":9091", "Port to listen on")
	flag.StringVarP(&dataDir, "dataDir", "d", "data", "Directory containing tiles")
	flag.Parse()

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	tileServer = gosdata.NewTileServer(dataDir)

	server := makeHTTPServer(port)
	log.Fatal(server.ListenAndServe())
}

func makeHTTPServer(listenPort string) *http.Server {
	http.Handle(
		"/metrics",
		promhttp.Handler())

	http.Handle(
		"/"+apiVersion+"/height/",
		middleware.MakeLoggingHandler(http.HandlerFunc(handle)))

	log.Println("Starting HTTP server on " + listenPort)

	s := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       10 * time.Minute,
		Handler:           http.DefaultServeMux,
		Addr:              listenPort,
	}

	return s
}

func handle(w http.ResponseWriter, r *http.Request) {
	gridRefStr := r.URL.Path[len("/"+apiVersion+"/height/"):]
	gridRef, err := osgridref.ParseOsGridRef(gridRefStr)
	if err != nil {
		handleError(w, r, gridRefStr, err)
		return
	}

	type Reply struct {
		OSGridRef string  `json:"osGridRef"`
		Easting   int     `json:"easting"`
		Northing  int     `json:"northing"`
		Height    float64 `json:"height"`
	}

	height, err := tileServer.Height(gridRefStr)
	if err != nil {
		handleError(w, r, gridRefStr, err)
		return
	}

	reply := Reply{
		OSGridRef: gridRef.StringN(8),
		Easting:   gridRef.Easting,
		Northing:  gridRef.Northing,
		Height:    height,
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		log.Printf("Failed to write response: %s", err)
		w.WriteHeader(http.StatusBadGateway)
	}
}

func handleError(w http.ResponseWriter, _ *http.Request, str string, err error) {
	w.WriteHeader(http.StatusBadRequest)
	if err != nil {
		fmt.Fprintf(w, "Invalid request: %q (%s)\n", str, err)
	} else {
		fmt.Fprintf(w, "Invalid request: %q\n", str)
	}
}
