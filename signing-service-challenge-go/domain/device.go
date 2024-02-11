package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

// Algorithm represents the supported algorithms.
type Algorithm int

// Constants representing supported algorithms.
const (
	ECC Algorithm = iota + 1
	RSA
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

func NewDeviceInstance(logger *zap.SugaredLogger) *InternalSignatureDevice {

	return &InternalSignatureDevice{}
}

// SignatureResponse represents the response body for signing a transaction.
type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type SignatureCounter int32

func (c *SignatureCounter) Increment() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *SignatureCounter) Get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
