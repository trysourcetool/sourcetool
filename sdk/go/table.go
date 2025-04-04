package sourcetool

import (
	"encoding/json"

	"github.com/gofrs/uuid/v5"

	websocketv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/websocket/v1"
	widgetv1 "github.com/trysourcetool/sourcetool-go/internal/generated/proto/widget/v1"
	"github.com/trysourcetool/sourcetool-go/internal/options"
	"github.com/trysourcetool/sourcetool-go/internal/session/state"
	"github.com/trysourcetool/sourcetool-go/table"
)

func (b *uiBuilder) Table(data any, opts ...table.Option) table.Value {
	tableOpts := &options.TableOptions{
		OnSelect:     table.SelectionBehaviorIgnore.String(),
		RowSelection: table.SelectionModeSingle.String(),
	}

	for _, o := range opts {
		o.Apply(tableOpts)
	}

	sess := b.session
	if sess == nil {
		return table.Value{}
	}
	page := b.page
	if page == nil {
		return table.Value{}
	}
	cursor := b.cursor
	if cursor == nil {
		return table.Value{}
	}
	path := cursor.getPath()

	widgetID := b.generatePageID(state.WidgetTypeTable, path)
	tableState := sess.State.GetTable(widgetID)
	if tableState == nil {
		tableState = &state.TableState{
			ID:    widgetID,
			Value: state.TableStateValue{},
		}
	}
	tableState.Data = data
	tableState.Header = tableOpts.Header
	tableState.Description = tableOpts.Description
	tableState.Height = tableOpts.Height
	tableState.ColumnOrder = tableOpts.ColumnOrder
	tableState.OnSelect = tableOpts.OnSelect
	tableState.RowSelection = tableOpts.RowSelection
	sess.State.Set(widgetID, tableState)

	tableProto, err := convertStateToTableProto(tableState)
	if err != nil {
		return table.Value{}
	}
	b.runtime.wsClient.Enqueue(uuid.Must(uuid.NewV4()).String(), &websocketv1.RenderWidget{
		SessionId: sess.ID.String(),
		PageId:    page.id.String(),
		Path:      convertPathToInt32Slice(path),
		Widget: &widgetv1.Widget{
			Id: widgetID.String(),
			Type: &widgetv1.Widget_Table{
				Table: tableProto,
			},
		},
	})

	cursor.next()

	value := table.Value{}
	if tableState.Value.Selection != nil {
		rows := make([]int, len(tableState.Value.Selection.Rows))
		for i, r := range tableState.Value.Selection.Rows {
			rows[i] = int(r)
		}
		value.Selection = &table.Selection{
			Row:  int(tableState.Value.Selection.Row),
			Rows: rows,
		}
	}

	return value
}

func convertStateToTableProto(state *state.TableState) (*widgetv1.Table, error) {
	if state == nil {
		return nil, nil
	}
	dataBytes, err := json.Marshal(state.Data)
	if err != nil {
		return nil, err
	}
	data := &widgetv1.Table{
		Data:         dataBytes,
		Header:       state.Header,
		Description:  state.Description,
		Height:       state.Height,
		ColumnOrder:  state.ColumnOrder,
		OnSelect:     state.OnSelect,
		RowSelection: state.RowSelection,
		Value:        &widgetv1.TableValue{},
	}
	if state.Value.Selection != nil {
		data.Value.Selection = &widgetv1.TableValueSelection{
			Row:  state.Value.Selection.Row,
			Rows: state.Value.Selection.Rows,
		}
	}
	return data, nil
}

func convertTableProtoToState(id uuid.UUID, data *widgetv1.Table) *state.TableState {
	if data == nil {
		return nil
	}
	tableState := &state.TableState{
		ID:           id,
		Data:         data.Data,
		Header:       data.Header,
		Description:  data.Description,
		Height:       data.Height,
		ColumnOrder:  data.ColumnOrder,
		OnSelect:     data.OnSelect,
		RowSelection: data.RowSelection,
		Value:        state.TableStateValue{},
	}
	if data.Value.Selection != nil {
		tableState.Value.Selection = &state.TableStateValueSelection{
			Row:  data.Value.Selection.Row,
			Rows: data.Value.Selection.Rows,
		}
	}
	return tableState
}
