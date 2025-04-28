import { PageHeader } from '@/components/common/page-header';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
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
import { Separator } from '@/components/ui/separator';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { useDispatch, useSelector } from '@/store';
import { usersStore } from '@/store/modules/users';
import { Ellipsis, Loader2, Plus } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import {
  createFileRoute,
  Link,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { boolean, number, object, string } from 'zod';
import type { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { useToast } from '@/hooks/use-toast';
import { Badge } from '@/components/ui/badge';
import type { User, UserInvitation, UserRole } from '@/api/modules/users';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
// import { useQueryState } from 'nuqs';
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
import { zodValidator } from '@tanstack/zod-adapter';

const InviteForm = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { t } = useTranslation('common');

  const isInviteWaiting = useSelector(
    (state) => state.users.isCreateUserInvitationsWaiting,
  );

  const schema = object({
    emails: string({
      required_error: t('zod_errors_email_required'),
    }).refine(
      (value) =>
        value
          .split(',')
          .every((email) =>
            email.trim().match(/[\w\-._]+@[\w\-._]+\.[A-Za-z]+/),
          ),
      {
        message: t('zod_errors_email_format'),
      },
    ),
    role: string({
      required_error: t('zod_errors_role_required'),
    }).refine((value) => ['admin', 'member', 'developer'].includes(value), {
      message: t('zod_errors_invalid_role'),
    }),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
    defaultValues: {
      emails: '',
      role: 'admin',
    },
  });

  const emailsArray = form
    .watch('emails')
    .split(',')
    .filter((email) => email.trim() !== '');

  const membersCount = emailsArray.length;

  const onSubmit = form.handleSubmit(async (data) => {
    if (isInviteWaiting) {
      return;
    }
    const resultAction = await dispatch(
      usersStore.asyncActions.createUserInvitations({
        data: {
          emails: data.emails.split(',').map((email) => email.trim()),
          role: data.role as UserRole,
        },
      }),
    );

    if (
      usersStore.asyncActions.createUserInvitations.fulfilled.match(
        resultAction,
      )
    ) {
      navigate({ to: '/users', search: { invite: undefined } });
      toast({
        title: t('routes_users_toast_invited'),
        description: t('routes_users_toast_invited_description'),
      });
      dispatch(usersStore.asyncActions.listUsers());
    } else {
      toast({
        title: t('routes_users_toast_invite_failed'),
        description: t('routes_users_toast_invite_failed_description'),
        variant: 'destructive',
      });
    }
  });

  console.log(form.formState.errors);

  return (
    <Form {...form}>
      <form onSubmit={onSubmit} className="block space-y-4">
        <FormField
          control={form.control}
          name="emails"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('routes_users_invite_email_label')}</FormLabel>
              <FormControl>
                <Input
                  placeholder={t('routes_users_invite_email_placeholder')}
                  {...field}
                />
              </FormControl>
              <p className="text-muted-foreground text-sm">
                {t('routes_users_invite_email_help')}
              </p>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="role"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('routes_users_invite_role_label')}</FormLabel>
              <FormControl>
                <Select {...field} onValueChange={field.onChange}>
                  <SelectTrigger>
                    <SelectValue
                      placeholder={t('routes_users_invite_role_placeholder')}
                    />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="admin">
                      {t('routes_users_invite_role_admin')}
                    </SelectItem>
                    <SelectItem value="member">
                      {t('routes_users_invite_role_member')}
                    </SelectItem>
                    <SelectItem value="developer">
                      {t('routes_users_invite_role_developer')}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Separator />
        <Button type="submit" disabled={isInviteWaiting}>
          {isInviteWaiting && <Loader2 className="size-4 animate-spin" />}
          {t('routes_users_invite_button_text', { count: membersCount })}
        </Button>
      </form>
    </Form>
  );
};

