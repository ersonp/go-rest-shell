package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

type API struct {
	http.Handler
	Host string
	Port int
	log  *log.Logger
}

type Shell struct {
	Command string `json:"command"`
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
	http.HandleFunc("/shell", api.shellHandler)
}

func (api *API) shellHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the command from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the request body

	// Unmarshal the request body into a Shell struct
	var shell Shell
	if err := json.Unmarshal(body, &shell); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Execute the command
	cmd := exec.Command("sh", "-c", shell.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
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
	fmt.Fprintf(w, "Command output:\n%s\n", output)

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
