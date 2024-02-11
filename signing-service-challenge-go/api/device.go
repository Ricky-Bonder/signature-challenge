package api

import (
	"encoding/json"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func (s *Server) Signature(response http.ResponseWriter, request *http.Request) {
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

	generatedKeyPair, err := signatureDevice.Algorithm.Generate()
	if err != nil {
		return
	}
	fmt.Println("generated keypair:", generatedKeyPair.Public, " - ", generatedKeyPair.Private)
	signatureService.Devices[signatureDevice.ID] = signatureDevice

	signatureResponse := &domain.CreateSignatureDeviceResponse{
		ID:        signatureDevice.ID,
		Algorithm: signatureDevice.Algorithm.GetAlgorithm(),
		Label:     *signatureDevice.Label,
	}
	WriteAPIResponse(response, http.StatusCreated, signatureResponse)
}

//func (s *Server) CreateSignatureDeviceHandler(response http.ResponseWriter, request *http.Request) {
//	r := gin.Default()
//
//	if request.Method != http.MethodPost {
//		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{http.StatusText(http.StatusMethodNotAllowed)})
//		return
//	}
//
//	body, err := io.ReadAll(request.Body)
//	if err != nil {
//		WriteInternalError(response)
//		return
//	}
//
//	signatureService := domain.NewSignatureService()
//
//	r.POST("/signature-device", func(c *gin.Context) {
//		signatureService.Mutex.Lock()
//		defer signatureService.Mutex.Unlock()
//		var req domain.CreateSignatureDeviceRequest
//		if err := c.ShouldBindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		if req.Algorithm != "RSA" && req.Algorithm != "ECC" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid algorithm"})
//			return
//		}
//
//		switch req.Algorithm {
//		case "RSA":
//			handleRSAAlgorithm()
//		case "ECC":
//			handleECCAlgorithm()
//		}
//
//		//Parse request parameters
//		id := uuid.New().String()
//		label := request.FormValue("label")
//		counter := domain.SignatureCounter.Increment
//
//		// Create a new signature device
//		signatureDevice := &domain.InternalSignatureDevice{
//			ID:               id,
//			Algorithm:        req.Algorithm,
//			Label:            req.Label,
//			SignatureCounter: counter,
//		}
//
//		signatureService.Devices[id] = signatureDevice
//
//		// Respond with success
//		c.JSON(http.StatusCreated, gin.H{"id": id})
//
//		WriteAPIResponse(response, http.StatusCreated, signatureDevice)
//	})
//}

//func (s *Server) Device(response http.ResponseWriter, request *http.Request) {
//	r := gin.Default()
//
//	signatureService := domain.NewSignatureService()
//
//	r.POST("/signature-device", func(c *gin.Context) {
//		var req domain.CreateSignatureDeviceRequest
//		if err := c.ShouldBindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		id := uuid.New().String()
//
//		// Create a new signature device
//		signatureDevice := &domain.InternalSignatureDevice{
//			ID:               id,
//			Algorithm:        req.Algorithm,
//			Label:            req.Label,
//			SignatureCounter: new(int),
//		}
//
//		// Store the signature device
//		signatureService.Mutex.Lock()
//		defer signatureService.Mutex.Unlock()
//		signatureService.Devices[id] = signatureDevice
//
//		// Respond with success
//		c.JSON(http.StatusCreated, gin.H{"id": id})
//	})
//
//	// Endpoint to sign a transaction
//	r.POST("/sign-transaction", func(c *gin.Context) {
//		// Retrieve device ID and data from the request
//		deviceID := c.PostForm("deviceId")
//		data := c.PostForm("data")
//
//		// Find the signature device
//		signatureService.Mutex.Lock()
//		thisDevice, found := signatureService.Devices[deviceID]
//		signatureService.Mutex.Unlock()
//		if !found {
//			c.JSON(http.StatusNotFound, gin.H{"error": "Signature device not found"})
//			return
//		}
//
//		// Construct secured data to be signed
//		_ = "_" + data + "_" + base64.StdEncoding.EncodeToString([]byte(thisDevice.ID))
//
//		//logger.Info("Secured data to be signed: " + securedDataToBeSigned)
//		// Increment signature counter
//		*thisDevice.SignatureCounter++
//
//		// Generate signature (for demonstration, we use a placeholder)
//		signature := "base64_encoded_signature"
//
//		// Construct signed data
//		signedData := strconv.Itoa(*thisDevice.SignatureCounter) + "_" + data + "_" + base64.StdEncoding.EncodeToString([]byte(thisDevice.ID))
//
//		// Respond with signature response
//		c.JSON(http.StatusOK, domain.SignatureResponse{
//			Signature:  signature,
//			SignedData: signedData,
//		})
//	})
//}
