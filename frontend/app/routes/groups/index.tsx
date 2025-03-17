import { PageHeader } from '@/components/common/page-header';
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
import { groupsStore } from '@/store/modules/groups';
import dayjs from 'dayjs';
import { Ellipsis, Plus } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { Link, useNavigate } from 'react-router';
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
import { Input } from '@/components/ui/input';

export default function Groups() {
  const isInitialLoading = useRef(false);
  const [isInitialLoaded, setIsInitialLoaded] = useState(false);
  const { toast } = useToast();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const [search, setSearch] = useState('');
  const groups = useSelector(groupsStore.selector.getGroups);
  const userGroups = useSelector(groupsStore.selector.getUserGroups);
  const isDeleteGroupWaiting = useSelector(
    (state) => state.groups.isDeleteGroupWaiting,
  );

  const [page, setPage] = useQueryState('page', {
    parse: (query: string) => parseInt(query, 10),
    serialize: (value) => value.toString(),
  });

  const pageSize = 10;

  const pageCount = Math.ceil(groups.length / pageSize);

  const filteredGroups = useMemo(() => {
    if (!search) {
      return groups;
    }
    const fuse = new Fuse(groups, {
      keys: ['name'],
    });
    return fuse.search(search).map((result) => result.item);
  }, [groups, search]);

  const slicedGroups = useMemo(() => {
    return filteredGroups.slice(
      (page || 1) * pageSize - pageSize,
      (page || 1) * pageSize,
    );
  }, [filteredGroups, page]);

  const countGroupUsers = (groupId: string) => {
    return userGroups.filter((userGroup) => userGroup.groupId === groupId)
      .length;
  };

  const onDelete = async (groupId: string) => {
    if (isDeleteGroupWaiting) {
      return;
    }
    const resultAction = await dispatch(
      groupsStore.asyncActions.deleteGroup({ groupId }),
    );
    if (groupsStore.asyncActions.deleteGroup.fulfilled.match(resultAction)) {
      toast({
        title: t('routes_groups_toast_deleted'),
      });
    } else {
      toast({
        title: t('routes_groups_toast_delete_failed'),
        description: (resultAction.error as any)?.message ?? '',
        variant: 'destructive',
      });
    }
  };

  useEffect(() => {
    setBreadcrumbsState?.([{ label: t('breadcrumbs_permissions') }]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        await dispatch(groupsStore.asyncActions.listGroups());
        isInitialLoading.current = false;
        setIsInitialLoaded(true);
      })();
    }
  }, [dispatch]);

  return (
    <div>
      <PageHeader label={t('routes_groups_page_header')} />
      <div className="flex w-screen flex-col gap-6 p-6 md:w-auto">
        <div className="flex items-center justify-between gap-2 pt-6 md:flex-row">
          <p className="text-xl font-bold text-foreground">
            {t('routes_groups_title')}
          </p>
          {groups.length === 0 && (
            <Button asChild>
              <Link to={$path('/groups/new')}>
                <Plus />
                {t('routes_groups_create_new')}
              </Link>
            </Button>
          )}
        </div>

        {isInitialLoaded && groups.length > 0 && (
          <div className="flex justify-between gap-2">
            <div className="hidden max-w-72 flex-1 md:block">
              <Input
                placeholder={t('routes_groups_search_placeholder')}
                value={search}
                onChange={(e) => {
                  setSearch(e.target.value);
                }}
              />
            </div>
            <Button asChild>
              <Link to={$path('/groups/new')}>
                <Plus />
                {t('routes_groups_create_new')}
              </Link>
            </Button>
          </div>
        )}

        {groups.length === 0 && (
          <div className="flex flex-col gap-1 rounded-md border p-6">
            <h3 className="text-lg leading-7 font-bold">
              {t('routes_groups_no_groups_title')}
            </h3>
            <p className="text-sm font-normal text-muted-foreground">
              {t('routes_groups_no_groups_description')}
            </p>
          </div>
        )}

        {groups.length > 0 && (
          <div className="w-full overflow-auto rounded-md border">
            <Table className="md:table-fixed">
              <TableHeader>
                <TableRow>
                  <TableHead>{t('routes_groups_table_name')}</TableHead>
                  <TableHead>{t('routes_groups_table_slug')}</TableHead>
                  <TableHead>{t('routes_groups_table_created_at')}</TableHead>
                  <TableHead>{t('routes_groups_table_users')}</TableHead>
                  <TableHead className="w-16"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {slicedGroups.map((group) => (
                  <TableRow
                    key={group.id}
                    className="cursor-pointer"
                    onClick={() =>
                      navigate($path('/groups/:groupId', { groupId: group.id }))
                    }
                  >
                    <TableCell className="font-medium">{group.name}</TableCell>
                    <TableCell>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <span className="block truncate">{group.slug}</span>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>{group.slug}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </TableCell>

                    <TableCell className="truncate">
                      {dayjs
                        .unix(Number(group.createdAt))
                        .format('DD MMM YYYY')}
                    </TableCell>
                    <TableCell className="truncate">
                      {t('routes_groups_users_count', {
                        count: countGroupUsers(group.id),
                      })}
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
                              onDelete(group.id);
                            }}
                            className="text-destructive"
                          >
                            {t('routes_groups_delete')}
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
        )}
      </div>
    </div>
  );
}
