import { Checkbox } from '@/components/ui/checkbox';
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import clsx from 'clsx';
import { useMemo, useState } from 'react';
import type { FC } from 'react';
import { useDeepCompareEffect } from 'use-deep-compare';

export const WidgetTable: FC<{
  widgetId: string;
}> = ({ widgetId }) => {
  const dispatch = useDispatch();
  const [page, setPage] = useState<number>(1);
  const widget = useSelector((state) =>
    widgetsStore.selector.getWidget(state, widgetId),
  );
  const state = useSelector((state) =>
    widgetsStore.selector.getWidgetState(state, widgetId),
  );
  const isWidgetWaiting = useSelector((state) => state.widgets.isWidgetWaiting);

  const handleRowClick = (row: number) => {
    if (isWidgetWaiting || state.type !== 'table') {
      return;
    }
    const value = {
      selection: {
        row: widget.widget?.table?.rowSelection === 'single' ? row : 0,
        rows:
          widget.widget?.table?.rowSelection === 'multiple'
            ? state.value?.selection?.rows?.includes(row)
              ? (state.value?.selection?.rows?.filter((r) => r !== row) ?? [])
              : [...(state.value?.selection?.rows || []), row]
            : [],
      },
    };
    dispatch(
      widgetsStore.actions.setWidgetState({
        widgetId,
        widgetType: 'table',
        value,
      }),
    );
    dispatch(
      widgetsStore.actions.setWidgetValue({
        widgetId,
        widgetType: 'table',
        value,
      }),
    );
  };

  const tableData = useMemo(() => {
    if (!widget) {
      return {
        keys: [],
        data: [],
        pageCount: 0,
      };
    }
    if (!widget.widget?.table) {
      return {
        keys: [],
        data: [],
        pageCount: 0,
      };
    }
    const data = JSON.parse(atob(widget.widget.table.data ?? '')) ?? [];
    const keys = Object.keys(data?.[0] ?? {}) as string[];
    console.log({ data });
    return {
      keys: widget.widget.table.columnOrder
        ? keys.sort((a, b) => {
            const columnOrder = widget.widget?.table?.columnOrder ?? [];
            const aIndex = columnOrder?.indexOf(a) ?? -1;
            const bIndex = columnOrder?.indexOf(b) ?? -1;
            if (aIndex === -1 && bIndex === -1) {
              return 0;
            }
            if (aIndex === -1) {
              return 1;
            }
            if (bIndex === -1) {
              return -1;
            }
            return aIndex - bIndex;
          })
        : keys,
      data: widget.widget.table.height
        ? data.slice(
            (page - 1) * widget.widget.table.height,
            page * widget.widget.table.height,
          )
        : data,
      pageCount: Math.ceil(data.length / (widget.widget.table.height || 10)),
    };
  }, [widget, page]);

  useDeepCompareEffect(() => {
    setPage(1);
  }, [widget?.widget?.table?.data]);

  return (
    widget &&
    widget.widget?.table && (
      <div className="grid grid-cols-1 gap-4">
        {(widget.widget.table.header || widget.widget.table.description) && (
          <div className="space-y-1">
            {widget.widget.table.header && (
              <p className="text-foreground text-xl font-bold">
                {widget.widget.table.header}
              </p>
            )}
            {widget.widget.table.description && (
              <p className="text-muted-foreground text-sm">
                {widget.widget.table.description}
              </p>
            )}
          </div>
        )}
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                {widget.widget?.table?.rowSelection === 'multiple' && (
                  <TableHead></TableHead>
                )}
                {tableData.keys.map((key) => (
                  <TableHead key={key} className="truncate">
                    {key}
                  </TableHead>
                ))}
              </TableRow>
            </TableHeader>
            <TableBody>
              {tableData.data.map((row: any, index: number) => (
                <TableRow
                  key={index}
                  data-state={
                    state?.type === 'table' &&
                    (widget.widget?.table?.rowSelection === 'single'
                      ? state.value?.selection?.row === index
                      : state.value?.selection?.rows?.includes(index))
                      ? 'selected'
                      : 'default'
                  }
                  onClick={() => {
                    if (widget.widget?.table?.rowSelection === 'single') {
                      handleRowClick(index);
                    }
                  }}
                  className="cursor-pointer"
                >
                  {widget.widget?.table?.rowSelection === 'multiple' && (
                    <TableCell>
                      <Checkbox
                        checked={
                          state?.type === 'table' &&
                          state.value?.selection?.rows?.includes(index)
                        }
                        onCheckedChange={() => {
                          handleRowClick(index);
                        }}
                      />
                    </TableCell>
                  )}
                  {tableData.keys.map((key) => (
                    <TableCell key={key} className="truncate">
                      {typeof row[key] === 'string'
                        ? row[key]
                        : JSON.stringify(row[key])}
                    </TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <div className="bg-muted border-t p-4">
            <Pagination className="flex justify-end">
              <PaginationContent>
                {page !== 1 && (
                  <PaginationItem>
                    <PaginationPrevious
                      className={clsx(page !== 1 && 'cursor-pointer')}
                      onClick={() => {
                        if (page === 1) {
                          return;
                        }
                        setPage(page - 1);
                      }}
                    />
                  </PaginationItem>
                )}
                {Array.from({ length: tableData.pageCount }).map((_, index) => {
                  if (index > (page || 1) + 2 || index < (page || 1) - 3) {
                    return null;
                  }
                  if (index === (page || 1) + 2 || index === (page || 1) - 3) {
                    return (
                      <PaginationItem key={index}>
                        <PaginationEllipsis />
                      </PaginationItem>
                    );
                  }
                  return (
                    <PaginationItem key={index}>
                      <PaginationLink
                        onClick={() => setPage(index + 1)}
                        className={clsx(
                          page !== index + 1 && 'cursor-pointer',
                          page === index + 1 && 'pointer-events-none',
                        )}
                        isActive={
                          page === index + 1 || (index === 0 && page === null)
                        }
                      >
                        {index + 1}
                      </PaginationLink>
                    </PaginationItem>
                  );
                })}
                {tableData.pageCount > 1 &&
                  page !== tableData.pageCount &&
                  tableData.pageCount !== 1 && (
                    <PaginationItem>
                      <PaginationNext
                        className={clsx(
                          page !== tableData.pageCount && 'cursor-pointer',
                        )}
                        onClick={() => {
                          if (page === tableData.pageCount) {
                            return;
                          }
                          setPage(page + 1);
                        }}
                      />
                    </PaginationItem>
                  )}
              </PaginationContent>
            </Pagination>
          </div>
        </div>
      </div>
    )
  );
};