function Users() {
  const isInitialLoading = useRef(false);
  const dispatch = useDispatch();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const { toast } = useToast();
  const [search, setSearch] = useState('');
  const searchParams = useSearch({ from: '/_auth/users/' });
  const page = searchParams.page || 1;
  const invite = searchParams.invite || false;
  const [userToDelete, setUserToDelete] = useState<User | null>(null);

  const users = useSelector(usersStore.selector.getUsers);
  const userInvitations = useSelector(usersStore.selector.getUserInvitations);
  const me = useSelector(usersStore.selector.getUserMe);
  const isInviteWaiting = useSelector(
    (state) => state.users.isCreateUserInvitationsWaiting,
  );
  const isInvitationsResendWaiting = useSelector(
    (state) => state.users.isResendUserInvitationWaiting,
  );
  const isDeleteUserWaiting = useSelector(
    (state) => state.users.isDeleteUserWaiting,
  );

  const filteredUsers = useMemo(() => {
    if (!search) {
      return users;
    }
    const fuse = new Fuse(users, {
      keys: ['firstName', 'lastName', 'email'],
    });

    return fuse.search(search).map((result) => result.item);
  }, [users, search]);

  const filteredUserInvitations = useMemo(() => {
    if (!search) {
      return userInvitations;
    }
    const fuse = new Fuse(userInvitations, {
      keys: ['email'],
    });

    return fuse.search(search).map((result) => result.item);
  }, [userInvitations, search]);

  const pageSize = 10;

  const pageCount = Math.ceil(
    (filteredUsers.length + filteredUserInvitations.length) / pageSize,
  );

  const slicedUsers = useMemo(() => {
    return [
      ...filteredUserInvitations.map((userInvitation) => ({
        ...userInvitation,
        type: 'invitation',
      })),
      ...filteredUsers.map((user) => ({
        ...user,
        type: 'active',
      })),
    ].slice((page || 1) * pageSize - pageSize, (page || 1) * pageSize) as (
      | ({
          type: 'invitation';
        } & UserInvitation)
      | ({
          type: 'active';
        } & User)
    )[];
  }, [filteredUsers, filteredUserInvitations, page]);

  const handleResendInvitation = async (invitationId: string) => {
    const resultAction = await dispatch(
      usersStore.asyncActions.resendUserInvitation({
        invitationId,
      }),
    );

    if (
      usersStore.asyncActions.resendUserInvitation.fulfilled.match(resultAction)
    ) {
      toast({
        title: t('routes_users_toast_invitation_resent'),
        description: t('routes_users_toast_invitation_resent_description'),
      });
    } else {
      toast({
        title: t('routes_users_toast_invitation_resend_failed'),
        description: t(
          'routes_users_toast_invitation_resend_failed_description',
        ),
        variant: 'destructive',
      });
    }
  };

  const handleDeleteUser = async () => {
    if (!userToDelete || isDeleteUserWaiting) {
      return;
    }
    const resultAction = await dispatch(
      usersStore.asyncActions.deleteUser({
        userId: userToDelete.id,
      }),
    );

    if (usersStore.asyncActions.deleteUser.fulfilled.match(resultAction)) {
      toast({
        title: t('routes_users_toast_user_deleted'),
        description: t('routes_users_toast_user_deleted_description', {
          name: `${userToDelete.firstName} ${userToDelete.lastName}`.trim(),
        }),
      });
      dispatch(usersStore.asyncActions.listUsers()); // Refresh the list
    } else {
      toast({
        title: t('routes_users_toast_user_delete_failed'),
        description: t('routes_users_toast_user_delete_failed_description'),
        variant: 'destructive',
      });
    }
    setUserToDelete(null); // Close dialog
  };

  useEffect(() => {
    setBreadcrumbsState?.([{ label: t('breadcrumbs_users') }]);
  }, [setBreadcrumbsState, t]);

  useEffect(() => {
    if (!isInitialLoading.current) {
      (async () => {
        isInitialLoading.current = true;
        await dispatch(usersStore.asyncActions.listUsers());
        isInitialLoading.current = false;
      })();
    }
  }, [dispatch]);

  return (
    <div>
      <PageHeader label={t('routes_users_page_header')} />
      <div className="flex flex-col gap-4 px-4 py-6 md:gap-6 md:px-6">
        <div className="flex flex-col justify-between gap-2 md:flex-row md:pt-6">
          <p className="text-foreground text-xl font-bold">
            {t('routes_users_title')}
          </p>
        </div>

        <div className="flex flex-col-reverse justify-between gap-3 md:flex-row md:gap-2">
          <div className="w-full max-w-full flex-1 md:w-auto md:max-w-72">
            <Input
              placeholder={t('routes_users_search_placeholder')}
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
              }}
            />
          </div>
          <div>
            <Button asChild>
              <Link to={'/users'} search={{ invite: true }}>
                <Plus />
                {t('routes_users_invite_button')}
              </Link>
            </Button>
          </div>
        </div>

        <div className="w-full overflow-auto rounded-md border">
          <Table className="md:table-fixed">
            <TableHeader>
              <TableRow>
                <TableHead>{t('routes_users_table_name')}</TableHead>
                <TableHead>{t('routes_users_table_email')}</TableHead>
                <TableHead>{t('routes_users_table_permission')}</TableHead>
                <TableHead></TableHead>
                <TableHead className="w-16"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {slicedUsers.map((user) => {
                if (user.type === 'invitation') {
                  return (
                    <TableRow key={user.id}>
                      <TableCell className="text-muted-foreground truncate font-normal">
                        {t('routes_users_invitation_sent')}
                      </TableCell>
                      <TableCell className="text-muted-foreground truncate font-normal">
                        {user.email}
                      </TableCell>
                      <TableCell></TableCell>
                      <TableCell>
                        <Button
                          variant={'outline'}
                          type="button"
                          className="cursor-pointer"
                          onClick={() => handleResendInvitation(user.id)}
                          disabled={isInvitationsResendWaiting}
                        >
                          {t('routes_users_resend_invitation')}
                        </Button>
                      </TableCell>
                      <TableCell></TableCell>
                    </TableRow>
                  );
                }
                if (user.type === 'active') {
                  return (
                    <TableRow
                      key={user.id}
                      className="cursor-pointer"
                      onClick={() => {
                        navigate({
                          to: '/users/$userId',
                          params: {
                            userId: user.id,
                          },
                        });
                      }}
                    >
                      <TableCell className="truncate font-medium">
                        {user.firstName} {user.lastName}
                      </TableCell>
                      <TableCell className="truncate">{user.email}</TableCell>
                      <TableCell>
                        <Badge
                          variant="secondary"
                          className="text-xs"
                          style={{
                            textTransform: 'capitalize',
                          }}
                        >
                          {user.role}
                        </Badge>
                      </TableCell>
                      <TableCell></TableCell>
                      <TableCell align="right">
                        {me?.role === 'admin' && me?.id !== user.id && (
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon">
                                <Ellipsis className="size-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem
                                onClick={(e) => {
                                  e.stopPropagation(); // Prevent row click
                                  setUserToDelete(user);
                                }}
                              >
                                {t('routes_users_delete')}
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        )}
                      </TableCell>
                    </TableRow>
                  );
                }
              })}
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
                          to: '.',
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
                          to: '.',
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
                  if (index === (page || 1) + 2 || index === (page || 1) - 3) {
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
                          navigate({ to: '.', search: { page: index + 1 } })
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

                {page !== pageCount && pageCount > 1 && pageCount !== 1 && (
                  <PaginationItem>
                    <PaginationNext
                      onClick={() => {
                        if (page === pageCount) {
                          return;
                        }
                        navigate({
                          to: '.',
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
                          to: '.',
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
      </div>

      <Dialog
        open={invite}
        onOpenChange={(open) => {
          if (!open && !isInviteWaiting) {
            navigate({ to: '/users', search: { invite: undefined } });
          }
        }}
      >
        <DialogContent
          className="max-h-[calc(100svh-48px)] max-w-lg overflow-y-auto"
          data-state={invite ? 'open' : 'closed'}
        >
          <DialogHeader>
            <DialogTitle className="text-2xl font-bold">
              {t('routes_users_invite_dialog_title')}
            </DialogTitle>
            <DialogDescription>
              {t('routes_users_invite_dialog_description')}
            </DialogDescription>
          </DialogHeader>
          <InviteForm />
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={!!userToDelete}
        onOpenChange={(open) => {
          if (!open) {
            setUserToDelete(null);
          }
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {t('routes_users_delete_confirm_title')}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {t('routes_users_delete_confirm_description', {
                name: `${userToDelete?.firstName} ${userToDelete?.lastName}`.trim(),
                email: userToDelete?.email,
              })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={isDeleteUserWaiting}>
              {t('common_cancel')}
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteUser}
              disabled={isDeleteUserWaiting}
              className="bg-foreground text-background hover:bg-foreground/90"
            >
              {isDeleteUserWaiting && (
                <Loader2 className="mr-2 size-4 animate-spin" />
              )}
              {t('routes_users_delete_confirm_button')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}

export const Route = createFileRoute('/_auth/users/')({
  component: Users,
  validateSearch: zodValidator(
    object({
      page: number().optional(),
      invite: boolean().optional(),
    }),
  ),
});
