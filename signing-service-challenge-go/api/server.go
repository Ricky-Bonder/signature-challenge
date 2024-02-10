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
func NewServer(URL string, listenAddress string, storage persistence.Storage) *Server {
	return &Server{
		URL:           URL,
		listenAddress: listenAddress,
		storage:       storage,
		// TODO: add services / further dependencies here ...
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))
	mux.HandleFunc("/api/v0/create-signature-device", s.CreateSignatureDeviceHandler)
	// TODO: register further HandlerFuncs here ...

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

func (s *Server) CreateSignatureDeviceHandler(response http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	WriteAPIResponse(response, http.StatusOK, "success")

	// Parse request parameters
	//id := uuid.New().String() // Generate a unique UUID for the id
	//algorithm := r.FormValue("algorithm")
	//label := r.FormValue("label")
	//
	//// Process the request and create the signature device
	//
	//// Return the response with the generated id and any other relevant information
	//response := domain.SignatureDevice{
	//		ID               string `json:"id"`,
	//	Algorithm        Algorithm `json:"algorithm"`
	//	Label            *string   `json:"label"`
	//	SignatureCounter *int      `json:"signatureCounter"`
	//
	//}
	// Return the response to the client
	// ...
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
