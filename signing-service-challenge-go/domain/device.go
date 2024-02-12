package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"sync"
)

// Algorithm represents the supported algorithms.
type Algorithm int

// Constants representing supported algorithms.
const (
	ECC Algorithm = iota + 1
	RSA
)

var (
	onceSignatureService      sync.Once
	singletonSignatureService *SignatureService
)

// AlgorithmNames maps algorithm strings to their corresponding Algorithm constants.
var AlgorithmNames = map[string]Algorithm{
	"ECC": ECC,
	"RSA": RSA,
}

// CreateSignatureDeviceRequest represents the request body for creating a signature device.
type CreateSignatureDeviceRequest struct {
	Algorithm string  `json:"algorithm"`
	Label     *string `json:"label"`
}

type InternalSignatureDevice struct {
	ID               string                  `json:"id"`
	Algorithm        crypto.KeyPairGenerator `json:"algorithm"`
	Label            *string                 `json:"label"`
	SignatureCounter int32                   `json:"signatureCounter"`
}

type CreateSignatureDeviceResponse struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type SignatureService struct {
	Devices map[string]*InternalSignatureDevice
	Mutex   sync.Mutex
}

func NewSignatureService() *SignatureService {
	return &SignatureService{
		Devices: make(map[string]*InternalSignatureDevice),
	}
}

func GetSignatureService() *SignatureService {
	onceSignatureService.Do(func() {
		singletonSignatureService = NewSignatureService()
	})
	return singletonSignatureService
}

type SignTransactionRequest struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}
