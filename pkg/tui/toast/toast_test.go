package toast

import (
	"testing"
	"time"

	"github.com/kyanite/design"
)

func TestNewModelNoPanic(t *testing.T) {
	dm := design.NewManager("amber-night")
	m := New(dm)
	if m == nil {
		t.Fatal("expected non-nil Model")
	}
}

func TestModelEmptyView(t *testing.T) {
	dm := design.NewManager("amber-night")
	m := New(dm)
	if m.View() != "" {
		t.Error("empty model should render empty string")
	}
}

func TestAddToast(t *testing.T) {
	dm := design.NewManager("amber-night")
	m := New(dm)
	cmd := m.Add(SuccessMsg("test"))
	if cmd == nil {
		t.Fatal("expected non-nil command")
	}
	// Simulate the Msg arriving
	_ = m.Update(cmd())
	if !m.HasToasts() {
		t.Error("expected toasts after add")
	}
}

func TestClearAll(t *testing.T) {
	dm := design.NewManager("amber-night")
	m := New(dm)
	_ = m.Update(m.Add(InfoMsg("one")))
	_ = m.Update(m.Add(InfoMsg("two")))
	m.ClearAll()
	if m.HasToasts() {
		t.Error("expected no toasts after ClearAll")
	}
}

func TestDismissRemovesToast(t *testing.T) {
	dm := design.NewManager("amber-night")
	m := New(dm)
	cmd := m.Add(WarningMsg("fade"))
	_ = m.Update(cmd()) // call the cmd to produce the Msg, then update

	// Find the item ID and dismiss it
	if len(m.items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(m.items))
	}
	id := m.items[0].ID
	_ = m.Update(dismissMsg{id: id})

	if m.HasToasts() {
		t.Error("expected toast to be dismissed")
	}
}

func TestDurationConstants(t *testing.T) {
	if DurationDefault != 3*time.Second {
		t.Errorf("DurationDefault = %v, want 3s", DurationDefault)
	}
	if DurationLong != 5*time.Second {
		t.Errorf("DurationLong = %v, want 5s", DurationLong)
	}
	if DurationShort != 2*time.Second {
		t.Errorf("DurationShort = %v, want 2s", DurationShort)
	}
}
