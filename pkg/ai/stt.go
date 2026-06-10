package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// STTClient handles speech-to-text via local whisper.cpp.
type STTClient struct {
	bin    string
	model  string
	lang   string
}

// NewSTTClient creates a new STT client using local whisper.cpp.
func NewSTTClient(bin, model, lang string) *STTClient {
	return &STTClient{
		bin:   bin,
		model: model,
		lang:  lang,
	}
}

// TranscribeFile transcribes an audio file to text.
// The file must be in WAV format (16-bit PCM, 16kHz mono recommended).
func (c *STTClient) TranscribeFile(ctx context.Context, audioPath string) (*Transcription, error) {
	start := time.Now()

	args := []string{
		"-m", c.model,
		"-f", audioPath,
		"--output-json",
		"--no-timestamps",
	}
	if c.lang != "" {
		args = append(args, "-l", c.lang)
	}

	cmd := exec.CommandContext(ctx, c.bin, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %v\nstderr: %s", ErrSTTFailed, err, stderr.String())
	}

	text := parseWhisperOutput(stdout.String())

	return &Transcription{
		Text:     strings.TrimSpace(text),
		Language: c.lang,
		Duration: time.Since(start),
	}, nil
}

// TranscribePCM transcribes raw PCM audio data by piping to whisper-stream.
// The input must be 16-bit PCM, 16kHz mono.
func (c *STTClient) TranscribePCM(ctx context.Context, pcmData []byte) (*Transcription, error) {
	start := time.Now()

	args := []string{
		"-m", c.model,
		"--step", "0",
		"--length", "30000",
	}
	if c.lang != "" {
		args = append(args, "-l", c.lang)
	}

	cmd := exec.CommandContext(ctx, c.bin, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewReader(pcmData)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %v\nstderr: %s", ErrSTTFailed, err, stderr.String())
	}

	text := parseWhisperOutput(stdout.String())

	return &Transcription{
		Text:     strings.TrimSpace(text),
		Language: c.lang,
		Duration: time.Since(start),
	}, nil
}

// IsAvailable checks if whisper.cpp is installed and the model exists.
func (c *STTClient) IsAvailable() error {
	_, err := exec.LookPath(c.bin)
	if err != nil {
		return fmt.Errorf("whisper binary %q not found in PATH: %w", c.bin, err)
	}
	if _, err := os.Stat(c.model); err != nil {
		return fmt.Errorf("%w: model file %q not found: %v", ErrSTTNotInstalled, c.model, err)
	}
	return nil
}

// whisperJSONOutput matches the JSON output format from whisper.cpp --output-json.
type whisperJSONOutput struct {
	Transcription []struct {
		Timestamps struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"timestamps"`
		Text string `json:"text"`
	} `json:"transcription"`
}

// parseWhisperOutput extracts the transcription text from whisper's output.
// It tries JSON first, then falls back to plain text.
func parseWhisperOutput(output string) string {
	output = strings.TrimSpace(output)

	// Try JSON parse
	var jsonOut whisperJSONOutput
	if err := json.Unmarshal([]byte(output), &jsonOut); err == nil {
		var parts []string
		for _, seg := range jsonOut.Transcription {
			if seg.Text != "" {
				parts = append(parts, strings.TrimSpace(seg.Text))
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, " ")
		}
	}

	// Fallback: whisper-cli with --output-json writes JSON to a sidecar file
	// and human-readable text to stdout. If JSON parse failed, treat as plain text.
	lines := strings.Split(output, "\n")
	var textLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Strip whisper timestamp prefix: [HH:MM:SS.mmm --> HH:MM:SS.mmm]
		if idx := strings.Index(line, "]"); idx >= 0 && strings.HasPrefix(line, "[") {
			line = strings.TrimSpace(line[idx+1:])
		}
		if line != "" {
			textLines = append(textLines, line)
		}
	}
	return strings.Join(textLines, " ")
}
