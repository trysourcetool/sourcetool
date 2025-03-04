import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../../';

export const getPageIds = createSelector(
  (state: RootState) => state.pages,
  (values) => values.pages.ids,
);

export const getPageEntities = createSelector(
  (state: RootState) => state.pages,
  (values) => values.pages.entities,
);

export const getPages = createSelector(
  (state: RootState) => state.pages,
  (values) => values.pages.ids.map((id) => values.pages.entities[id]),
);

export const getPage = createSelector(
  (state: RootState, pageId: string) => state.pages.pages.entities[pageId],
  (values) => values || null,
);

export const getPageFromPath = createSelector(
  (state: RootState, path: string) => ({
    pages: getPages(state),
    path,
  }),
  ({ pages, path }) => {
    const page = pages.find((page) => page.route === path);
    console.log({ page });
    return page || null;
  },
);

export const getPermissionPages = createSelector(
  (state: RootState) => ({
    account: state.users.me,
    pages: state.pages,
    groups: state.groups,
  }),
  ({ account, pages, groups }) => {
    const userGroups = groups.userGroups.ids
      .map((id) => groups.userGroups.entities[id])
      .filter((userGroup) => userGroup.userId === account?.id);

    const groupPages = groups.groupPages.ids
      .map((id) => groups.groupPages.entities[id])
      .filter((groupPage) =>
        userGroups.some((userGroup) => userGroup.groupId === groupPage.groupId),
      );

    console.log(
      { groupPages, userGroups },
      pages.pages,
      pages.pages.ids
        .filter((id) => groupPages.some((page) => page.pageId === id))
        .map((id) => pages.pages.entities[id]),
    );

    return pages.pages.ids
      .filter(
        (id) =>
          !groups.groups.ids.length ||
          groupPages.some((page) => page.pageId === id),
      )
      .map((id) => pages.pages.entities[id]);
  },
);
