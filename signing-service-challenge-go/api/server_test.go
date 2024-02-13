package api

import (
	"bytes"
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandlerRR(t *testing.T) {

	s := NewServer("http://localhost", ":8080", &persistence.DeviceStorage{})

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

	bodyBytes, err := io.ReadAll(rr.Result().Body)

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

	assert.Equal(t, "pass", status)
	assert.Equal(t, "v0", version)
}

func TestCreateSignatureDeviceHandler(t *testing.T) {
	storage := persistence.GetSignatureDeviceStorage()
	s := NewServer("http://localhost", ":8080", storage)

	body := []byte(`{
		"algorithm": "RSA",
		"label": "My device"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/api/v0/signature-device", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	s.CreateSignatureDevice(rr, req)

	if err != nil {
		t.Error("error on Health POST request")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error("error closing response body")
		}
	}(rr.Result().Body)

	bodyBytes, err := io.ReadAll(rr.Result().Body)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response map[string]map[string]string
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)
	data, ok := response["data"]
	if !ok {
		t.Error("response does not contain 'data' field")
	}
	algorithm, ok := data["algorithm"]
	if !ok {
		t.Error("data field does not contain 'status' field")
	}
	label, ok := data["label"]
	if !ok {
		t.Error("data field does not contain 'version' field")
	}

	assert.Equal(t, "RSA", algorithm)
	assert.Equal(t, "My device", label)

	devices := domain.GetSignatureService().Devices

	var device *domain.InternalSignatureDevice
	for _, v := range devices {
		device = v
		break // Exit the loop after the first iteration
	}

	assert.Equal(t, int32(0), device.SignatureCounter)

}

func TestCreateTwoSignatureDeviceHandler(t *testing.T) {
	storage := persistence.GetSignatureDeviceStorage()
	s := NewServer("http://localhost", ":8080", storage)

	body := []byte(`{
		"algorithm": "RSA",
		"label": "My device"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/api/v0/signature-device", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	s.CreateSignatureDevice(rr, req)

	if err != nil {
		t.Error("error on Health POST request")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error("error closing response body")
		}
	}(rr.Result().Body)

	bodyBytes, err := io.ReadAll(rr.Result().Body)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response map[string]map[string]string
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)
	data, ok := response["data"]
	if !ok {
		t.Error("response does not contain 'data' field")
	}
	algorithm, ok := data["algorithm"]
	if !ok {
		t.Error("data field does not contain 'status' field")
	}
	label, ok := data["label"]
	if !ok {
		t.Error("data field does not contain 'version' field")
	}

	assert.Equal(t, "RSA", algorithm)
	assert.Equal(t, "My device", label)

	body2 := []byte(`{
		"algorithm": "ECC",
		"label": "Another device"
	}`)
	req2, err2 := http.NewRequest(http.MethodPost, "/api/v0/signature-device", bytes.NewBuffer(body2))
	rr = httptest.NewRecorder()
	s.CreateSignatureDevice(rr, req2)

	if err2 != nil {
		t.Error("error on Health POST request")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error("error closing response body")
		}
	}(rr.Result().Body)

	bodyBytes, err = io.ReadAll(rr.Result().Body)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response2 map[string]map[string]string
	err = json.Unmarshal(bodyBytes, &response2)
	assert.NoError(t, err)
	data, ok = response2["data"]
	if !ok {
		t.Error("response does not contain 'data' field")
	}
	algorithm, ok = data["algorithm"]
	if !ok {
		t.Error("data field does not contain 'status' field")
	}
	label, ok = data["label"]
	if !ok {
		t.Error("data field does not contain 'version' field")
	}

	assert.Equal(t, "ECC", algorithm)
	assert.Equal(t, "Another device", label)
}

func TestSignTransactionHandler(t *testing.T) {
	storage := persistence.GetSignatureDeviceStorage()
	s := NewServer("http://localhost", ":8080", storage)

	body := []byte(`{
		"algorithm": "RSA",
		"label": "My device"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/api/v0/signature-device", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	s.CreateSignatureDevice(rr, req)

	if err != nil {
		t.Error("error on signature-device POST request")
	}
	assert.NoError(t, err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error("error closing response body")
		}
	}(rr.Result().Body)

	bodyBytes, err := io.ReadAll(rr.Result().Body)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response map[string]map[string]string
	err = json.Unmarshal(bodyBytes, &response)
	assert.NoError(t, err)
	data, ok := response["data"]
	if !ok {
		t.Error("response does not contain 'data' field")
	}
	deviceId, ok := data["id"]
	if !ok {
		t.Error("data field does not contain 'id' field")
	}
	algorithm, ok := data["algorithm"]
	if !ok {
		t.Error("data field does not contain 'algorithm' field")
	}
	label, ok := data["label"]
	if !ok {
		t.Error("data field does not contain 'label' field")
	}

	assert.Equal(t, "RSA", algorithm)
	assert.Equal(t, "My device", label)

	transactionBody := []byte(`{
		"id": "` + deviceId + `",
		"label": "Another device"
	}`)
	req2, err2 := http.NewRequest(http.MethodPost, "/api/v0/sign-transaction", bytes.NewBuffer(transactionBody))
	rr = httptest.NewRecorder()
	s.SignTransaction(rr, req2)

	if err2 != nil {
		t.Error("error on sign-transaction POST request")
	}

	assert.NoError(t, err2)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error("error closing response body")
		}
	}(rr.Result().Body)

	bodyBytes, err = io.ReadAll(rr.Result().Body)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response2 map[string]map[string]string
	err = json.Unmarshal(bodyBytes, &response2)
	assert.NoError(t, err)
	data, ok = response2["data"]
	if !ok {
		t.Error("response does not contain 'data' field")
	}
	algorithm, ok = data["signature"]
	if !ok {
		t.Error("data field does not contain 'signature' field")
	}
	label, ok = data["signed_data"]
	if !ok {
		t.Error("data field does not contain 'signed-data' field")
	}

}
