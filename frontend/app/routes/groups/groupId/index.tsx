import { PageHeader } from '@/components/common/page-header';
import { array, object, string } from 'zod';
import type { z } from 'zod';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from '@/store';
import { useEffect, useMemo, useRef, useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import Fuse from 'fuse.js';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';

import { Button } from '@/components/ui/button';
import { Link, useParams } from 'react-router';
import { $path } from 'safe-routes';
import { Ellipsis, Loader2 } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { groupsStore } from '@/store/modules/groups';
import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarTrigger,
} from '@/components/ui/menubar';
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
import { Badge } from '@/components/ui/badge';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { usersStore } from '@/store/modules/users';

const searchOptions = {
  includeMatches: true,
  includeScore: true,
  keys: ['firstName', 'lastName'],
  threshold: 0.3,
};

export default function GroupEdit() {
  const isInitialLoading = useRef(false);
  const { toast } = useToast();
  const dispatch = useDispatch();
  const { groupId } = useParams();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const [search, setSearch] = useState('');
  const { t } = useTranslation('common');
  const users = useSelector(usersStore.selector.getUsers);

  const group = useSelector((state) =>
    groupsStore.selector.getGroup(state, groupId ?? ''),
  );
  const userGroups = useSelector(groupsStore.selector.getUserGroups);
  const isUpdateGroupWaiting = useSelector(
    (state) => state.groups.isUpdateGroupWaiting,
  );

  const schema = object({
    name: string({
      required_error: t('zod_errors_name_required'),
    }),
    userIds: array(string()),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const userIds = form.watch('userIds');

  const filteredUsers = useMemo(() => {
    console.log({ userIds });
    const filteredUsers = users.filter((user) => {
      return !userIds?.includes(user.id);
    });
    if (!search) {
      return filteredUsers;
    }
    const fuse = new Fuse(filteredUsers, searchOptions);
    return fuse.search(search).map((result) => result.item);
  }, [users, search, userIds]);

  const selectedUsers = useMemo(() => {
    return users.filter((user) => userIds?.includes(user.id));
  }, [users, userIds]);

  const onSubmit = form.handleSubmit(async (data) => {
    console.log({ data });
    if (isUpdateGroupWaiting || !groupId) {
      return;
    }
    const resultAction = await dispatch(
      groupsStore.asyncActions.updateGroup({
        groupId,
        data: {
          name: data.name,
          userIds: data.userIds || [],
        },
      }),
    );
    if (groupsStore.asyncActions.updateGroup.fulfilled.match(resultAction)) {
      toast({
        title: t('routes_groups_edit_toast_updated'),
      });
      await dispatch(groupsStore.asyncActions.listGroups());
    } else {
      toast({
        title: t('routes_groups_edit_toast_update_failed'),
        description: (resultAction.error as any)?.message ?? '',
        variant: 'destructive',
      });
    }
  });

  useEffect(() => {
    setBreadcrumbsState?.([
      { label: t('breadcrumbs_groups'), to: $path('/groups') },
      { label: t('breadcrumbs_edit_group') },
    ]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      isInitialLoading.current = true;
      (async () => {
        await dispatch(groupsStore.asyncActions.listGroups());
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch, groupId]);

  useEffect(() => {
    if (group) {
      form.reset({
        name: group.name,
        userIds: userGroups
          .filter((userGroup) => userGroup.groupId === groupId)
          .map((userGroup) => userGroup.userId),
      });
    }
  }, [group?.id, userGroups]);

  return (
    <div>
      <PageHeader label={t('routes_groups_edit_page_header')} />
      {group && (
        <Form {...form}>
          <form
            className="flex flex-col gap-6 px-4 py-6 md:px-6"
            onSubmit={onSubmit}
          >
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('routes_groups_edit_name_label')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('routes_groups_edit_name_placeholder')}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormItem>
              <FormLabel>{t('routes_groups_edit_slug_label')}</FormLabel>
              <Input value={group.slug} disabled={true} />
            </FormItem>

            <FormField
              control={form.control}
              name="userIds"
              render={({ field }) => (
                <div className="flex flex-col gap-6">
                  <p className="text-xl font-bold text-foreground">
                    {t('routes_groups_edit_users_title')}
                  </p>
                  <div className="flex">
                    <Menubar>
                      <MenubarMenu>
                        <MenubarTrigger>
                          {t('routes_groups_edit_search_user')}
                        </MenubarTrigger>
                        <MenubarContent>
                          <Input
                            placeholder={t(
                              'routes_groups_edit_search_placeholder',
                            )}
                            value={search}
                            onChange={(e) => {
                              e.stopPropagation();
                              e.preventDefault();
                              setSearch(e.target.value);
                            }}
                            onKeyDown={(e) => {
                              e.stopPropagation();
                            }}
                          />

                          {filteredUsers.map((user) => (
                            <MenubarItem
                              key={user.id}
                              onSelect={() => {
                                field.onChange([
                                  ...(field.value || []),
                                  user.id,
                                ]);
                              }}
                            >
                              <span className="flex-1">
                                {user.firstName} {user.lastName}
                              </span>
                              <Button variant="ghost" size="sm">
                                {t('routes_groups_edit_add_button')}
                              </Button>
                            </MenubarItem>
                          ))}
                        </MenubarContent>
                      </MenubarMenu>
                    </Menubar>
                  </div>
                  <div className="rounded-md border">
                    <Table className="md:table-fixed">
                      <TableHeader>
                        <TableRow>
                          <TableHead>
                            {t('routes_groups_edit_table_name')}
                          </TableHead>
                          <TableHead>
                            {t('routes_groups_edit_table_email')}
                          </TableHead>
                          <TableHead>
                            {t('routes_groups_edit_table_permission')}
                          </TableHead>
                          <TableHead className="w-16"></TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {selectedUsers.map((user) => (
                          <TableRow key={user.id}>
                            <TableCell className="truncate font-medium">
                              {user.firstName} {user.lastName}
                            </TableCell>
                            <TableCell className="truncate">
                              <TooltipProvider>
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <span className="block truncate">
                                      {user.email}
                                    </span>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p>{user.email}</p>
                                  </TooltipContent>
                                </Tooltip>
                              </TooltipProvider>
                            </TableCell>

                            <TableCell className="truncate">
                              <Badge variant="secondary">{user.role}</Badge>
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
                                    onClick={() =>
                                      field.onChange(
                                        field.value.filter(
                                          (id) => id !== user.id,
                                        ),
                                      )
                                    }
                                  >
                                    {t('routes_groups_edit_remove')}
                                  </DropdownMenuItem>
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </div>
                </div>
              )}
            />

            <div className="flex flex-row justify-start gap-3">
              <Button type="submit">
                {isUpdateGroupWaiting && (
                  <Loader2 className="size-4 animate-spin" />
                )}
                {t('routes_groups_edit_save_button')}
              </Button>
              <Button variant="outline" asChild>
                <Link to={$path('/groups')}>
                  {t('routes_groups_edit_cancel_button')}
                </Link>
              </Button>
            </div>
          </form>
        </Form>
      )}
    </div>
  );
}
