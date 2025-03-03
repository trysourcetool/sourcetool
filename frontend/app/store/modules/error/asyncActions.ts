/* eslint-disable no-useless-return */

import { createAsyncThunk } from '@reduxjs/toolkit';

export const handleError = createAsyncThunk(
  'error/handleError',
  async (data: any, { dispatch }) => {},
);
