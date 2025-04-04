package sourcetool

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gofrs/uuid/v5"

	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
	"github.com/trysourcetool/sourcetool-go/table"
)

type testData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestConvertStateToTableProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := []testData{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}
	selection := &state.TableStateValueSelection{
		Row:  0,
		Rows: []int32{0},
	}

	tableState := &state.TableState{
		ID:           id,
		Data:         data,
		Header:       "Test Table",
		Description:  "Test Description",
		OnSelect:     table.SelectionBehaviorRerun.String(),
		RowSelection: table.SelectionModeSingle.String(),
		Value: state.TableStateValue{
			Selection: selection,
		},
	}

	tableData, err := convertStateToTableProto(tableState)
	if err != nil {
		t.Fatalf("convertStateToTableProto returned error: %v", err)
	}

	if tableData == nil {
		t.Fatal("convertStateToTableProto returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Header", tableData.Header, tableState.Header},
		{"Description", tableData.Description, tableState.Description},
		{"OnSelect", tableData.OnSelect, tableState.OnSelect},
		{"RowSelection", tableData.RowSelection, tableState.RowSelection},
		{"Selection.Row", tableData.Value.Selection.Row, tableState.Value.Selection.Row},
		{"Selection.Rows", tableData.Value.Selection.Rows, tableState.Value.Selection.Rows},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}
	if !reflect.DeepEqual(tableData.Data, dataBytes) {
		t.Errorf("Data = %v, want %v", tableData.Data, data)
	}
}

func TestConvertTableProtoToState(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	data := []testData{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}
	selection := &widgetv1.TableValueSelection{
		Row:  0,
		Rows: []int32{0},
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	tableData := &widgetv1.Table{
		Data:         dataBytes,
		Header:       "Test Table",
		Description:  "Test Description",
		OnSelect:     table.SelectionBehaviorRerun.String(),
		RowSelection: table.SelectionModeSingle.String(),
		Value: &widgetv1.TableValue{
			Selection: selection,
		},
	}

	state := convertTableProtoToState(id, tableData)

	if state == nil {
		t.Fatal("convertTableProtoToState returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ID", state.ID, id},
		{"Header", state.Header, tableData.Header},
		{"Description", state.Description, tableData.Description},
		{"OnSelect", state.OnSelect, tableData.OnSelect},
		{"RowSelection", state.RowSelection, tableData.RowSelection},
		{"Selection.Row", state.Value.Selection.Row, tableData.Value.Selection.Row},
		{"Selection.Rows", state.Value.Selection.Rows, tableData.Value.Selection.Rows},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	dataBytes, err = json.Marshal(data)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}
	if !reflect.DeepEqual(state.Data, dataBytes) {
		t.Errorf("Data = %v, want %v", state.Data, data)
	}
}

func TestTable(t *testing.T) {
	sessionID := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())
	sess := session.New(sessionID, pageID)

	mockWS := mock.NewClient()

	builder := &uiBuilder{
		context: context.Background(),
		session: sess,
		cursor:  newCursor(),
		page: &page{
			id: pageID,
		},
		runtime: &runtime{
			wsClient: mockWS,
		},
	}

	data := []testData{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}
	header := "Test Table"
	description := "Test Description"

	builder.Table(data,
		table.Header(header),
		table.Description(description),
		table.OnSelect(table.SelectionBehaviorRerun),
		table.RowSelection(table.SelectionModeSingle),
	)

	messages := mockWS.Messages()
	if len(messages) != 1 {
		t.Errorf("WebSocket messages count = %d, want 1", len(messages))
	}
	msg := messages[0]
	if v := msg.GetRenderWidget(); v == nil {
		t.Fatal("WebSocket message type = nil, want RenderWidget")
	}

	widgetID := builder.generatePageID(state.WidgetTypeTable, []int{0})
	state := sess.State.GetTable(widgetID)
	if state == nil {
		t.Fatal("Table state not found")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Header", state.Header, header},
		{"Description", state.Description, description},
		{"OnSelect", state.OnSelect, table.SelectionBehaviorRerun.String()},
		{"RowSelection", state.RowSelection, table.SelectionModeSingle.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	if !reflect.DeepEqual(state.Data, data) {
		t.Errorf("Data = %v, want %v", state.Data, data)
	}
}

func TestTable_DefaultValues(t *testing.T) {
	sessionID := uuid.Must(uuid.NewV4())
	pageID := uuid.Must(uuid.NewV4())
	sess := session.New(sessionID, pageID)

	mockWS := mock.NewClient()

	builder := &uiBuilder{
		context: context.Background(),
		session: sess,
		cursor:  newCursor(),
		page: &page{
			id: pageID,
		},
		runtime: &runtime{
			wsClient: mockWS,
		},
	}

	data := []testData{
		{ID: 1, Name: "Test 1"},
		{ID: 2, Name: "Test 2"},
	}

	builder.Table(data)

	widgetID := builder.generatePageID(state.WidgetTypeTable, []int{0})
	state := sess.State.GetTable(widgetID)
	if state == nil {
		t.Fatal("Table state not found")
	}

	if state.OnSelect != table.SelectionBehaviorIgnore.String() {
		t.Errorf("Default OnSelect = %v, want %v", state.OnSelect, table.SelectionBehaviorIgnore)
	}
	if state.RowSelection != table.SelectionModeSingle.String() {
		t.Errorf("Default RowSelection = %v, want %v", state.RowSelection, table.SelectionModeSingle)
	}
	if state.Header != "" {
		t.Errorf("Default Header = %v, want empty string", state.Header)
	}
	if state.Description != "" {
		t.Errorf("Default Description = %v, want empty string", state.Description)
	}
}
