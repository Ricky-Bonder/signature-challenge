package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// RSASigner is a signer for RSA.
type RSASigner struct {
	KeyPair RSAKeyPair
}

// NewRSASigner creates a new RSASigner.
func NewRSASigner(keyPair RSAKeyPair) *RSASigner {
	return &RSASigner{KeyPair: keyPair}
}

// Sign signs the data using RSA.
func SignRSA(keypair *RSAKeyPair, dataToBeSigned []byte) ([]byte, error) {
	hashed := crypto.SHA256.New()
	_, err := hashed.Write(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	hashedData := hashed.Sum(nil)
	if err != nil {
		panic(err)
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, keypair.Private, crypto.SHA256, hashedData)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}
	return signature, nil
}

// ECDSASigner is a signer for ECDSA.
type ECDSASigner struct {
	keyPair ECCKeyPair
}

// NewECDSASigner creates a new ECDSASigner.
func NewECDSASigner(keyPair ECCKeyPair) *ECDSASigner {
	return &ECDSASigner{keyPair: keyPair}
}

// Sign signs the data using ECDSA.
func SignECC(keypair *ECCKeyPair, dataToBeSigned []byte) ([]byte, error) {
	hashed := crypto.SHA256.New()
	hashed.Write(dataToBeSigned)
	hashedData := hashed.Sum(nil)

	signature, err := ecdsa.SignASN1(rand.Reader, keypair.Private, hashedData)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}
	return signature, nil
}
