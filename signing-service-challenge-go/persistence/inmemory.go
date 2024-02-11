package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

type Storage interface {
	GetSignatureDevice(id string) (*domain.InternalSignatureDevice, error)
	CreateSignatureDevice(device *domain.InternalSignatureDevice) error
	GetAllSignatureDevices() ([]*domain.InternalSignatureDevice, error)
}

type MemoryStorage struct {
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (m *MemoryStorage) GetSignatureDevice(id string) (*domain.InternalSignatureDevice, error) {

	return nil, nil
}

func (m *MemoryStorage) CreateSignatureDevice(device *domain.InternalSignatureDevice) error {
	return nil
}

func (m *MemoryStorage) GetAllSignatureDevices() ([]*domain.InternalSignatureDevice, error) {
	return nil, nil
}
