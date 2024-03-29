package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"net/http"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	URL           string
	listenAddress string
	storage       persistence.Storage
}

// NewServer is a factory to instantiate a new Server.
func NewServer(
	URL string,
	listenAddress string,
	storage persistence.Storage,
) *Server {
	return &Server{
		URL:           URL,
		listenAddress: listenAddress,
		storage:       storage,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))
	mux.HandleFunc("/api/v0/create-signature-device", s.CreateSignatureDevice)
	mux.HandleFunc("/api/v0/sign-transaction", s.SignTransaction)
	mux.HandleFunc("/api/v0/get-signature-device", s.GetSignatureDevice)
	mux.HandleFunc("/api/v0/get-all-devices", s.GetAllSignatureDevices)

	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	if err != nil {
		return
	}
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	_, err = w.Write(bytes)
	if err != nil {
		return
	}
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
