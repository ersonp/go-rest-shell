package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/ersonp/go-rest-shell/pkg/shell"
)

type API struct {
	Host string
	Port int
	log  *log.Logger
	Mux  *http.ServeMux
}

type Executor struct {
	Command string `json:"command"`
}

func (ex *Executor) Validate() error {
	if ex.Command == "" {
		return errors.New("command is required")
	}
	return nil
}

func New(host string, port int, logger *log.Logger) *API {
	api := &API{
		Host: host,
		Port: port,
		log:  logger,
		Mux:  http.NewServeMux(),
	}
	api.initRoutes()
	return api
}

func (api *API) initRoutes() {
	api.Mux.HandleFunc("POST /api/cmd", api.cmdHandler)
}

func (api *API) cmdHandler(w http.ResponseWriter, r *http.Request) {

	// Parse the command from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.log.Printf("Failed to read request body: %v, Method: %s, URL: %s, RemoteAddr: %s",
			err, r.Method, r.URL, r.RemoteAddr)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the request body

	// Unmarshal the request body into a Shell struct
	var ex Executor
	if err := json.Unmarshal(body, &ex); err != nil {
		api.log.Printf("Failed to parse request body: %v, Method: %s, URL: %s, RemoteAddr: %s",
			err, r.Method, r.URL, r.RemoteAddr)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Validate the Executor
	if err := ex.Validate(); err != nil {
		api.log.Printf("Invalid request body: %v, Method: %s, URL: %s, RemoteAddr: %s",
			err, r.Method, r.URL, r.RemoteAddr)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute the command depending on the OS
	cmd := shell.Execute(ex.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		api.log.Printf("Failed to execute command: %v, Command: %s, Method: %s, URL: %s, RemoteAddr: %s",
			err, ex.Command, r.Method, r.URL, r.RemoteAddr)
		// Check if the command failed because it was not found
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 127 {
				http.Error(w, "Command not found", http.StatusNotFound)
				return
			}
		}
		// Return a generic error if the command execution failed
		http.Error(w, "Failed to execute command", http.StatusInternalServerError)
		return
	}

	// Write the command output to the response
	fmt.Fprintf(w, "%s\n", output)

}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.log.Printf("%s - %s %s", r.RemoteAddr, r.Method, r.URL)
	api.Mux.ServeHTTP(w, r)
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
