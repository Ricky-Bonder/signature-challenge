package api

import (
	"bytes"
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {

	s := NewServer("http://localhost", ":8080", persistence.MemoryStorage{})

	if s.listenAddress != ":8080" {
		t.Error("listenAddress should be :8080")
	}

	if s.storage == nil {
		t.Error("storage should not be nil")
	}

	err := s.Run()
	if err != nil {
		t.Error("Running server should not return an error")
	}
}

func TestHealthHandlerRR(t *testing.T) {

	s := NewServer("http://localhost", ":8080", persistence.MemoryStorage{})

	req, err := http.NewRequest(http.MethodGet, "/api/v0/health", nil)
	if err != nil {
		t.Error("error on Health GET request:	", err)
	}
	rr := httptest.NewRecorder()

	s.Health(rr, req)

	if err != nil {
		t.Error("error on Health GET request")
	}

	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Health should return 200, but was ", rr.Result().StatusCode)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error("error closing response body")
		}
	}(rr.Result().Body)

	//expectedBody := "{\n  \"data\": {\n    \"status\": \"pass\",\n    \"version\": \"v0\"\n  }\n}"
	//
	bodyBytes, err := io.ReadAll(rr.Result().Body)
	//
	//assert.Equal(t, expectedBody, string(bodyBytes))

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]map[string]string
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)
	data, ok := response["data"]
	if !ok {
		t.Error("response does not contain 'data' field")
	}

	status, ok := data["status"]
	if !ok {
		t.Error("data field does not contain 'status' field")
	}

	version, ok := data["version"]
	if !ok {
		t.Error("data field does not contain 'version' field")
	}

	// Perform assertions
	assert.Equal(t, "pass", status)
	assert.Equal(t, "v0", version)

}

func TestCreateSignatureDeviceHandler(t *testing.T) {
	router := gin.Default()
	NewServer("http://localhost", ":8080", persistence.MemoryStorage{})

	body := []byte(`{
		"algorithm": "RSA",
		"label": "My device"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/signature-device", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response gin.H
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "pass", response["status"])
	assert.Equal(t, "v0", response["version"])

}
