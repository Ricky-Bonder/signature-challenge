package api

import (
	"encoding/base64"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

func (s *Server) Device(response http.ResponseWriter, request *http.Request) {
	r := gin.Default()

	signatureService := domain.NewSignatureService()

	r.POST("/signature-device", func(c *gin.Context) {
		var req domain.CreateSignatureDeviceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := uuid.New().String()

		// Create a new signature device
		signatureDevice := &domain.SignatureDevice{
			ID:               id,
			Algorithm:        req.Algorithm,
			Label:            req.Label,
			SignatureCounter: new(int),
		}

		// Store the signature device
		signatureService.Mutex.Lock()
		defer signatureService.Mutex.Unlock()
		signatureService.Devices[id] = signatureDevice

		// Respond with success
		c.JSON(http.StatusCreated, gin.H{"id": id})
	})

	// Endpoint to sign a transaction
	r.POST("/sign-transaction", func(c *gin.Context) {
		// Retrieve device ID and data from the request
		deviceID := c.PostForm("deviceId")
		data := c.PostForm("data")

		// Find the signature device
		signatureService.Mutex.Lock()
		thisDevice, found := signatureService.Devices[deviceID]
		signatureService.Mutex.Unlock()
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Signature device not found"})
			return
		}

		// Construct secured data to be signed
		_ = "_" + data + "_" + base64.StdEncoding.EncodeToString([]byte(thisDevice.ID))

		//logger.Info("Secured data to be signed: " + securedDataToBeSigned)
		// Increment signature counter
		*thisDevice.SignatureCounter++

		// Generate signature (for demonstration, we use a placeholder)
		signature := "base64_encoded_signature"

		// Construct signed data
		signedData := strconv.Itoa(*thisDevice.SignatureCounter) + "_" + data + "_" + base64.StdEncoding.EncodeToString([]byte(thisDevice.ID))

		// Respond with signature response
		c.JSON(http.StatusOK, domain.SignatureResponse{
			Signature:  signature,
			SignedData: signedData,
		})
	})
}
