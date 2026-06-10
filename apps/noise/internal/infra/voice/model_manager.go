// Deprecated: Model management is now handled by pkg/ai STTClient configuration. Use infra/brain/client.go instead.
package voice

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/kyanite/noise/internal/logging"
)

// Common errors for model management
var (
	ErrModelNotFound    = errors.New("model not found")
	ErrDownloadFailed   = errors.New("model download failed")
	ErrChecksumMismatch = errors.New("model checksum mismatch")
	ErrInvalidModel     = errors.New("invalid model name")
)

// Model constants
const (
	// Model names
	ModelTinyEN  = "ggml-tiny.en.bin"  // ~75MB, fastest, English only
	ModelBaseEN  = "ggml-base.en.bin"  // ~142MB, recommended balance
	ModelSmallEN = "ggml-small.en.bin" // ~466MB, best accuracy for English

	// Multilingual models (larger, support all languages)
	ModelTiny  = "ggml-tiny.bin"  // ~75MB
	ModelBase  = "ggml-base.bin"  // ~142MB
	ModelSmall = "ggml-small.bin" // ~466MB

	// Download URLs
	baseDownloadURL = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/"
)

// ModelInfo contains information about a whisper model
type ModelInfo struct {
	Name         string
	Size         int64  // Size in bytes
	SHA256       string // Expected checksum
	Description  string
	Multilingual bool
}

// AvailableModels returns information about all available models
var AvailableModels = map[string]ModelInfo{
	ModelTinyEN: {
		Name:         ModelTinyEN,
		Size:         75000000,
		SHA256:       "", // Will be verified after download
		Description:  "Tiny English-only model (~75MB) - Fastest, good for quick dictation",
		Multilingual: false,
	},
	ModelBaseEN: {
		Name:         ModelBaseEN,
		Size:         142000000,
		SHA256:       "",
		Description:  "Base English-only model (~142MB) - Recommended balance of speed/accuracy",
		Multilingual: false,
	},
	ModelSmallEN: {
		Name:         ModelSmallEN,
		Size:         466000000,
		SHA256:       "",
		Description:  "Small English-only model (~466MB) - Best accuracy for English",
		Multilingual: false,
	},
	ModelTiny: {
		Name:         ModelTiny,
		Size:         75000000,
		SHA256:       "",
		Description:  "Tiny multilingual model (~75MB) - Fastest, supports 99 languages",
		Multilingual: true,
	},
	ModelBase: {
		Name:         ModelBase,
		Size:         142000000,
		SHA256:       "",
		Description:  "Base multilingual model (~142MB) - Good balance, 99 languages",
		Multilingual: true,
	},
	ModelSmall: {
		Name:         ModelSmall,
		Size:         466000000,
		SHA256:       "",
		Description:  "Small multilingual model (~466MB) - Best accuracy, 99 languages",
		Multilingual: true,
	},
}

// ProgressCallback is called during download with progress updates
type ProgressCallback func(downloaded, total int64)

// ModelManager handles whisper model download and management
type ModelManager struct {
	modelsDir string
	logger    *logging.Logger
	mu        sync.Mutex

	// Active download tracking
	downloading map[string]bool
	downloadMu  sync.Mutex
}

// NewModelManager creates a new model manager
func NewModelManager(modelsDir string, logger *logging.Logger) (*ModelManager, error) {
	if logger == nil {
		logger = logging.GetDefaultLogger()
	}

	// Create models directory if it doesn't exist
	if err := os.MkdirAll(modelsDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create models directory: %w", err)
	}

	return &ModelManager{
		modelsDir:   modelsDir,
		logger:      logger,
		downloading: make(map[string]bool),
	}, nil
}

// NewModelManagerWithAutoSetup creates a model manager and ensures data directories exist
func NewModelManagerWithAutoSetup(dataDir string, logger *logging.Logger) (*ModelManager, error) {
	modelsDir := filepath.Join(dataDir, "models")
	return NewModelManager(modelsDir, logger)
}

// GetModelPath returns the full path to a model file
func (mm *ModelManager) GetModelPath(modelName string) string {
	return filepath.Join(mm.modelsDir, modelName)
}

// IsModelAvailable checks if a model is downloaded and available
func (mm *ModelManager) IsModelAvailable(modelName string) bool {
	path := mm.GetModelPath(modelName)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// Check if file is not empty
	return info.Size() > 0
}

// GetAvailableModels returns a list of downloaded models
func (mm *ModelManager) GetAvailableModels() []string {
	var models []string
	for name := range AvailableModels {
		if mm.IsModelAvailable(name) {
			models = append(models, name)
		}
	}
	return models
}

