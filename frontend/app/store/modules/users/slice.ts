import {
  createEntityAdapter,
  createSlice,
  type EntityState,
} from '@reduxjs/toolkit';
import * as asyncActions from './asyncActions';
import type { User, UserInvitation } from '@/api/modules/users';
import { groupsStore } from '../groups';
import { pagesStore } from '../pages';
import { organizationsStore } from '../organizations';

// =============================================
// schema

const usersAdapter = createEntityAdapter<User, string>({
  selectId: (user) => user.id,
});

const userInvitationsAdapter = createEntityAdapter<UserInvitation, string>({
  selectId: (userInvitation) => userInvitation.id,
});

// =============================================
// State

export type State = {
  me: User | null;
  users: EntityState<User, string>;
  userInvitations: EntityState<UserInvitation, string>;
  isGetMeWaiting: boolean;
  isRefreshTokenWaiting: boolean;
  isListUsersWaiting: boolean;
  isRequestMagicLinkWaiting: boolean;
  isRegisterWithMagicLinkWaiting: boolean;
  isAuthenticateWithMagicLinkWaiting: boolean;
  isAuthChecked: boolean;
  isAuthSucceeded: boolean;
  isAuthFailed: boolean;
  isInviteWaiting: boolean;
  isSaveAuthWaiting: boolean;
  isSignoutWaiting: boolean;
  isObtainAuthTokenWaiting: boolean;
  isOauthGoogleAuthWaiting: boolean;
  isUpdateUserWaiting: boolean;
  isUpdateUserEmailWaiting: boolean;
  isUsersSendUpdateEmailInstructionsWaiting: boolean;
  isInvitationsResendWaiting: boolean;
  isRequestInvitationMagicLinkWaiting: boolean;
  isAuthenticateWithInvitationMagicLinkWaiting: boolean;
  isRegisterWithInvitationMagicLinkWaiting: boolean;
};

const initialState: State = {
  me: null,
  users: usersAdapter.getInitialState(),
  userInvitations: userInvitationsAdapter.getInitialState(),
  isGetMeWaiting: false,
  isRefreshTokenWaiting: false,
  isListUsersWaiting: false,
  isRequestMagicLinkWaiting: false,
  isRegisterWithMagicLinkWaiting: false,
  isAuthenticateWithMagicLinkWaiting: false,
  isAuthChecked: false,
  isAuthSucceeded: false,
  isAuthFailed: false,
  isInviteWaiting: false,
  isObtainAuthTokenWaiting: false,
  isSaveAuthWaiting: false,
  isSignoutWaiting: false,
  isOauthGoogleAuthWaiting: false,
  isUpdateUserWaiting: false,
  isUpdateUserEmailWaiting: false,
  isUsersSendUpdateEmailInstructionsWaiting: false,
  isInvitationsResendWaiting: false,
  isRequestInvitationMagicLinkWaiting: false,
  isAuthenticateWithInvitationMagicLinkWaiting: false,
  isRegisterWithInvitationMagicLinkWaiting: false,
};
// =============================================
// slice

