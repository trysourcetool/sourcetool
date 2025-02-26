package sourcetool

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/datetimeinput"
	"github.com/trysourcetool/sourcetool-go/internal/conv"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	websocketv1 "github.com/trysourcetool/sourcetool-proto/go/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-proto/go/widget/v1"
)

func (b *uiBuilder) DateTimeInput(label string, opts ...datetimeinput.Option) *time.Time {
	dateTimeInputOpts := &options.DateTimeInputOptions{
		Label:        label,
		Placeholder:  "",
		DefaultValue: nil,
		Required:     false,
		Disabled:     false,
		Format:       "YYYY/MM/DD HH:MM:SS",
		MaxValue:     nil,
		MinValue:     nil,
		Location:     time.Local,
	}

	for _, o := range opts {
		o.Apply(dateTimeInputOpts)
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

	widgetID := b.generateDateTimeInputID(label, path)
	dateTimeInputState := sess.State.GetDateTimeInput(widgetID)
	if dateTimeInputState == nil {
		dateTimeInputState = &state.DateTimeInputState{
			ID:    widgetID,
			Value: dateTimeInputOpts.DefaultValue,
		}
	}
	dateTimeInputState.Label = dateTimeInputOpts.Label
	dateTimeInputState.Placeholder = dateTimeInputOpts.Placeholder
	dateTimeInputState.DefaultValue = dateTimeInputOpts.DefaultValue
	dateTimeInputState.Required = dateTimeInputOpts.Required
	dateTimeInputState.Disabled = dateTimeInputOpts.Disabled
	dateTimeInputState.Format = dateTimeInputOpts.Format
	dateTimeInputState.MaxValue = dateTimeInputOpts.MaxValue
	dateTimeInputState.MinValue = dateTimeInputOpts.MinValue
	dateTimeInputState.Location = dateTimeInputOpts.Location
	sess.State.Set(widgetID, dateTimeInputState)

	dateTimeInput := convertStateToDateTimeInputProto(dateTimeInputState)
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_DateTimeInput{
				DateTimeInput: dateTimeInput,
			},
		},
	})

	cursor.next()

	return dateTimeInputState.Value
}

func (b *uiBuilder) generateDateTimeInputID(label string, path path) uuid.UUID {
	page := b.page
	if page == nil {
		return uuid.Nil
	}
	return uuid.NewV5(page.id, state.WidgetTypeDateTimeInput.String()+"-"+label+"-"+path.String())
}

func convertDateTimeInputProtoToState(id uuid.UUID, data *widgetv1.DateTimeInput, location *time.Location) (*state.DateTimeInputState, error) {
	if data == nil {
		return nil, nil
	}

	parseDate := func(dateStr string) (*time.Time, error) {
		if dateStr == "" {
			return nil, nil
		}
		t, err := time.ParseInLocation(time.DateTime, dateStr, location)
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

	return &state.DateTimeInputState{
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

func convertStateToDateTimeInputProto(state *state.DateTimeInputState) *widgetv1.DateTimeInput {
	if state == nil {
		return nil
	}
	var value, defaultValue, maxValue, minValue string
	if state.Value != nil {
		value = state.Value.Format(time.DateTime)
	}
	if state.DefaultValue != nil {
		defaultValue = state.DefaultValue.Format(time.DateTime)
	}
	if state.MaxValue != nil {
		maxValue = state.MaxValue.Format(time.DateTime)
	}
	if state.MinValue != nil {
		minValue = state.MinValue.Format(time.DateTime)
	}
	return &widgetv1.DateTimeInput{
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
