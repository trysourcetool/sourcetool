package sourcetool

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"

	"github.com/trysourcetool/sourcetool-go/columns"
	"github.com/trysourcetool/sourcetool-go/internal/session"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/internal/websocket/mock"
)

func TestConvertStateToColumnsProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	columnsState := &state.ColumnsState{
		ID:      id,
		Columns: 3,
	}

	data := convertStateToColumnsProto(columnsState)

	if data == nil {
		t.Fatal("convertStateToColumnsProto returned nil")
	}

	if int(data.Columns) != columnsState.Columns {
		t.Errorf("Columns = %v, want %v", data.Columns, columnsState.Columns)
	}
}

func TestConvertColumnsProtoToState(t *testing.T) {
	data := &widgetv1.Columns{
		Columns: 3,
	}

	state := convertColumnsProtoToState(uuid.Must(uuid.NewV4()), data)

	if state == nil {
		t.Fatal("convertColumnsProtoToState returned nil")
	}

	if state.Columns != int(data.Columns) {
		t.Errorf("Columns = %v, want %v", state.Columns, data.Columns)
	}
}

func TestConvertStateToColumnItemProto(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	columnItemState := &state.ColumnItemState{
		ID:     id,
		Weight: 0.5,
	}

	data := convertStateToColumnItemProto(columnItemState)

	if data == nil {
		t.Fatal("convertStateToColumnItemProto returned nil")
	}

	if data.Weight != columnItemState.Weight {
		t.Errorf("Weight = %v, want %v", data.Weight, columnItemState.Weight)
	}
}

func TestConvertColumnItemProtoToState(t *testing.T) {
	data := &widgetv1.ColumnItem{
		Weight: 0.5,
	}

	state := convertColumnItemProtoToState(uuid.Must(uuid.NewV4()), data)

	if state == nil {
		t.Fatal("convertColumnItemProtoToState returned nil")
	}

	if state.Weight != data.Weight {
		t.Errorf("Weight = %v, want %v", state.Weight, data.Weight)
	}
}

func TestColumns(t *testing.T) {
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

	cols := 3
	builders := builder.Columns(cols)

	if builders == nil {
		t.Fatal("Columns returned nil")
	}
	if len(builders) != cols {
		t.Errorf("Builders length = %v, want %v", len(builders), cols)
	}

	messages := mockWS.Messages()
	expectedMessages := cols + 1 // columns widget + column items
	if len(messages) != expectedMessages {
		t.Errorf("WebSocket messages count = %d, want %d", len(messages), expectedMessages)
	}

	widgetID := builder.generatePageID(state.WidgetTypeColumns, []int{0})
	columnsState := sess.State.GetColumns(widgetID)
	if columnsState == nil {
		t.Fatal("Columns state not found")
	}

	if columnsState.Columns != cols {
		t.Errorf("Columns = %v, want %v", columnsState.Columns, cols)
	}

	for i := 0; i < cols; i++ {
		columnPath := []int{0, i}
		columnID := builder.generatePageID(state.WidgetTypeColumnItem, columnPath)
		columnState := sess.State.Get(columnID)
		if columnState == nil {
			t.Fatalf("Column item state not found for index %d", i)
		}

		expectedWeight := 1.0 / float64(cols)
		columnItemState, ok := columnState.(*state.ColumnItemState)
		if !ok {
			t.Fatalf("Column item state[%d] is not *columnitem.State", i)
		}
		if columnItemState.Weight != expectedWeight {
			t.Errorf("Column item weight[%d] = %v, want %v", i, columnItemState.Weight, expectedWeight)
		}
	}
}

func TestColumns_WithWeight(t *testing.T) {
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

	cols := 3
	weights := []int{2, 1, 1}
	totalWeight := 4

	builders := builder.Columns(cols, columns.Weight(weights...))

	if builders == nil {
		t.Fatal("Columns returned nil")
	}

	for i := 0; i < cols; i++ {
		columnPath := []int{0, i}
		columnID := builder.generatePageID(state.WidgetTypeColumnItem, columnPath)
		columnState := sess.State.Get(columnID)
		if columnState == nil {
			t.Fatalf("Column item state not found for index %d", i)
		}

		expectedWeight := float64(weights[i]) / float64(totalWeight)
		columnItemState, ok := columnState.(*state.ColumnItemState)
		if !ok {
			t.Fatalf("Column item state[%d] is not *columnitem.State", i)
		}
		if columnItemState.Weight != expectedWeight {
			t.Errorf("Column item weight[%d] = %v, want %v", i, columnItemState.Weight, expectedWeight)
		}
	}
}

func TestColumns_InvalidInput(t *testing.T) {
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

	tests := []struct {
		name    string
		cols    int
		weights []int
	}{
		{"Zero columns", 0, nil},
		{"Negative columns", -1, nil},
		{"Invalid weights length", 3, []int{1, 1}},
		{"Zero weights", 3, []int{0, 0, 0}},
		{"Negative weights", 3, []int{-1, 1, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builders []UIBuilder
			if tt.weights == nil {
				builders = builder.Columns(tt.cols)
			} else {
				builders = builder.Columns(tt.cols, columns.Weight(tt.weights...))
			}

			if tt.cols <= 0 && builders != nil {
				t.Error("Expected nil builders for invalid column count")
			}

			if tt.cols > 0 && builders != nil {
				for i := 0; i < tt.cols; i++ {
					columnPath := []int{0, i}
					columnID := builder.generatePageID(state.WidgetTypeColumnItem, columnPath)
					columnState := sess.State.Get(columnID)
					if columnState == nil {
						t.Fatalf("Column item state not found for index %d", i)
					}

					expectedWeight := 1.0 / float64(tt.cols)
					columnItemState, ok := columnState.(*state.ColumnItemState)
					if !ok {
						t.Fatalf("Column item state[%d] is not *columnitem.State", i)
					}
					if columnItemState.Weight != expectedWeight {
						t.Errorf("Column item weight[%d] = %v, want %v", i, columnItemState.Weight, expectedWeight)
					}
				}
			}
		})
	}
}
