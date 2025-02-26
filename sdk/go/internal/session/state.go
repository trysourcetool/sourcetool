package session

import (
	"sync"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/session/state"
)

type WidgetState interface {
	IsWidgetState()
	GetType() state.WidgetType
}

type StateData map[uuid.UUID]WidgetState

type State struct {
	// data map[uuid.UUID]any // ui ID -> options state
	data StateData
	mu   sync.RWMutex
}

func newState() *State {
	return &State{
		data: make(map[uuid.UUID]WidgetState),
	}
}

func (s *State) Get(id uuid.UUID) WidgetState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.data[id]
	if !ok {
		return nil
	}

	return state
}

func (s *State) GetTextInput(id uuid.UUID) *state.TextInputState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.TextInputState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetNumberInput(id uuid.UUID) *state.NumberInputState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.NumberInputState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetDateInput(id uuid.UUID) *state.DateInputState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.DateInputState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetDateTimeInput(id uuid.UUID) *state.DateTimeInputState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.DateTimeInputState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetTimeInput(id uuid.UUID) *state.TimeInputState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.TimeInputState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetSelectbox(id uuid.UUID) *state.SelectboxState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.SelectboxState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetMultiSelect(id uuid.UUID) *state.MultiSelectState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.MultiSelectState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetCheckbox(id uuid.UUID) *state.CheckboxState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.CheckboxState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetCheckboxGroup(id uuid.UUID) *state.CheckboxGroupState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.CheckboxGroupState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetRadio(id uuid.UUID) *state.RadioState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.RadioState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetTextArea(id uuid.UUID) *state.TextAreaState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.TextAreaState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetTable(id uuid.UUID) *state.TableState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.TableState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetButton(id uuid.UUID) *state.ButtonState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.ButtonState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetColumns(id uuid.UUID) *state.ColumnsState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.ColumnsState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetMarkdown(id uuid.UUID) *state.MarkdownState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.MarkdownState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) GetForm(id uuid.UUID) *state.FormState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.data[id]
	if !ok {
		return nil
	}

	v, ok := st.(*state.FormState)
	if !ok {
		return nil
	}

	return v
}

func (s *State) ResetStates() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[uuid.UUID]WidgetState)
}

func (s *State) ResetButtons() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, st := range s.data {
		switch st.GetType() {
		case state.WidgetTypeButton:
			buttonState, ok := st.(*state.ButtonState)
			if ok {
				buttonState.Value = false
				s.data[id] = buttonState
			}
		case state.WidgetTypeForm:
			formState, ok := st.(*state.FormState)
			if ok {
				formState.Value = false
				s.data[id] = formState
			}
		}
	}
}

func (s *State) Set(id uuid.UUID, state WidgetState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = state
}

func (s *State) SetStates(states map[uuid.UUID]WidgetState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, state := range states {
		s.data[id] = state
	}
}
