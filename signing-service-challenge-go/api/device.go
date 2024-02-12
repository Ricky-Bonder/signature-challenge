package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"
)

func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	signatureService := domain.GetSignatureService()
	signatureService.Mutex.Lock()
	defer signatureService.Mutex.Unlock()
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		WriteInternalError(response)
		return
	}

	var data domain.CreateSignatureDeviceRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(response, "Failed to parse JSON body", http.StatusBadRequest)
		return
	}

	// Access the field values
	fmt.Println("Received field 1:", data.Algorithm)
	fmt.Println("Received field 2:", data.Label)

	var generator crypto.KeyPairGenerator
	algo := domain.AlgorithmNames[data.Algorithm]
	switch algo {
	case domain.RSA:
		generator = &crypto.RSAGenerator{}
	case domain.ECC:
		generator = &crypto.ECCGenerator{}
	default:
		WriteErrorResponse(response, http.StatusNotImplemented, []string{
			http.StatusText(http.StatusNotImplemented),
		})
		return
	}

	counter := domain.Increment().Get()
	fmt.Println("signature counter:", counter)

	signatureDevice := &domain.InternalSignatureDevice{
		ID:               uuid.New().String(),
		Algorithm:        generator,
		Label:            data.Label,
		SignatureCounter: counter,
	}

	err = s.storage.CreateSignatureDevice(signatureDevice)
	if err != nil {
		return
	}

	signatureService.Devices[signatureDevice.ID] = signatureDevice
	signatureResponse := CreateSignatureDeviceResponse(signatureDevice.ID, signatureDevice.Algorithm.GetAlgorithm(), *signatureDevice.Label)
	WriteAPIResponse(response, http.StatusCreated, signatureResponse)
}

func (s *Server) SignTransaction(response http.ResponseWriter, request *http.Request) {
	signatureService := domain.GetSignatureService()
	signatureService.Mutex.Lock()
	defer signatureService.Mutex.Unlock()
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		WriteInternalError(response)
		return
	}

	var data domain.SignTransactionRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(response, "Failed to parse JSON body", http.StatusBadRequest)
		return
	}

	fmt.Println("Received field 1:", data.ID)
	fmt.Println("Received field 2:", data.Data)

	device := signatureService.Devices[data.ID]
	fmt.Println("device:", device)

	if device == nil {
		WriteErrorResponse(response, http.StatusNotFound, []string{
			http.StatusText(http.StatusNotFound),
		})
		return
	}

	var dataToSign string
	if device.SignatureCounter == 0 {
		dataToSign = data.Data + "" + strconv.Itoa(int(device.SignatureCounter)) + "" + base64.StdEncoding.EncodeToString([]byte(device.ID))
	} else {
		//TODO get last signature
		//dataToSign = data.Data + "" + strconv.Itoa(int(device.SignatureCounter)) + "" + base64.StdEncoding.EncodeToString([]byte(LAST_SIGNATURE))
	}

	keypair, err := device.Algorithm.Generate()
	var signatureResponse *domain.SignatureResponse
	if device.Algorithm.GetAlgorithm() == "RSA" {
		rsaKeyPair, err := crypto.CastToRSAKeyPair(keypair)
		if err != nil {
			return
		}
		signature, err := crypto.SignRSA(rsaKeyPair, []byte(dataToSign))
		if err != nil {
			return
		}
		signatureResponse = &domain.SignatureResponse{
			Signature:  base64.StdEncoding.EncodeToString(signature),
			SignedData: string(signature),
		}
	} else if device.Algorithm.GetAlgorithm() == "ECC" {
		eccKeyPair, err := crypto.CastToECCKeyPair(keypair)
		if err != nil {
			return
		}
		signature, err := crypto.SignECC(eccKeyPair, []byte(dataToSign))
		if err != nil {
			return
		}
		signatureResponse = &domain.SignatureResponse{
			Signature:  base64.StdEncoding.EncodeToString(signature),
			SignedData: string(signature),
		}
	}

	WriteAPIResponse(response, http.StatusCreated, signatureResponse)

}

func (s *Server) GetSignatureDevice(response http.ResponseWriter, request *http.Request) {
	signatureService := domain.GetSignatureService()
	signatureService.Mutex.Lock()
	defer signatureService.Mutex.Unlock()
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	queryParams := request.URL.Query()

	// Get the value of a specific query parameter (e.g., "id")
	id := queryParams.Get("id")

	signatureDevice, err := persistence.GetMemoryStorage().GetSignatureDevice(id)
	if err != nil {
		return
	}
	fmt.Println("device:", signatureDevice)

	if signatureDevice == nil {
		WriteErrorResponse(response, http.StatusNotFound, []string{
			http.StatusText(http.StatusNotFound),
		})
		return
	}

	signatureResponse := CreateSignatureDeviceResponse(signatureDevice.ID, signatureDevice.Algorithm.GetAlgorithm(), *signatureDevice.Label)

	WriteAPIResponse(response, http.StatusFound, signatureResponse)
}

func (s *Server) GetAllSignatureDevices(response http.ResponseWriter, request *http.Request) {
	signatureService := domain.GetSignatureService()
	signatureService.Mutex.Lock()
	defer signatureService.Mutex.Unlock()
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	signatureDevices, err := persistence.GetMemoryStorage().GetAllSignatureDevices()
	if err != nil {
		return
	}
	fmt.Println("list all devices:", signatureDevices)

	if signatureDevices == nil {
		WriteErrorResponse(response, http.StatusNotFound, []string{
			http.StatusText(http.StatusNotFound),
		})
		return
	}

	var signatureResponse []*domain.CreateSignatureDeviceResponse
	for _, signatureDevice := range signatureDevices {
		signatureResponse = append(signatureResponse, CreateSignatureDeviceResponse(signatureDevice.ID, signatureDevice.Algorithm.GetAlgorithm(), *signatureDevice.Label))
	}

	WriteAPIResponse(response, http.StatusFound, signatureResponse)
}

func CreateSignatureDeviceResponse(id string, algorithm string, label string) *domain.CreateSignatureDeviceResponse {
	return &domain.CreateSignatureDeviceResponse{
		ID:        id,
		Algorithm: algorithm,
		Label:     label,
	}
}
