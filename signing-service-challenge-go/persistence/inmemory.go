package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

type Storage interface {
}

type MemoryStorage struct {
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (m *MemoryStorage) GetSignatureDevice(id string) (*domain.SignatureDevice, error) {

	return nil, nil
}

func (m *MemoryStorage) CreateSignatureDevice(device *domain.SignatureDevice) error {
	return nil
}

func (m *MemoryStorage) GetAllSignatureDevices() ([]*domain.SignatureDevice, error) {
	return nil, nil
}