export const slice = createSlice({
  extraReducers: (builder) => {
    builder
      // refreshToken
      .addCase(asyncActions.refreshToken.pending, (state) => {
        state.isRefreshTokenWaiting = true;
      })
      .addCase(asyncActions.refreshToken.fulfilled, (state) => {
        state.isRefreshTokenWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(asyncActions.refreshToken.rejected, (state) => {
        state.isRefreshTokenWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // obtainAuthToken
      .addCase(asyncActions.obtainAuthToken.pending, (state) => {
        state.isSaveAuthWaiting = true;
      })
      .addCase(asyncActions.obtainAuthToken.fulfilled, (state) => {
        state.isSaveAuthWaiting = false;
      })
      .addCase(asyncActions.obtainAuthToken.rejected, (state) => {
        state.isSaveAuthWaiting = false;
      })

      // saveAuth
      .addCase(asyncActions.saveAuth.pending, (state) => {
        state.isSaveAuthWaiting = true;
      })
      .addCase(asyncActions.saveAuth.fulfilled, (state) => {
        state.isSaveAuthWaiting = false;
      })
      .addCase(asyncActions.saveAuth.rejected, (state) => {
        state.isSaveAuthWaiting = false;
      })

      // getUsersMe
      .addCase(asyncActions.getUsersMe.pending, (state) => {
        state.isGetMeWaiting = true;
      })
      .addCase(asyncActions.getUsersMe.fulfilled, (state, action) => {
        state.me = action.payload.user;
        state.isGetMeWaiting = false;
      })
      .addCase(asyncActions.getUsersMe.rejected, (state) => {
        state.isGetMeWaiting = false;
      })

      // listUsers
      .addCase(asyncActions.listUsers.pending, (state) => {
        state.isListUsersWaiting = true;
      })
      .addCase(asyncActions.listUsers.fulfilled, (state, action) => {
        usersAdapter.setAll(state.users, action.payload.users);
        userInvitationsAdapter.setAll(
          state.userInvitations,
          action.payload.userInvitations,
        );
        state.isListUsersWaiting = false;
      })
      .addCase(asyncActions.listUsers.rejected, (state) => {
        state.isListUsersWaiting = false;
      })

      // requestMagicLink
      .addCase(asyncActions.requestMagicLink.pending, (state) => {
        state.isRequestMagicLinkWaiting = true;
      })
      .addCase(asyncActions.requestMagicLink.fulfilled, (state) => {
        state.isRequestMagicLinkWaiting = false;
      })
      .addCase(asyncActions.requestMagicLink.rejected, (state) => {
        state.isRequestMagicLinkWaiting = false;
      })

      // authenticateWithMagicLink
      .addCase(asyncActions.authenticateWithMagicLink.pending, (state) => {
        state.isAuthenticateWithMagicLinkWaiting = true;
      })
      .addCase(asyncActions.authenticateWithMagicLink.fulfilled, (state) => {
        state.isAuthenticateWithMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(asyncActions.authenticateWithMagicLink.rejected, (state) => {
        state.isAuthenticateWithMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // registerWithMagicLink
      .addCase(asyncActions.registerWithMagicLink.pending, (state) => {
        state.isRegisterWithMagicLinkWaiting = true;
      })
      .addCase(asyncActions.registerWithMagicLink.fulfilled, (state) => {
        state.isRegisterWithMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(asyncActions.registerWithMagicLink.rejected, (state) => {
        state.isRegisterWithMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // signout
      .addCase(asyncActions.signout.pending, (state) => {
        state.isSignoutWaiting = true;
      })
      .addCase(asyncActions.signout.fulfilled, (state) => {
        state.isSignoutWaiting = false;
      })
      .addCase(asyncActions.signout.rejected, (state) => {
        state.isSignoutWaiting = false;
      })

      // invite
      .addCase(asyncActions.invite.pending, (state) => {
        state.isInviteWaiting = true;
      })
      .addCase(asyncActions.invite.fulfilled, (state) => {
        state.isInviteWaiting = false;
      })
      .addCase(asyncActions.invite.rejected, (state) => {
        state.isInviteWaiting = false;
      })

      // oauthGoogleAuthCodeUrl
      .addCase(asyncActions.oauthGoogleAuthCodeUrl.pending, (state) => {
        state.isOauthGoogleAuthWaiting = true;
      })
      .addCase(asyncActions.oauthGoogleAuthCodeUrl.fulfilled, (state) => {
        state.isOauthGoogleAuthWaiting = false;
      })
      .addCase(asyncActions.oauthGoogleAuthCodeUrl.rejected, (state) => {
        state.isOauthGoogleAuthWaiting = false;
      })

      // getUserGroups
      .addCase(groupsStore.asyncActions.listGroups.pending, () => {})
      .addCase(
        groupsStore.asyncActions.listGroups.fulfilled,
        (state, action) => {
          usersAdapter.setAll(state.users, action.payload.users);
        },
      )
      .addCase(groupsStore.asyncActions.listGroups.rejected, () => {})

      // updateUser
      .addCase(asyncActions.updateUser.pending, (state) => {
        state.isUpdateUserWaiting = true;
      })
      .addCase(asyncActions.updateUser.fulfilled, (state, action) => {
        state.isUpdateUserWaiting = false;
        if (state.me) {
          state.me = action.payload.user;
        }
      })
      .addCase(asyncActions.updateUser.rejected, (state) => {
        state.isUpdateUserWaiting = false;
      })

      // updateUserEmail
      .addCase(asyncActions.updateUserEmail.pending, (state) => {
        state.isUpdateUserEmailWaiting = true;
      })
      .addCase(asyncActions.updateUserEmail.fulfilled, (state, action) => {
        state.isUpdateUserEmailWaiting = false;
        state.me = action.payload.user;
      })
      .addCase(asyncActions.updateUserEmail.rejected, (state) => {
        state.isUpdateUserEmailWaiting = false;
      })

      // usersSendUpdateEmailInstructions
      .addCase(
        asyncActions.usersSendUpdateEmailInstructions.pending,
        (state) => {
          state.isUsersSendUpdateEmailInstructionsWaiting = true;
        },
      )
      .addCase(
        asyncActions.usersSendUpdateEmailInstructions.fulfilled,
        (state) => {
          state.isUsersSendUpdateEmailInstructionsWaiting = false;
        },
      )
      .addCase(
        asyncActions.usersSendUpdateEmailInstructions.rejected,
        (state) => {
          state.isUsersSendUpdateEmailInstructionsWaiting = false;
        },
      )

      // listPages
      .addCase(pagesStore.asyncActions.listPages.pending, () => {})
      .addCase(pagesStore.asyncActions.listPages.fulfilled, (state, action) => {
        usersAdapter.setAll(state.users, action.payload.users);
      })
      .addCase(pagesStore.asyncActions.listPages.rejected, () => {})

      // updateOrganizationUser
      .addCase(
        organizationsStore.asyncActions.updateOrganizationUser.pending,
        () => {},
      )
      .addCase(
        organizationsStore.asyncActions.updateOrganizationUser.fulfilled,
        (state, action) => {
          usersAdapter.updateOne(state.users, {
            id: action.payload.id,
            changes: action.payload,
          });
        },
      )
      .addCase(
        organizationsStore.asyncActions.updateOrganizationUser.rejected,
        () => {},
      )

      // invitationsResend
      .addCase(asyncActions.invitationsResend.pending, (state) => {
        state.isInvitationsResendWaiting = true;
      })
      .addCase(asyncActions.invitationsResend.fulfilled, (state) => {
        state.isInvitationsResendWaiting = false;
      })
      .addCase(asyncActions.invitationsResend.rejected, (state) => {
        state.isInvitationsResendWaiting = false;
      })

      // requestInvitationMagicLink
      .addCase(asyncActions.requestInvitationMagicLink.pending, (state) => {
        state.isRequestInvitationMagicLinkWaiting = true;
      })
      .addCase(asyncActions.requestInvitationMagicLink.fulfilled, (state) => {
        state.isRequestInvitationMagicLinkWaiting = false;
      })
      .addCase(asyncActions.requestInvitationMagicLink.rejected, (state) => {
        state.isRequestInvitationMagicLinkWaiting = false;
      })

      // authenticateWithInvitationMagicLink
      .addCase(asyncActions.authenticateWithInvitationMagicLink.pending, (state) => {
        state.isAuthenticateWithInvitationMagicLinkWaiting = true;
      })
      .addCase(asyncActions.authenticateWithInvitationMagicLink.fulfilled, (state) => {
        state.isAuthenticateWithInvitationMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(asyncActions.authenticateWithInvitationMagicLink.rejected, (state) => {
        state.isAuthenticateWithInvitationMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      })

      // registerWithInvitationMagicLink
      .addCase(asyncActions.registerWithInvitationMagicLink.pending, (state) => {
        state.isRegisterWithInvitationMagicLinkWaiting = true;
      })
      .addCase(asyncActions.registerWithInvitationMagicLink.fulfilled, (state) => {
        state.isRegisterWithInvitationMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = true;
        state.isAuthFailed = false;
      })
      .addCase(asyncActions.registerWithInvitationMagicLink.rejected, (state) => {
        state.isRegisterWithInvitationMagicLinkWaiting = false;
        state.isAuthChecked = true;
        state.isAuthSucceeded = false;
        state.isAuthFailed = true;
      });
  },
  initialState,
  name: 'users',
  reducers: {},
});
