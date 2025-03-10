import { PageHeader } from '@/components/common/page-header';
import { Button } from '@/components/ui/button';
import dayjs from 'dayjs';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import Fuse from 'fuse.js';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useDispatch, useSelector } from '@/store';
import { environmentsStore } from '@/store/modules/environments';
import { Ellipsis, Loader2, Plus } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { $path } from 'safe-routes';
import { Link, useNavigate } from 'react-router';
import clsx from 'clsx';
import { useTranslation } from 'react-i18next';
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
import { useQueryState } from 'nuqs';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useToast } from '@/hooks/use-toast';
import type { ErrorResponse } from '@/api/instance';
import { Input } from '@/components/ui/input';

export default function Environments() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const { toast } = useToast();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const [search, setSearch] = useState('');
  const environments = useSelector(environmentsStore.selector.getEnvironments);
  const isDeleteEnvironmentWaiting = useSelector(
    (state) => state.environments.isDeleteEnvironmentWaiting,
  );
  const [selectEnvironmentId, setSelectEnvironmentId] = useState<string | null>(
    null,
  );
  const navigate = useNavigate();
  const [page, setPage] = useQueryState('page', {
    parse: (query: string) => parseInt(query, 10),
    serialize: (value) => value.toString(),
  });

  const pageSize = 10;

  const pageCount = Math.ceil(environments.length / pageSize);

  const filteredEnvironments = useMemo(() => {
    if (!search) {
      return environments;
    }
    const fuse = new Fuse(environments, {
      keys: ['name'],
    });

    return fuse.search(search).map((result) => result.item);
  }, [environments, search]);

  const slicedEnvironments = useMemo(() => {
    return filteredEnvironments.slice(
      (page || 1) * pageSize - pageSize,
      (page || 1) * pageSize,
    );
  }, [environments, page]);

  const handleDeleteEnvironment = async () => {
    if (isDeleteEnvironmentWaiting || !selectEnvironmentId) {
      return;
    }

    const resultAction = await dispatch(
      environmentsStore.asyncActions.deleteEnvironment({
        environmentId: selectEnvironmentId,
      }),
    );
    if (
      environmentsStore.asyncActions.deleteEnvironment.fulfilled.match(
        resultAction,
      )
    ) {
      setSelectEnvironmentId(null);
      toast({
        title: t('routes_environments_toast_deleted'),
        description: t('routes_environments_toast_deleted_description'),
      });
    } else {
      console.error(resultAction);
      toast({
        title: t('routes_environments_toast_delete_failed'),
        description:
          (resultAction.payload as ErrorResponse)?.detail ||
          t('routes_environments_toast_delete_failed_description'),
        variant: 'destructive',
      });
    }
  };

  useEffect(() => {
    setBreadcrumbsState?.([{ label: t('breadcrumbs_environment') }]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        await dispatch(environmentsStore.asyncActions.listEnvironments());
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch]);

  return (
    <div>
      <PageHeader label={t('routes_environments_page_header')} />
      <div className="flex w-screen flex-col gap-6 p-6 md:w-auto">
        <div className="flex flex-col justify-between gap-2 pt-6 md:flex-row">
          <p className="text-xl font-bold text-foreground">
            {t('routes_environments_title')}
          </p>
        </div>
        <div className="flex justify-between gap-2">
          <div className="hidden max-w-72 flex-1 md:block">
            <Input
              placeholder={t('routes_environments_search_placeholder')}
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
              }}
            />
          </div>
          <Button asChild>
            <Link to={$path('/environments/new')}>
              <Plus />
              {t('routes_environments_create_new')}
            </Link>
          </Button>
        </div>
        <div className="w-full overflow-auto rounded-md border">
          <Table className="md:table-fixed">
            <TableHeader>
              <TableRow>
                <TableHead>{t('routes_environments_table_name')}</TableHead>
                <TableHead>{t('routes_environments_table_color')}</TableHead>
                <TableHead>
                  {t('routes_environments_table_created_at')}
                </TableHead>
                <TableHead className="w-16"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredEnvironments.map((environment) => (
                <TableRow
                  key={environment.id}
                  className="cursor-pointer"
                  onClick={() => {
                    navigate(
                      $path('/environments/:environmentId', {
                        environmentId: environment.id,
                      }),
                    );
                  }}
                >
                  <TableCell className="truncate">{environment.name}</TableCell>
                  <TableCell>
                    <div
                      className="size-5 rounded-full"
                      style={{ backgroundColor: environment.color }}
                    />
                  </TableCell>
                  <TableCell className="truncate">
                    {dayjs
                      .unix(Number(environment.createdAt))
                      .format('DD MMM YYYY')}
                  </TableCell>
                  <TableCell align="right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                          <Ellipsis className="size-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem
                          onClick={(e) => {
                            e.stopPropagation();
                            if (
                              environment.slug === 'production' ||
                              environment.slug === 'development'
                            ) {
                              return;
                            }
                            setSelectEnvironmentId(environment.id);
                          }}
                          className={clsx(
                            environment.slug === 'production' ||
                              environment.slug === 'development'
                              ? 'text-popover-foreground opacity-50'
                              : 'text-destructive',
                          )}
                        >
                          {t('routes_environments_delete')}
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <div className="border-t bg-muted p-4">
            <Pagination className="flex justify-end">
              <PaginationContent>
                {page !== 1 && page !== null && (
                  <PaginationItem>
                    <PaginationPrevious
                      className={clsx(
                        (page !== 1 || page === null) && 'cursor-pointer',
                      )}
                      onClick={() => {
                        if (page === 1 || page === null) {
                          return;
                        }
                        setPage((page || 1) - 1);
                      }}
                    />
                  </PaginationItem>
                )}
                {Array.from({ length: pageCount }).map((_, index) => {
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

                {page !== pageCount && pageCount !== 1 && (
                  <PaginationItem>
                    <PaginationNext
                      onClick={() => {
                        if (page === pageCount) {
                          return;
                        }
                        setPage((page || 1) + 1);
                      }}
                      className={clsx(page !== pageCount && 'cursor-pointer')}
                    />
                  </PaginationItem>
                )}
              </PaginationContent>
            </Pagination>
          </div>
        </div>
      </div>
      <Dialog
        open={selectEnvironmentId !== null}
        onOpenChange={(boolean) => {
          if (isDeleteEnvironmentWaiting) {
            return;
          }
          if (!boolean) {
            setSelectEnvironmentId(null);
          }
        }}
      >
        <DialogContent
          onCloseAutoFocus={(event) => {
            event.preventDefault();
            document.body.style.pointerEvents = '';
          }}
          className="max-h-[calc(100svh-48px)] max-w-lg overflow-y-auto"
          data-state={selectEnvironmentId !== null ? 'open' : 'closed'}
        >
          <DialogHeader>
            <DialogTitle>
              {t('routes_environments_delete_dialog_title')}
            </DialogTitle>
            <DialogDescription>
              {t('routes_environments_delete_dialog_description')}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="flex justify-end gap-2">
            <Button
              variant="outline"
              onClick={() => {
                if (isDeleteEnvironmentWaiting) {
                  return;
                }
                setSelectEnvironmentId(null);
              }}
            >
              {t('routes_environments_delete_dialog_cancel')}
            </Button>
            <Button
              onClick={handleDeleteEnvironment}
              disabled={isDeleteEnvironmentWaiting}
            >
              {isDeleteEnvironmentWaiting && (
                <Loader2 className="size-4 animate-spin" />
              )}
              {t('routes_environments_delete_dialog_delete')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