// EnsureModel ensures a model is available, downloading if necessary
func (mm *ModelManager) EnsureModel(modelName string, progress ProgressCallback) (string, error) {
	// Validate model name
	if _, ok := AvailableModels[modelName]; !ok {
		return "", fmt.Errorf("%w: %s", ErrInvalidModel, modelName)
	}

	path := mm.GetModelPath(modelName)

	// Check if already downloaded
	if mm.IsModelAvailable(modelName) {
		mm.logger.Infof("Model %s already available at %s", modelName, path)
		return path, nil
	}

	// Download the model
	if err := mm.DownloadModel(modelName, progress); err != nil {
		return "", err
	}

	return path, nil
}

// DownloadModel downloads a model from Hugging Face
func (mm *ModelManager) DownloadModel(modelName string, progress ProgressCallback) error {
	// Check if already downloading
	mm.downloadMu.Lock()
	if mm.downloading[modelName] {
		mm.downloadMu.Unlock()
		return fmt.Errorf("model %s is already being downloaded", modelName)
	}
	mm.downloading[modelName] = true
	mm.downloadMu.Unlock()

	defer func() {
		mm.downloadMu.Lock()
		delete(mm.downloading, modelName)
		mm.downloadMu.Unlock()
	}()

	// Validate model name
	modelInfo, ok := AvailableModels[modelName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrInvalidModel, modelName)
	}

	url := baseDownloadURL + modelName
	destPath := mm.GetModelPath(modelName)
	tempPath := destPath + ".download"

	mm.logger.Infof("Downloading model %s from %s", modelName, url)

	// Create HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: HTTP %d", ErrDownloadFailed, resp.StatusCode)
	}

	// Create temp file
	out, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		out.Close()
		// Clean up temp file on error
		if _, err := os.Stat(tempPath); err == nil {
			os.Remove(tempPath)
		}
	}()

	// Get content length for progress
	contentLength := resp.ContentLength
	if contentLength <= 0 {
		contentLength = modelInfo.Size // Use estimated size
	}

	// Download with progress tracking
	hasher := sha256.New()
	var downloaded int64

	reader := resp.Body
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			// Write to file
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to file: %w", writeErr)
			}
			// Update hash
			hasher.Write(buf[:n])

			downloaded += int64(n)

			// Report progress
			if progress != nil {
				progress(downloaded, contentLength)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%w: %v", ErrDownloadFailed, err)
		}
	}

	// Close file before rename
	out.Close()

	// Verify checksum if available
	if modelInfo.SHA256 != "" {
		actualHash := hex.EncodeToString(hasher.Sum(nil))
		if actualHash != modelInfo.SHA256 {
			os.Remove(tempPath)
			return fmt.Errorf("%w: expected %s, got %s", ErrChecksumMismatch, modelInfo.SHA256, actualHash)
		}
	}

	// Rename temp file to final destination
	if err := os.Rename(tempPath, destPath); err != nil {
		return fmt.Errorf("failed to rename downloaded file: %w", err)
	}

	mm.logger.Infof("Model %s downloaded successfully to %s", modelName, destPath)
	return nil
}

// DeleteModel removes a downloaded model
func (mm *ModelManager) DeleteModel(modelName string) error {
	path := mm.GetModelPath(modelName)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete model: %w", err)
	}
	mm.logger.Infof("Deleted model %s", modelName)
	return nil
}

// GetModelInfo returns information about a model
func (mm *ModelManager) GetModelInfo(modelName string) (*ModelInfo, error) {
	info, ok := AvailableModels[modelName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidModel, modelName)
	}

	// Update with actual file size if downloaded
	if mm.IsModelAvailable(modelName) {
		fileInfo, err := os.Stat(mm.GetModelPath(modelName))
		if err == nil {
			info.Size = fileInfo.Size()
		}
	}

	return &info, nil
}

// GetModelsDir returns the models directory path
func (mm *ModelManager) GetModelsDir() string {
	return mm.modelsDir
}

// IsDownloading returns whether a model is currently being downloaded
func (mm *ModelManager) IsDownloading(modelName string) bool {
	mm.downloadMu.Lock()
	defer mm.downloadMu.Unlock()
	return mm.downloading[modelName]
}

// FormatModelSize returns a human-readable size string
func FormatModelSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// RecommendedModel returns the recommended model for a given use case
func RecommendedModel(englishOnly bool, preferSpeed bool) string {
	if englishOnly {
		if preferSpeed {
			return ModelTinyEN
		}
		return ModelBaseEN // Best balance for English
	}
	if preferSpeed {
		return ModelTiny
	}
	return ModelBase // Best balance for multilingual
}
