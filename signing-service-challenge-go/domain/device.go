package domain

import (
	"go.uber.org/zap"
	"sync"
)

// Algorithm represents the supported algorithms.
type Algorithm int

// Constants representing supported algorithms.
const (
	ECC Algorithm = iota + 1
	RSA
)

//// Validate checks if the given algorithm is supported.
//func (a Algorithm) Validate() (crypto.ECCKeyPair, crypto.RSAKeyPair, error) {
//	switch a {
//	case ECC:
//		return crypto.ECCKeyPair{}, nil, nil
//	case RSA:
//		return crypto.RSAKeyPair, nil
//	default:
//		return nil, errors.New("unsupported algorithm")
//	}
//}
//
//func (a Algorithm) ValidateRSA() (crypto.RSAGenerator, error) {
//	switch a {
//	case RSA:
//		return crypto.RSAKeyPair{}, nil
//	default:
//		return nil, errors.New("unsupported algorithm")
//	}
//}

type SignatureDevice struct {
	ID               string    `json:"id"`
	Algorithm        Algorithm `json:"algorithm"`
	Label            *string   `json:"label"`
	SignatureCounter *int      `json:"signatureCounter"`
	Logger           *zap.SugaredLogger
}

type CreateSignatureDeviceResponse struct {
	// Define your response fields here
}

type SignatureService struct {
	Devices map[string]*SignatureDevice
	Mutex   sync.Mutex
}

func NewSignatureService() *SignatureService {
	return &SignatureService{
		Devices: make(map[string]*SignatureDevice),
	}
}

func NewDeviceInstance(logger *zap.SugaredLogger) *SignatureDevice {

	return &SignatureDevice{
		Logger: logger,
	}
}

// CreateSignatureDeviceRequest represents the request body for creating a signature device.
type CreateSignatureDeviceRequest struct {
	Algorithm Algorithm `json:"algorithm"`
	Label     *string   `json:"label,omitempty"`
}

// SignatureResponse represents the response body for signing a transaction.
type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}
