package main

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"go.uber.org/zap"
	"log"
)

const (
	ServerURL     = "http://localhost"
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
	storage := persistence.GetSignatureDeviceStorage() // For further implementation change this to DB storage
	signatureStorage := persistence.GetSignatureStorage()
	server := api.NewServer(ServerURL, ListenAddress, storage, signatureStorage)
	domain.NewSignatureService()

	// Run the server
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}

}
