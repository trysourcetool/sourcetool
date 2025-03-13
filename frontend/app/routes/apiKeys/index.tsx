import { PageHeader } from '@/components/common/page-header';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
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
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useToast } from '@/hooks/use-toast';
import { useDispatch, useSelector } from '@/store';
import { apiKeysStore } from '@/store/modules/apiKeys';
import dayjs from 'dayjs';
import { Copy, Ellipsis, Loader2, Plus } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, Link } from 'react-router';
import { $path } from 'safe-routes';
import { useQueryState } from 'nuqs';
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
import clsx from 'clsx';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import type { ErrorResponse } from '@/api/instance';
import { Input } from '@/components/ui/input';

export default function ApiKeys() {
  const isInitialLoading = useRef(false);
  const { toast } = useToast();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const [search, setSearch] = useState('');
  const apiKeys = useSelector(apiKeysStore.selector.getApiKeys);
  const devKey = useSelector(apiKeysStore.selector.getDevKey);
  const [selectApiKeyId, setSelectApiKeyId] = useState<string | null>(null);
  const isDeleteApiKeyWaiting = useSelector(
    (state) => state.apiKeys.isDeleteApiKeyWaiting,
  );

  const [page, setPage] = useQueryState('page', {
    parse: (query: string) => parseInt(query, 10),
    serialize: (value) => value.toString(),
  });

  const pageSize = 10;

  const pageCount = Math.ceil(apiKeys.length / pageSize);

  const filteredApiKeys = useMemo(() => {
    if (!search) {
      return apiKeys;
    }
    const fuse = new Fuse(apiKeys, {
      keys: ['name'],
    });

    return fuse.search(search).map((result) => result.item);
  }, [apiKeys, search]);

  const slicedApiKeys = useMemo(() => {
    return filteredApiKeys.slice(
      (page || 1) * pageSize - pageSize,
      (page || 1) * pageSize,
    );
  }, [filteredApiKeys, page]);

  const onCopy = async (apiKey: string) => {
    try {
      await navigator.clipboard.writeText(apiKey);
      toast({
        title: t('routes_apikeys_toast_copied'),
      });
    } catch (error) {
      toast({
        title: t('routes_apikeys_toast_copy_failed'),
        description: (error as any)?.message ?? '',
        variant: 'destructive',
      });
    }
  };

  const handleDeleteApiKey = async () => {
    if (isDeleteApiKeyWaiting || !selectApiKeyId) {
      return;
    }
    const resultAction = await dispatch(
      apiKeysStore.asyncActions.deleteApiKey({ apiKeyId: selectApiKeyId }),
    );
    if (apiKeysStore.asyncActions.deleteApiKey.fulfilled.match(resultAction)) {
      setSelectApiKeyId(null);
      toast({
        title: t('routes_apikeys_toast_deleted'),
      });
    } else {
      toast({
        title: t('routes_apikeys_toast_delete_failed'),
        description: (resultAction.payload as ErrorResponse)?.detail ?? '',
        variant: 'destructive',
      });
    }
  };

  useEffect(() => {
    setBreadcrumbsState?.([{ label: t('breadcrumbs_api_keys') }]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        await dispatch(apiKeysStore.asyncActions.listApiKeys());
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch]);

  return (
    <div>
      <PageHeader label={t('routes_apikeys_page_header')} />
      <div className="flex w-screen flex-col gap-6 p-6 md:w-auto">
        <div className="flex flex-col justify-between gap-2 pt-6 md:flex-row">
          <p className="text-xl font-bold text-foreground">
            {t('routes_apikeys_personal_key_title')}
          </p>
        </div>

        <div className="w-full overflow-auto rounded-md border">
          <Table className="md:table-fixed">
            <TableHeader>
              <TableRow>
                <TableHead>{t('routes_apikeys_table_api_key')}</TableHead>
                <TableHead className="w-[72px]"></TableHead>
                <TableHead>{t('routes_apikeys_table_environment')}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {devKey && (
                <TableRow key={devKey.id}>
                  <TableCell>
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <span className="block truncate">{devKey.key}</span>
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>{devKey.key}</p>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  </TableCell>
                  <TableCell align="center">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="cursor-pointer"
                      onClick={() => onCopy(devKey.key)}
                    >
                      <Copy className="size-4" />
                    </Button>
                  </TableCell>
                  <TableCell>
                    <Badge
                      style={{ backgroundColor: devKey.environment.color }}
                    >
                      {devKey.environment.name}
                    </Badge>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>

        <div className="flex items-center justify-between gap-2 pt-6 md:flex-row">
          <p className="text-xl font-bold text-foreground">
            {t('routes_apikeys_live_mode_keys_title')}
          </p>
          {apiKeys.length === 0 && (
            <Button asChild>
              <Link to={$path('/apiKeys/new')}>
                <Plus />
                {t('routes_apikeys_create_new')}
              </Link>
            </Button>
          )}
        </div>

        {apiKeys.length === 0 && (
          <div className="flex flex-col gap-1 rounded-md border p-6">
            <h3 className="text-lg leading-7 font-bold">
              {t('routes_apikeys_no_keys_title')}
            </h3>
            <p className="text-sm font-normal text-muted-foreground">
              {t('routes_apikeys_no_keys_description')}
            </p>
          </div>
        )}

        {apiKeys.length > 0 && (
          <div className="flex justify-between gap-2">
            <div className="hidden max-w-72 flex-1 md:block">
              <Input
                placeholder={t('routes_apikeys_search_placeholder')}
                value={search}
                onChange={(e) => {
                  setSearch(e.target.value);
                }}
              />
            </div>
            <Button asChild>
              <Link to={$path('/apiKeys/new')}>
                <Plus />
                {t('routes_apikeys_create_new')}
              </Link>
            </Button>
          </div>
        )}

        {apiKeys.length > 0 && (
          <div className="w-full overflow-auto rounded-md border">
            <Table className="md:table-fixed">
              <TableHeader>
                <TableRow>
                  <TableHead>{t('routes_apikeys_table_name')}</TableHead>
                  <TableHead>{t('routes_apikeys_table_api_key')}</TableHead>
                  <TableHead className="w-[72px]"></TableHead>
                  <TableHead>{t('routes_apikeys_table_environment')}</TableHead>
                  <TableHead>{t('routes_apikeys_table_created_at')}</TableHead>
                  <TableHead className="w-16"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {slicedApiKeys.map((apiKey) => (
                  <TableRow
                    key={apiKey.id}
                    className="cursor-pointer"
                    onClick={() =>
                      navigate(
                        $path('/apiKeys/:apiKeyId', { apiKeyId: apiKey.id }),
                      )
                    }
                  >
                    <TableCell className="font-medium">{apiKey.name}</TableCell>
                    <TableCell>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <span className="block truncate">{apiKey.key}</span>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>{apiKey.key}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </TableCell>
                    <TableCell align="center">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="cursor-pointer"
                        onClick={(e) => {
                          e.stopPropagation();
                          onCopy(apiKey.key);
                        }}
                      >
                        <Copy className="size-4" />
                      </Button>
                    </TableCell>
                    <TableCell>
                      <Badge
                        style={{ backgroundColor: apiKey.environment.color }}
                      >
                        {apiKey.environment.name}
                      </Badge>
                    </TableCell>
                    <TableCell className="truncate">
                      {dayjs
                        .unix(Number(apiKey.createdAt))
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
                              setSelectApiKeyId(apiKey.id);
                            }}
                            className="text-destructive"
                          >
                            {t('routes_apikeys_delete')}
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
                    if (
                      index === (page || 1) + 2 ||
                      index === (page || 1) - 3
                    ) {
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

                  {page !== pageCount && pageCount > 1 && pageCount !== 1 && (
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
        )}
      </div>
      <Dialog
        open={selectApiKeyId !== null}
        onOpenChange={(boolean) => {
          if (isDeleteApiKeyWaiting) {
            return;
          }
          if (!boolean) {
            setSelectApiKeyId(null);
          }
        }}
      >
        <DialogContent
          onCloseAutoFocus={(event) => {
            event.preventDefault();
            document.body.style.pointerEvents = '';
          }}
          className="max-h-[calc(100svh-48px)] max-w-lg overflow-y-auto"
          data-state={selectApiKeyId !== null ? 'open' : 'closed'}
        >
          <DialogHeader>
            <DialogTitle>{t('routes_apikeys_delete_dialog_title')}</DialogTitle>
            <DialogDescription>
              {t('routes_apikeys_delete_dialog_description')}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="flex justify-end gap-2">
            <Button
              variant="outline"
              onClick={() => {
                if (isDeleteApiKeyWaiting) {
                  return;
                }
                setSelectApiKeyId(null);
              }}
            >
              {t('routes_apikeys_delete_dialog_cancel')}
            </Button>
            <Button
              onClick={handleDeleteApiKey}
              disabled={isDeleteApiKeyWaiting}
            >
              {isDeleteApiKeyWaiting && (
                <Loader2 className="size-4 animate-spin" />
              )}
              {t('routes_apikeys_delete_dialog_delete')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
