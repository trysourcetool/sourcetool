package sourcetool

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/options"
	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/pb/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/pb/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/ptrconv"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/timeinput"
)

func (b *uiBuilder) TimeInput(label string, opts ...timeinput.Option) *time.Time {
	timeInputOpts := &options.TimeInputOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		Location:     time.Local,
	}

	for _, o := range opts {
		o.Apply(timeInputOpts)
	}

	sess := b.session
	if sess == nil {
		return nil
	}
	page := b.page
	if page == nil {
		return nil
	}
	cursor := b.cursor
	if cursor == nil {
		return nil
	}
	path := cursor.getPath()

	widgetID := b.generatePageID(state.WidgetTypeTimeInput, path)
	timeInputState := sess.State.GetTimeInput(widgetID)
	if timeInputState == nil {
		timeInputState = &state.TimeInputState{
			ID:    widgetID,
			Value: timeInputOpts.DefaultValue,
		}
	}
	timeInputState.Label = timeInputOpts.Label
	timeInputState.Placeholder = timeInputOpts.Placeholder
	timeInputState.DefaultValue = timeInputOpts.DefaultValue
	timeInputState.Required = timeInputOpts.Required
	timeInputState.Disabled = timeInputOpts.Disabled
	timeInputState.Location = timeInputOpts.Location
	sess.State.Set(widgetID, timeInputState)

	timeInput := convertStateToTimeInputProto(timeInputState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_TimeInput{
				TimeInput: timeInput,
			},
		},
	})

	cursor.next()

	return timeInputState.Value
}

func convertTimeInputProtoToState(id uuid.UUID, data *widgetv1.TimeInput, location *time.Location) (*state.TimeInputState, error) {
	if data == nil {
		return nil, nil
	}

	parseTime := func(timeStr string) (*time.Time, error) {
		if timeStr == "" {
			return nil, nil
		}
		t, err := time.ParseInLocation(time.TimeOnly, timeStr, location)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time %q: %v", timeStr, err)
		}
		return &t, nil
	}

	value, err := parseTime(ptrconv.StringValue(data.Value))
	if err != nil {
		return nil, err
	}

	defaultValue, err := parseTime(ptrconv.StringValue(data.DefaultValue))
	if err != nil {
		return nil, err
	}

	return &state.TimeInputState{
		ID:           id,
		Value:        value,
		Label:        data.Label,
		DefaultValue: defaultValue,
		Placeholder:  data.Placeholder,
		Required:     data.Required,
		Disabled:     data.Disabled,
		Location:     location,
	}, nil
}

func convertStateToTimeInputProto(state *state.TimeInputState) *widgetv1.TimeInput {
	if state == nil {
		return nil
	}
	var value, defaultValue string
	if state.Value != nil {
		value = state.Value.Format(time.TimeOnly)
	}
	if state.DefaultValue != nil {
		defaultValue = state.DefaultValue.Format(time.TimeOnly)
	}
	return &widgetv1.TimeInput{
		Value:        ptrconv.StringPtr(value),
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: ptrconv.StringPtr(defaultValue),
		Required:     state.Required,
		Disabled:     state.Disabled,
	}
}
