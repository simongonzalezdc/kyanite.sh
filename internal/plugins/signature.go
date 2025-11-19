package plugins

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// SignatureVerifier handles plugin signature verification
type SignatureVerifier struct {
	trustedKeys map[string]*rsa.PublicKey
	logger      Logger
}

// NewSignatureVerifier creates a new signature verifier
func NewSignatureVerifier(logger Logger) *SignatureVerifier {
	return &SignatureVerifier{
		trustedKeys: make(map[string]*rsa.PublicKey),
		logger:      logger,
	}
}

// LoadTrustedKey loads a trusted public key
func (v *SignatureVerifier) LoadTrustedKey(keyID string, keyPath string) error {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("failed to read key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("key is not an RSA public key")
	}

	v.trustedKeys[keyID] = rsaKey
	v.logger.Debugf("Loaded trusted key: %s", keyID)
	return nil
}

// LoadTrustedKeysFromDir loads all trusted keys from a directory
func (v *SignatureVerifier) LoadTrustedKeysFromDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			v.logger.Debugf("Trusted keys directory does not exist: %s", dir)
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".pem" {
			continue
		}

		keyID := entry.Name()[:len(entry.Name())-4]
		keyPath := filepath.Join(dir, entry.Name())

		if err := v.LoadTrustedKey(keyID, keyPath); err != nil {
			v.logger.Warnf("Failed to load key %s: %v", keyID, err)
			continue
		}
	}

	v.logger.Infof("Loaded %d trusted keys", len(v.trustedKeys))
	return nil
}

// VerifyPlugin verifies a plugin's signature
func (v *SignatureVerifier) VerifyPlugin(pluginPath string) error {
	// 1. Calculate plugin file hash
	hash, err := v.calculatePluginHash(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to calculate plugin hash: %w", err)
	}

	v.logger.Debugf("Plugin hash: %x", hash)

	// 2. Load signature file
	sigPath := pluginPath + ".sig"
	signature, keyID, err := v.loadSignature(sigPath)
	if err != nil {
		return fmt.Errorf("failed to load signature: %w", err)
	}

	v.logger.Debugf("Signature loaded for key: %s", keyID)

	// 3. Get trusted public key
	pubKey, ok := v.trustedKeys[keyID]
	if !ok {
		return fmt.Errorf("unknown signing key: %s", keyID)
	}

	// 4. Verify signature
	if err := v.verifySignature(pubKey, hash, signature); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	v.logger.Infof("Plugin signature verified successfully: %s", pluginPath)
	return nil
}

// calculatePluginHash calculates SHA-256 hash of plugin file
func (v *SignatureVerifier) calculatePluginHash(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	return hash[:], nil
}

// loadSignature loads signature and key ID from signature file
func (v *SignatureVerifier) loadSignature(path string) ([]byte, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	// Signature file format: keyID:base64signature
	// For simplicity, we expect PEM-encoded signature with key ID in header
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, "", fmt.Errorf("failed to decode signature PEM")
	}

	keyID, ok := block.Headers["KeyID"]
	if !ok {
		return nil, "", fmt.Errorf("signature missing KeyID header")
	}

	return block.Bytes, keyID, nil
}

// verifySignature verifies RSA signature
func (v *SignatureVerifier) verifySignature(pubKey *rsa.PublicKey, hash, signature []byte) error {
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash, signature)
}

// Helper to update SecurityManager.VerifyPluginSignature
func (sm *SecurityManager) VerifyPluginSignature(path string) error {
	// Initialize verifier if not already done
	verifier := NewSignatureVerifier(sm.logger)

	// Load trusted keys from config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		sm.logger.Warnf("Cannot get home directory: %v", err)
		return fmt.Errorf("cannot verify signature: %w", err)
	}

	keysDir := filepath.Join(homeDir, ".config", "noise", "plugin-keys")
	if err := verifier.LoadTrustedKeysFromDir(keysDir); err != nil {
		sm.logger.Warnf("Cannot load trusted keys: %v", err)
	}

	// If no trusted keys, skip verification but log warning
	if len(verifier.trustedKeys) == 0 {
		sm.logger.Warnf("No trusted keys found - skipping signature verification for: %s", path)
		sm.logger.Info("To enable signature verification, add public keys to: " + keysDir)
		return nil
	}

	// Verify signature
	return verifier.VerifyPlugin(path)
}
