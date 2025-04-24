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
import {
  createFileRoute,
  Link,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
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
import { zodValidator } from '@tanstack/zod-adapter';
import { number, object } from 'zod';

export default function Groups() {
  const isInitialLoading = useRef(false);
  const [isInitialLoaded, setIsInitialLoaded] = useState(false);
  const { toast } = useToast();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const [search, setSearch] = useState('');
  const searchParams = useSearch({ from: '/_auth/groups/' });
  const page = searchParams.page || 1;
  const groups = useSelector(groupsStore.selector.getGroups);
  const userGroups = useSelector(groupsStore.selector.getUserGroups);
  const isDeleteGroupWaiting = useSelector(
    (state) => state.groups.isDeleteGroupWaiting,
  );

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
      <div className="flex w-screen flex-col gap-4 px-4 py-6 md:w-auto md:gap-6 md:px-6">
        <div className="flex flex-col items-start justify-between gap-4 md:flex-row md:pt-6">
          <p className="text-foreground text-xl font-bold">
            {t('routes_groups_title')}
          </p>
          {groups.length === 0 && (
            <Button asChild>
              <Link to={'/groups/new'}>
                <Plus />
                {t('routes_groups_create_new')}
              </Link>
            </Button>
          )}
        </div>

        {isInitialLoaded && groups.length > 0 && (
          <div className="flex flex-col-reverse justify-between gap-2 md:flex-row">
            <div className="w-full max-w-full flex-1 md:w-auto md:max-w-72">
              <Input
                placeholder={t('routes_groups_search_placeholder')}
                value={search}
                onChange={(e) => {
                  setSearch(e.target.value);
                }}
              />
            </div>
            <div>
              <Button asChild>
                <Link to={'/groups/new'}>
                  <Plus />
                  {t('routes_groups_create_new')}
                </Link>
              </Button>
            </div>
          </div>
        )}

        {groups.length === 0 && (
          <div className="flex flex-col gap-1 rounded-md border p-6">
            <h3 className="text-lg font-bold leading-7">
              {t('routes_groups_no_groups_title')}
            </h3>
            <p className="text-muted-foreground text-sm font-normal">
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
                      navigate({
                        to: '/groups/$groupId',
                        params: { groupId: group.id },
                      })
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
            <div className="bg-muted border-t p-4">
              <Pagination className="flex justify-end">
                <PaginationContent>
                  {page !== 1 && page !== null && (
                    <PaginationItem>
                      <PaginationPrevious
                        className={clsx(
                          (page !== 1 || page === null) && 'cursor-pointer',
                          'hidden md:inline-flex',
                        )}
                        onClick={() => {
                          if (page === 1 || page === null) {
                            return;
                          }
                          navigate({
                            to: '/groups',
                            search: { page: (page || 1) - 1 },
                          });
                        }}
                      />
                      <PaginationLink
                        onClick={() => {
                          if (page === 1 || page === null) {
                            return;
                          }
                          navigate({
                            to: '/groups',
                            search: { page: (page || 1) - 1 },
                          });
                        }}
                        isActive
                        className={clsx(
                          page !== pageCount && 'cursor-pointer',
                          'inline-flex w-auto px-4 md:hidden',
                        )}
                        aria-label="Go to previous page"
                      >
                        <span>Previous</span>
                      </PaginationLink>
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
                        <PaginationItem key={index} className="hidden md:block">
                          <PaginationEllipsis />
                        </PaginationItem>
                      );
                    }
                    return (
                      <PaginationItem key={index} className="hidden md:block">
                        <PaginationLink
                          onClick={() =>
                            navigate({
                              to: '/groups',
                              search: { page: index + 1 },
                            })
                          }
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
                          navigate({
                            to: '/groups',
                            search: { page: (page || 1) + 1 },
                          });
                        }}
                        className={clsx(
                          page !== pageCount && 'cursor-pointer',
                          'hidden md:inline-flex',
                        )}
                      />
                      <PaginationLink
                        onClick={() => {
                          if (page === pageCount) {
                            return;
                          }
                          navigate({
                            to: '/groups',
                            search: { page: (page || 1) + 1 },
                          });
                        }}
                        isActive
                        className={clsx(
                          page !== pageCount && 'cursor-pointer',
                          'inline-flex w-auto px-4 md:hidden',
                        )}
                        aria-label="Go to next page"
                      >
                        <span>Next</span>
                      </PaginationLink>
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

export const Route = createFileRoute('/_auth/groups/')({
  component: Groups,
  validateSearch: zodValidator(
    object({
      page: number().optional(),
    }),
  ),
});
