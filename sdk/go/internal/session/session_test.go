package session

import (
	"sync"
	"testing"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session/state"
)

func TestSession_New(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())

	session := New(id, pageID)

	if session.ID != id {
		t.Errorf("session.ID = %v, want %v", session.ID, id)
	}
	if session.PageID != pageID {
		t.Errorf("session.PageID = %v, want %v", session.PageID, pageID)
	}
	if session.State == nil {
		t.Error("session.State is nil")
	}
}

func TestSessionManager_GetSetDelete(t *testing.T) {
	manager := NewSessionManager()
	id := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())
	session := New(id, pageID)

	// Test SetSession
	manager.SetSession(session)

	// Test GetSession
	got := manager.GetSession(id)
	if got != session {
		t.Errorf("GetSession(%v) = %v, want %v", id, got, session)
	}

	// Test DisconnectSession
	manager.DisconnectSession(id)
	got = manager.GetSession(id)
	if got != nil {
		t.Errorf("GetSession(%v) after delete = %v, want nil", id, got)
	}
}

func TestSessionManager_ConcurrentAccess(t *testing.T) {
	manager := NewSessionManager()
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := uuid.Must(uuid.NewV4())
			pageID := uuid.Must(uuid.NewV4())
			session := New(id, pageID)

			manager.SetSession(session)
			got := manager.GetSession(id)
			if got != session {
				t.Errorf("GetSession(%v) = %v, want %v", id, got, session)
			}
			manager.DisconnectSession(id)
		}()
	}

	wg.Wait()
}

func TestState_SetGet(t *testing.T) {
	s := newState()
	id := uuid.Must(uuid.NewV4())

	radioState := &state.RadioState{
		ID:      id,
		Label:   "Test Radio",
		Options: []string{"Option 1", "Option 2"},
	}

	// Test Set and Get
	s.Set(id, radioState)
	got := s.Get(id)

	if got == nil {
		t.Fatal("Get() returned nil")
	}
	if got.GetType() != state.WidgetTypeRadio {
		t.Errorf("got.GetType() = %v, want %v", got.GetType(), state.WidgetTypeRadio)
	}
}

func TestState_ResetStates(t *testing.T) {
	s := newState()
	id := uuid.Must(uuid.NewV4())

	radioState := &state.RadioState{
		ID:      id,
		Label:   "Test Radio",
		Options: []string{"Option 1", "Option 2"},
	}

	s.Set(id, radioState)
	s.ResetStates()

	if got := s.Get(id); got != nil {
		t.Errorf("Get(%v) after reset = %v, want nil", id, got)
	}
}

func TestState_ResetButtons(t *testing.T) {
	s := newState()
	buttonID := uuid.Must(uuid.NewV4())
	formID := uuid.Must(uuid.NewV4())

	// Set initial states
	buttonState := &state.ButtonState{
		ID:    buttonID,
		Value: true,
	}
	formState := &state.FormState{
		ID:    formID,
		Value: true,
	}

	s.Set(buttonID, buttonState)
	s.Set(formID, formState)

	// Reset buttons
	s.ResetButtons()

	// Check button state
	if got := s.GetButton(buttonID); got.Value {
		t.Error("button value after reset = true, want false")
	}

	// Check form state
	if got := s.GetForm(formID); got.Value {
		t.Error("form value after reset = true, want false")
	}
}

func TestState_SetStates(t *testing.T) {
	s := newState()
	id1 := uuid.Must(uuid.NewV4())
	id2 := uuid.Must(uuid.NewV4())

	states := map[uuid.UUID]WidgetState{
		id1: &state.RadioState{
			ID:      id1,
			Label:   "Radio 1",
			Options: []string{"Option 1", "Option 2"},
		},
		id2: &state.RadioState{
			ID:      id2,
			Label:   "Radio 2",
			Options: []string{"Option 3", "Option 4"},
		},
	}

	s.SetStates(states)

	// Verify both states were set
	for id := range states {
		if got := s.Get(id); got == nil {
			t.Errorf("Get(%v) = nil, want non-nil", id)
		}
	}
}

func TestState_ConcurrentAccess(t *testing.T) {
	s := newState()
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := uuid.Must(uuid.NewV4())

			radioState := &state.RadioState{
				ID:      id,
				Label:   "Test Radio",
				Options: []string{"Option 1", "Option 2"},
			}

			s.Set(id, radioState)
			got := s.Get(id)
			if got == nil {
				t.Errorf("Get(%v) = nil, want non-nil", id)
			}
		}()
	}

	wg.Wait()
}
