package bootstrap

import (
	"testing"

	"github.com/kyanite/appnames"
)

func TestInitReturnsAppName(t *testing.T) {
	// Init will fail (no config file, no brain) but should not panic.
	app := Init(appnames.Focus)
	if app == nil {
		t.Fatal("expected non-nil App")
	}
	if app.AppName != appnames.Focus {
		t.Errorf("expected app name %q, got %q", appnames.Focus, app.AppName)
	}
	if app.SessionID == "" {
		t.Error("expected non-empty session ID")
	}
}

func TestOthersExcludesSelf(t *testing.T) {
	app := Init(appnames.Prism)
	others := app.Others()
	for _, o := range others {
		if o == appnames.Prism {
			t.Error("Others should not contain self")
		}
	}
}

func TestShutdownNilSafe(t *testing.T) {
	var app *App
	app.Shutdown("test", nil) // should not panic
}

func TestSaveCrossAppContextNilSafe(t *testing.T) {
	var app *App
	app.SaveCrossAppContext(nil, "test", "test", 0.5) // should not panic
}
