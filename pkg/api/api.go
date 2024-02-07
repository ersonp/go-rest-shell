package api

import (
	"fmt"
	"log"
	"net/http"
)

type API struct {
	http.Handler
	Host string
	Port int
	log  *log.Logger
}

func New(host string, port int, logger *log.Logger) *API {
	api := &API{
		Host: host,
		Port: port,
		log:  logger,
	}
	api.initRoutes()
	return api
}

func (api *API) initRoutes() {
	http.HandleFunc("/", api.helloHandler)
}

func (api *API) helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Handle the POST request
	fmt.Fprintf(w, "Hello, you've posted: %s\n", r.URL.Path)
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.log.Printf("%s - %s %s", r.RemoteAddr, r.Method, r.URL)
	http.DefaultServeMux.ServeHTTP(w, r)
}

func (api *API) Start() error {
	addr := fmt.Sprintf("%s:%d", api.Host, api.Port)
	api.log.Printf("API listening on %s", addr)
	err := http.ListenAndServe(addr, api)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
