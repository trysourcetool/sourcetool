import { expect, test } from 'vitest';
import { Cursor } from './uibuilder';

test('coursor path management', () => {
  const cursor = new Cursor();

  // Test initial path
  const initialPath = cursor.getPath();
  expect(initialPath.length).toEqual(1);

  // Test next()
  cursor.next();
  const nextPath = cursor.getPath();
  console.log({ nextPath });
  expect(nextPath.length).toEqual(1);
  expect(nextPath[0]).toEqual(1);

  // Add parent path
  cursor.parentPath.push(1);
  const parentPath = cursor.getPath();
  console.log({ parentPath });
  expect(parentPath.length).toEqual(2);
  expect(parentPath[0]).toEqual(1);
  expect(parentPath[1]).toEqual(1);
});
