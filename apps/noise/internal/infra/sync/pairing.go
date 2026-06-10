package sync

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/logging"
)

const (
	// PairingCodeExpiry is how long a pairing code is valid
	PairingCodeExpiry = 5 * time.Minute
	// PairingCodeLength is the length of the pairing code
	PairingCodeLength = 6
)

// PairedDevice represents a paired device
type PairedDevice struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	LastSeen time.Time `json:"last_seen"`
	PairedAt time.Time `json:"paired_at"`
}

// PairingManager handles device pairing
type PairingManager struct {
	activeCodes   map[string]time.Time // code -> expiry
	pairedDevices map[string]*PairedDevice
	logger        *logging.Logger
	mu            sync.RWMutex
}

// NewPairingManager creates a new pairing manager
func NewPairingManager(logger *logging.Logger) *PairingManager {
	if logger == nil {
		logger = logging.GetDefaultLogger()
	}
	pm := &PairingManager{
		activeCodes:   make(map[string]time.Time),
		pairedDevices: make(map[string]*PairedDevice),
		logger:        logger,
	}

	// Start cleanup goroutine
	go pm.cleanupExpiredCodes()

	return pm
}

// GenerateCode generates a new pairing code
func (pm *PairingManager) GenerateCode() string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Generate random 6-digit code
	code := generateNumericCode(PairingCodeLength)

	// Ensure uniqueness - regenerate if code exists and hasn't expired
	for {
		expiry, exists := pm.activeCodes[code]
		if !exists || !expiry.After(time.Now()) {
			break
		}
		code = generateNumericCode(PairingCodeLength)
	}

	pm.activeCodes[code] = time.Now().Add(PairingCodeExpiry)
	pm.logger.Infof("Generated pairing code: %s (expires in %v)", code, PairingCodeExpiry)

	return code
}

// ValidateCode checks if a pairing code is valid
func (pm *PairingManager) ValidateCode(code string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	expiry, exists := pm.activeCodes[code]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(pm.activeCodes, code)
		return false
	}

	// Code is valid, remove it (one-time use)
	delete(pm.activeCodes, code)
	return true
}

// PairDevice creates a new paired device entry
func (pm *PairingManager) PairDevice(name string) *PairedDevice {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	device := &PairedDevice{
		ID:       generateDeviceID(),
		Name:     name,
		LastSeen: time.Now(),
		PairedAt: time.Now(),
	}

	pm.pairedDevices[device.ID] = device
	pm.logger.Infof("Device paired: %s (%s)", device.Name, device.ID)

	return device
}

// GetDevice returns a paired device by ID
func (pm *PairingManager) GetDevice(id string) *PairedDevice {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.pairedDevices[id]
}

// UpdateLastSeen updates the last seen time for a device
func (pm *PairingManager) UpdateLastSeen(id string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if device, ok := pm.pairedDevices[id]; ok {
		device.LastSeen = time.Now()
	}
}

// UnpairDevice removes a paired device
func (pm *PairingManager) UnpairDevice(id string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, ok := pm.pairedDevices[id]; ok {
		delete(pm.pairedDevices, id)
		pm.logger.Infof("Device unpaired: %s", id)
		return true
	}
	return false
}

// GetPairedDevices returns all paired devices
func (pm *PairingManager) GetPairedDevices() []*PairedDevice {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	devices := make([]*PairedDevice, 0, len(pm.pairedDevices))
	for _, device := range pm.pairedDevices {
		devices = append(devices, device)
	}
	return devices
}

// IsDevicePaired checks if a device ID is paired
func (pm *PairingManager) IsDevicePaired(id string) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	_, ok := pm.pairedDevices[id]
	return ok
}

// cleanupExpiredCodes periodically removes expired pairing codes
func (pm *PairingManager) cleanupExpiredCodes() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		pm.mu.Lock()
		now := time.Now()
		for code, expiry := range pm.activeCodes {
			if now.After(expiry) {
				delete(pm.activeCodes, code)
			}
		}
		pm.mu.Unlock()
	}
}

// generateNumericCode generates a random numeric code of the specified length
func generateNumericCode(length int) string {
	const digits = "0123456789"
	code := make([]byte, length)

	// Read random bytes
	if _, err := rand.Read(code); err != nil {
		// Fallback to time-based if crypto/rand fails
		code = []byte(fmt.Sprintf("%06d", time.Now().UnixNano()%1000000))
	} else {
		for i := range code {
			code[i] = digits[code[i]%10]
		}
	}

	return string(code)
}

// generateDeviceID generates a unique device ID
func generateDeviceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("device_%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", bytes)
}
