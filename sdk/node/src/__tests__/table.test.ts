import { expect, test } from 'vitest';
import { v4 as uuidv4 } from 'uuid';
import {
  TableState,
  TableOnSelect,
  TableRowSelection,
  TableStateValue,
} from '../session/state/table';
import {
  convertTableProtoToState,
  convertStateToTableProto,
  table,
} from '../uibuilder/widgets/table';
import { createSessionManager, newSession } from '../session';
import { Cursor, generateWidgetId } from '../uibuilder';
import { Page, PageManager } from '../page';
import { Runtime } from '../runtime';
import { Table as TableProto } from '../pb/widget/v1/widget_pb';
import { MockClient } from '../websocket/mock/websocket';

test('convertStateToTableProto', () => {
  const id = uuidv4();
  const data = [
    { ID: 1, Name: 'Test 1' },
    { ID: 2, Name: 'Test 2' },
  ];
  const header = 'Test Table';
  const description = 'Test Description';
  const height = 30;
  const columnOrder = ['ID', 'Name'];
  const onSelect = TableOnSelect.Rerun;
  const rowSelection = TableRowSelection.Multiple;
  const value: TableStateValue = { selection: { row: 0, rows: [0, 1] } };

  const state = new TableState(
    id,
    data,
    header,
    description,
    height,
    columnOrder,
    onSelect.toString(),
    rowSelection.toString(),
    value,
  );
  const proto = convertStateToTableProto(state);

  expect(JSON.parse(new TextDecoder().decode(proto.data))).toEqual(data);
  expect(proto.header).toBe(header);
  expect(proto.description).toBe(description);
  expect(proto.height).toBe(height);
  expect(proto.columnOrder).toEqual(columnOrder);
  expect(proto.onSelect).toBe(onSelect.toString());
  expect(proto.rowSelection).toBe(rowSelection.toString());
  expect(proto.value?.selection?.row).toBe(value.selection?.row);
  expect(proto.value?.selection?.rows).toEqual(value.selection?.rows);
});

test('convertTableProtoToState', () => {
  const id = uuidv4();
  const data = [
    { ID: 1, Name: 'Test 1' },
    { ID: 2, Name: 'Test 2' },
  ];
  const header = 'Test Table';
  const description = 'Test Description';
  const height = 30;
  const columnOrder = ['ID', 'Name'];
  const onSelect = TableOnSelect.Ignore.toString();
  const rowSelection = TableRowSelection.Single.toString();
  const value: TableStateValue = { selection: { row: 0, rows: [0, 1] } };

  const tempState = new TableState(
    id,
    data,
    header,
    description,
    height,
    columnOrder,
    onSelect,
    rowSelection,
    value,
  );
  const proto: TableProto = convertStateToTableProto(tempState);

  const state = convertTableProtoToState(id, proto);

  if (!state) {
    throw new Error('TableState not found');
  }

  expect(state.id).toBe(id);
  expect(state.data).toEqual(data);
  expect(state.header).toBe(header);
  expect(state.description).toBe(description);
  expect(state.height).toBe(height);
  expect(state.columnOrder).toEqual(columnOrder);
  expect(state.onSelect).toBe(onSelect);
  expect(state.rowSelection).toBe(rowSelection);
  expect(state.value.selection?.row).toEqual(value.selection?.row);
  expect(state.value.selection?.rows).toEqual(value.selection?.rows);
});

test('table interaction', () => {
  const sessionId = uuidv4();
  const pageId = uuidv4();
  const session = newSession(sessionId, pageId);

  const pageManager = new PageManager({
    [pageId]: new Page(
      pageId,
      'Test Page',
      '/test',
      [1, 2, 3],
      async () => {},
      ['test'],
    ),
  });

  const sessionManager = createSessionManager();
  const mockWS = new MockClient();
  const runtime = new Runtime(mockWS, sessionManager, pageManager);

  const page = pageManager.getPage(pageId);
  if (!page) {
    throw new Error('Page not found');
  }

  const tableData = [
    { ID: 1, Name: 'Test 1' },
    { ID: 2, Name: 'Test 2' },
  ];
  const options = {
    header: 'Test Table',
    description: 'Test Description',
    height: 30,
    columnOrder: ['ID', 'Name'],
    onSelect: TableOnSelect.Rerun,
    rowSelection: TableRowSelection.Single,
  };

  const cursor = new Cursor();

  table({ runtime, session, page, cursor }, tableData, options);

  const widgetId = generateWidgetId(page.id, 'table', [0]);
  const state = session.state.getTable(widgetId);

  if (!state) {
    throw new Error('TableState not found');
  }

  expect(state.id).toBe(widgetId);
  expect(state.data).toEqual(tableData);
  expect(state.header).toBe(options.header);
  expect(state.description).toBe(options.description);
  expect(state.height).toBe(options.height);
  expect(state.columnOrder).toEqual(options.columnOrder);
  expect(state.onSelect).toBe(options.onSelect.toString());
  expect(state.rowSelection).toBe(options.rowSelection.toString());
  expect(state.value).toEqual({});
});
