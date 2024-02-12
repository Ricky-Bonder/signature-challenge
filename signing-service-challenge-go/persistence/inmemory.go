package persistence

import (
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"sync"
)

// Storage Implement this interface in any persistence layer, such as a DB
type Storage interface {
	GetSignatureDevice(id string) (*domain.InternalSignatureDevice, error)
	CreateSignatureDevice(device *domain.InternalSignatureDevice) error
	GetAllSignatureDevices() ([]*domain.InternalSignatureDevice, error)
}

var (
	once                   sync.Once
	singletonMemoryStorage *MemoryStorage
)

type MemoryStorage struct {
	devices map[string]*domain.InternalSignatureDevice
	mutex   sync.RWMutex
}

// NewMemoryStorage creates a new instance of MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		devices: make(map[string]*domain.InternalSignatureDevice),
	}
}

func GetMemoryStorage() *MemoryStorage {
	once.Do(func() {
		singletonMemoryStorage = NewMemoryStorage()
	})
	return singletonMemoryStorage
}

// GetSignatureDevice retrieves a signature device by ID from memory storage.
func (m *MemoryStorage) GetSignatureDevice(id string) (*domain.InternalSignatureDevice, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	device, ok := m.devices[id]
	if !ok {
		return nil, errors.New("device not found")
	}
	return device, nil
}

// CreateSignatureDevice creates a signature device in memory storage.
func (m *MemoryStorage) CreateSignatureDevice(device *domain.InternalSignatureDevice) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.devices[device.ID]; exists {
		return errors.New("device with the same ID already exists")
	}
	m.devices[device.ID] = device
	return nil
}

// GetAllSignatureDevices retrieves all signature devices from memory storage.
func (m *MemoryStorage) GetAllSignatureDevices() ([]*domain.InternalSignatureDevice, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	devices := make([]*domain.InternalSignatureDevice, 0, len(m.devices))
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices, nil
}
