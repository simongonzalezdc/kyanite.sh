package audio

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"
)

type SoundEffect int

const (
	SoundTaskAdd SoundEffect = iota
	SoundTaskComplete
	SoundTimerStart
	SoundTimerStop
	SoundError
	SoundSuccess
	SoundNavigate
	SoundGlitch
)

type AudioPlayer struct {
	enabled atomic.Bool
}

var globalPlayer *AudioPlayer

func init() {
	globalPlayer = NewAudioPlayer()
}

func NewAudioPlayer() *AudioPlayer {
	p := &AudioPlayer{}
	p.enabled.Store(true) // Enable by default
	return p
}

func (a *AudioPlayer) PlaySound(effect SoundEffect) {
	if !a.enabled.Load() {
		return
	}

	// Simple system beep for cross-platform compatibility
	go func() {
		// T4-06: prevent a panic in the sound effect from crashing the UI.
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "audio: recovered panic in PlaySound: %v\n", r)
			}
		}()
		switch effect {
		case SoundTaskAdd:
			a.playBeep(800, 100) // High beep for add
		case SoundTaskComplete:
			a.playBeep(1200, 150) // Higher beep for complete
		case SoundTimerStart:
			a.playBeep(600, 100) // Lower beep for start
		case SoundTimerStop:
			a.playBeep(400, 150) // Even lower for stop
		case SoundError:
			a.playBeep(300, 200) // Low beep for error
		case SoundSuccess:
			a.playBeep(1000, 200) // Success sound
		case SoundNavigate:
			a.playBeep(700, 50) // Quick navigation sound
		case SoundGlitch:
			// Removed decorative glitch effect
			a.playBeep(400, 100) // Simple notification sound
		}
	}()
}

func (a *AudioPlayer) playBeep(frequency, duration int) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-c",
			fmt.Sprintf("[console]::beep(%d,%d)", frequency, duration))
	case "darwin":
		cmd = exec.Command("afplay", "/System/Library/Sounds/Ping.aiff")
		if duration > 100 {
			cmd = exec.Command("afplay", "/System/Library/Sounds/Glass.aiff")
		}
	case "linux":
		// Try to use beep command first, fallback to echo
		cmd = exec.Command("beep", "-f", fmt.Sprintf("%d", frequency), "-l", fmt.Sprintf("%d", duration))
		if err := cmd.Run(); err != nil {
			// Fallback to terminal bell
			fmt.Print("\a")
		}
		return
	default:
		// Fallback to terminal bell
		fmt.Print("\a")
		return
	}

	_ = cmd.Run() // intentionally ignore error to avoid breaking UI on missing player
}

func (a *AudioPlayer) Enable() {
	a.enabled.Store(true)
}

func (a *AudioPlayer) Disable() {
	a.enabled.Store(false)
}

func (a *AudioPlayer) IsEnabled() bool {
	return a.enabled.Load()
}

// Global convenience functions
func PlaySound(effect SoundEffect) {
	globalPlayer.PlaySound(effect)
}

func EnableAudio() {
	globalPlayer.Enable()
}

func DisableAudio() {
	globalPlayer.Disable()
}

func IsAudioEnabled() bool {
	return globalPlayer.IsEnabled()
}
