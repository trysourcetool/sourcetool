import type { UserRole } from '@/api/modules/users';
import { PageHeader } from '@/components/common/page-header';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarTrigger,
} from '@/components/ui/menubar';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useToast } from '@/hooks/use-toast';
import { useDispatch, useSelector } from '@/store';
import { groupsStore } from '@/store/modules/groups';
import { usersStore } from '@/store/modules/users';
import { zodResolver } from '@hookform/resolvers/zod';
import { Ellipsis, Loader2 } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { array, object, string } from 'zod';
import type { z } from 'zod';
import { createFileRoute, Link, useParams } from '@tanstack/react-router';

export default function UsersUserId() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const userId = useParams({
    from: '/_auth/users/$userId',
    select: (params) => params.userId,
  });
  const { toast } = useToast();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const [search, setSearch] = useState('');
  const { t } = useTranslation('common');
  const me = useSelector(usersStore.selector.getUserMe);
  const user = useSelector((state) =>
    usersStore.selector.getUser(state, userId ?? ''),
  );
  const isUpdateUserWaiting = useSelector(
    (state) => state.users.isUpdateUserWaiting,
  );
  const groups = useSelector((state) => groupsStore.selector.getGroups(state));
  const userGroups = useSelector((state) =>
    groupsStore.selector.getUserGroups(state),
  );

  const schema = object({
    role: string({
      required_error: t('zod_errors_role_required'),
    }),
    groupIds: array(string()).optional(),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const selectedGroupIds = form.watch('groupIds');

  console.log({ selectedGroupIds });

  const joinedGroups = useMemo(() => {
    const filterUserGroups = userGroups.map((userGroup) => userGroup.groupId);
    return groups.filter((group) => filterUserGroups.includes(group.id));
  }, [groups, userGroups]);

  const groupOptions = useMemo(() => {
    return groups.filter(
      (group) =>
        !(selectedGroupIds || []).some((groupId) => groupId === group.id),
    );
  }, [selectedGroupIds, groups]);

  const onSubmit = form.handleSubmit(async (data) => {
    if (isUpdateUserWaiting || !user) {
      return;
    }
    const resultAction = await dispatch(
      usersStore.asyncActions.updateUser({
        userId: user.id,
        data: {
          role: data.role as UserRole,
          groupIds: data.groupIds || [],
        },
      }),
    );
    if (usersStore.asyncActions.updateUser.fulfilled.match(resultAction)) {
      toast({
        title: t('routes_users_edit_toast_updated'),
      });
      dispatch(groupsStore.asyncActions.listGroups());
    } else {
      toast({
        title: t('routes_users_edit_toast_update_failed'),
        variant: 'destructive',
      });
    }
  });

  useEffect(() => {
    if (!isInitialLoading.current) {
      (async () => {
        isInitialLoading.current = true;
        await Promise.all([dispatch(groupsStore.asyncActions.listGroups())]);
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch]);

  useEffect(() => {
    setBreadcrumbsState?.([
      { label: t('breadcrumbs_users'), to: '/users' },
      { label: `${user?.firstName} ${user?.lastName}` },
    ]);
  }, [setBreadcrumbsState, user, t]);

  useEffect(() => {
    if (user) {
      console.log({ joinedGroups, user });
      form.reset({
        role: user?.role,
        groupIds: joinedGroups.map((group) => group.id),
      });
    }
  }, [user, joinedGroups]);

  return (
    user && (
      <div>
        <PageHeader
          label={`${user.firstName} ${user.lastName}`}
          description={user.email}
        />
        <div className="flex flex-col gap-6 p-6">
          <Form {...form}>
            <FormField
              control={form.control}
              name="role"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('routes_users_edit_role_label')}</FormLabel>
                  <FormControl>
                    <Select
                      value={field.value}
                      onValueChange={field.onChange}
                      disabled={user.id === me?.id || me?.role !== 'admin'}
                      defaultValue={user.role}
                    >
                      <SelectTrigger>
                        <SelectValue
                          placeholder={t('routes_users_edit_role_placeholder')}
                        />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="admin">
                          {t('routes_users_edit_role_admin')}
                        </SelectItem>
                        <SelectItem value="member">
                          {t('routes_users_edit_role_member')}
                        </SelectItem>
                        <SelectItem value="developer">
                          {t('routes_users_edit_role_developer')}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Separator />

            <div className="flex flex-col gap-2">
              <h2 className="text-2xl font-bold">
                {t('routes_users_edit_permission_groups_title')}
              </h2>
              <FormField
                control={form.control}
                name="groupIds"
                render={({ field }) => (
                  <div className="flex flex-col gap-6">
                    <div className="flex">
                      <Menubar>
                        <MenubarMenu>
                          <MenubarTrigger>
                            {t('routes_users_edit_search_group')}
                          </MenubarTrigger>
                          <MenubarContent>
                            <Input
                              placeholder={t(
                                'routes_users_edit_search_placeholder',
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

                            {groupOptions.map((group) => (
                              <MenubarItem
                                key={group.id}
                                onSelect={() => {
                                  field.onChange([
                                    ...(field.value || []),
                                    group.id,
                                  ]);
                                }}
                              >
                                <span className="flex-1">{group.name}</span>
                                <Button variant="ghost" size="sm">
                                  {t('routes_users_edit_add_button')}
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
                              {t('routes_users_edit_table_name')}
                            </TableHead>
                            <TableHead>
                              {t('routes_users_edit_table_slug')}
                            </TableHead>
                            <TableHead className="w-16"></TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {groups
                            .filter((group) =>
                              selectedGroupIds?.includes(group.id),
                            )
                            .map((group) => (
                              <TableRow key={group.id}>
                                <TableCell className="font-medium">
                                  {group.name}
                                </TableCell>
                                <TableCell>{group.slug}</TableCell>

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
                                            field.value?.filter(
                                              (id) => id !== group.id,
                                            ),
                                          )
                                        }
                                      >
                                        {t('routes_users_edit_remove')}
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
            </div>

            <div>
              <div className="flex flex-col justify-start gap-3 md:flex-row">
                <Button
                  type="button"
                  disabled={isUpdateUserWaiting}
                  onClick={onSubmit}
                >
                  {isUpdateUserWaiting && (
                    <Loader2 className="size-4 animate-spin" />
                  )}
                  {t('routes_users_edit_save_button')}
                </Button>
                <Button
                  variant="outline"
                  asChild
                  disabled={isUpdateUserWaiting}
                >
                  <Link to={'/users'}>
                    {t('routes_users_edit_cancel_button')}
                  </Link>
                </Button>
              </div>
            </div>
          </Form>
        </div>
      </div>
    )
  );
}

export const Route = createFileRoute('/_auth/users/$userId')({
  component: UsersUserId,
});
