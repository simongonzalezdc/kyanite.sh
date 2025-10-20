package audio

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
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
	enabled bool
}

var globalPlayer *AudioPlayer

func init() {
	globalPlayer = NewAudioPlayer()
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		enabled: true, // Enable by default
	}
}

func (a *AudioPlayer) PlaySound(effect SoundEffect) {
	if !a.enabled {
		return
	}

	// Simple system beep for cross-platform compatibility
	go func() {
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
			a.playGlitchEffect() // Special glitch sound
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
	
	_ = cmd.Run() // Ignore errors to not break the UI
}

func (a *AudioPlayer) playGlitchEffect() {
	// Rapid fire beeps for glitch effect
	go func() {
		for i := 0; i < 5; i++ {
			a.playBeep(200+i*100, 30)
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func (a *AudioPlayer) Enable() {
	a.enabled = true
}

func (a *AudioPlayer) Disable() {
	a.enabled = false
}

func (a *AudioPlayer) IsEnabled() bool {
	return a.enabled
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
