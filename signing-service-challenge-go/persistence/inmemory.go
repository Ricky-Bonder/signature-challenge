package persistence

import (
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"sync"
)

// DevicesStorage Implement this interface in any persistence layer, such as a DB
type DevicesStorage interface {
	GetSignatureDevice(id string) (*domain.InternalSignatureDevice, error)
	CreateSignatureDevice(device *domain.InternalSignatureDevice) error
	GetAllSignatureDevices() ([]*domain.InternalSignatureDevice, error)
}

type SignaturesStorage interface {
	GetLastSignature() (string, error)
	InsertSignature(signature *domain.SignatureResponse)
}

var (
	once                      sync.Once
	singletonMemoryStorage    *DeviceStorage
	singletonSignatureStorage *SignatureStorage
)

type DeviceStorage struct {
	devices map[string]*domain.InternalSignatureDevice
	mutex   sync.RWMutex
}

type SignatureStorage struct {
	signatures     []*domain.SignatureResponse
	lastInserted   string
	signatureMutex sync.RWMutex
}

// NewSignatureDeviceStorage creates a new instance of DeviceStorage.
func NewSignatureDeviceStorage() *DeviceStorage {
	return &DeviceStorage{
		devices: make(map[string]*domain.InternalSignatureDevice),
	}
}

func NewSignatureStorage() *SignatureStorage {
	return &SignatureStorage{
		signatures: make([]*domain.SignatureResponse, 0),
	}
}

func GetSignatureDeviceStorage() *DeviceStorage {
	once.Do(func() {
		singletonMemoryStorage = NewSignatureDeviceStorage()
	})
	return singletonMemoryStorage
}

func GetSignatureStorage() *SignatureStorage {
	once.Do(func() {
		singletonSignatureStorage = NewSignatureStorage()
	})
	return singletonSignatureStorage
}

// GetLastSignature retrieves the last inserted signature from memory storage.
func (s *SignatureStorage) GetLastSignature() (string, error) {
	s.signatureMutex.RLock()
	defer s.signatureMutex.RUnlock()

	if len(s.signatures) == 0 {
		return "0", errors.New("no devices found")
	}

	return s.lastInserted, nil
}

// InsertSignature creates a signature in memory storage.
func (s *SignatureStorage) InsertSignature(signature *domain.SignatureResponse) {
	//FIXME
	s.signatureMutex.Lock()
	defer s.signatureMutex.Unlock()

	s.signatures = append(s.signatures, signature)
	s.lastInserted = signature.Signature
}

// GetSignatureDevice retrieves a signature device by ID from memory storage.
func (m *DeviceStorage) GetSignatureDevice(id string) (*domain.InternalSignatureDevice, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	device, ok := m.devices[id]
	if !ok {
		return nil, errors.New("device not found")
	}
	return device, nil
}

// CreateSignatureDevice creates a signature device in memory storage.
func (m *DeviceStorage) CreateSignatureDevice(device *domain.InternalSignatureDevice) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.devices[device.ID]; exists {
		return errors.New("device with the same ID already exists")
	}
	m.devices[device.ID] = device
	return nil
}

// GetAllSignatureDevices retrieves all signature devices from memory storage.
func (m *DeviceStorage) GetAllSignatureDevices() ([]*domain.InternalSignatureDevice, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	devices := make([]*domain.InternalSignatureDevice, 0, len(m.devices))
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices, nil
}
