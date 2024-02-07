package main

import (
	"encoding/base64"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	device "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	zap.RedirectStdLog(logger.Desugar())
	zap.ReplaceGlobals(logger.Desugar())
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)
	logger.Info("Starting server on " + ListenAddress)
	server := api.NewServer(ListenAddress)

	r := gin.Default()

	signatureService := device.NewSignatureService()

	r.POST("/signature-device", func(c *gin.Context) {
		var req device.CreateSignatureDeviceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a unique ID for the device (for simplicity, we use a timestamp)
		id := base64.StdEncoding.EncodeToString([]byte(time.Now().String()))

		// Create a new signature device
		signatureDevice := &device.SignatureDevice{
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
		securedDataToBeSigned := "_" + data + "_" + base64.StdEncoding.EncodeToString([]byte(thisDevice.ID))

		logger.Info("Secured data to be signed: " + securedDataToBeSigned)
		// Increment signature counter
		*thisDevice.SignatureCounter++

		// Generate signature (for demonstration, we use a placeholder)
		signature := "base64_encoded_signature"

		// Construct signed data
		signedData := strconv.Itoa(*thisDevice.SignatureCounter) + "_" + data + "_" + base64.StdEncoding.EncodeToString([]byte(thisDevice.ID))

		// Respond with signature response
		c.JSON(http.StatusOK, device.SignatureResponse{
			Signature:  signature,
			SignedData: signedData,
		})
	})

	// Run the server
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}

}
