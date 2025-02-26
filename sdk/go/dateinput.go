package sourcetool

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/dateinput"
	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
)

func (b *uiBuilder) DateInput(label string, opts ...dateinput.Option) *time.Time {
	dateInputOpts := &options.DateInputOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		Format:       "YYYY/MM/DD",
		MaxValue:     nil,
		MinValue:     nil,
		Location:     time.Local,
	}

	for _, o := range opts {
		o.Apply(dateInputOpts)
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

	widgetID := b.generateDateInputID(label, path)
	dateInputState := sess.State.GetDateInput(widgetID)
	if dateInputState == nil {
		dateInputState = &state.DateInputState{
			ID:    widgetID,
			Value: dateInputOpts.DefaultValue,
		}
	}
	dateInputState.Label = dateInputOpts.Label
	dateInputState.Placeholder = dateInputOpts.Placeholder
	dateInputState.DefaultValue = dateInputOpts.DefaultValue
	dateInputState.Required = dateInputOpts.Required
	dateInputState.Disabled = dateInputOpts.Disabled
	dateInputState.Format = dateInputOpts.Format
	dateInputState.MaxValue = dateInputOpts.MaxValue
	dateInputState.MinValue = dateInputOpts.MinValue
	dateInputState.Location = dateInputOpts.Location
	sess.State.Set(widgetID, dateInputState)

	dateInput := convertStateToDateInputProto(dateInputState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_DateInput{
				DateInput: dateInput,
			},
		},
	})

	cursor.next()

	return dateInputState.Value
}

func (b *uiBuilder) generateDateInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeDateInput.String()+"-"+label+"-"+path.String())
}

func convertDateInputProtoToState(id uuid.UUID, data *widgetv1.DateInput, location *time.Location) (*state.DateInputState, error) {
	if data == nil {
		return nil, nil
	}

	parseDate := func(dateStr string) (*time.Time, error) {
		if dateStr == "" {
			return nil, nil
		}
		t, err := time.ParseInLocation(time.DateOnly, dateStr, location)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date %q: %v", dateStr, err)
		}
		return &t, nil
	}

	value, err := parseDate(conv.SafeValue(data.Value))
	if err != nil {
		return nil, err
	}

	defaultValue, err := parseDate(conv.SafeValue(data.DefaultValue))
	if err != nil {
		return nil, err
	}

	maxValue, err := parseDate(data.MaxValue)
	if err != nil {
		return nil, err
	}

	minValue, err := parseDate(data.MinValue)
	if err != nil {
		return nil, err
	}

	return &state.DateInputState{
		ID:           id,
		Value:        value,
		Label:        data.Label,
		DefaultValue: defaultValue,
		Placeholder:  data.Placeholder,
		Required:     data.Required,
		Disabled:     data.Disabled,
		Format:       data.Format,
		MaxValue:     maxValue,
		MinValue:     minValue,
		Location:     location,
	}, nil
}

func convertStateToDateInputProto(state *state.DateInputState) *widgetv1.DateInput {
	if state == nil {
		return nil
	}
	var value, defaultValue, maxValue, minValue string
	if state.Value != nil {
		value = state.Value.Format(time.DateOnly)
	}
	if state.DefaultValue != nil {
		defaultValue = state.DefaultValue.Format(time.DateOnly)
	}
	if state.MaxValue != nil {
		maxValue = state.MaxValue.Format(time.DateOnly)
	}
	if state.MinValue != nil {
		minValue = state.MinValue.Format(time.DateOnly)
	}
	return &widgetv1.DateInput{
		Value:        conv.NilValue(value),
		Label:        state.Label,
		Placeholder:  state.Placeholder,
		DefaultValue: conv.NilValue(defaultValue),
		Required:     state.Required,
		Disabled:     state.Disabled,
		Format:       state.Format,
		MaxValue:     maxValue,
		MinValue:     minValue,
	}
}
