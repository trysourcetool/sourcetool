// @generated by protoc-gen-es v2.0.0 with parameter "target=ts,import_extension=js,json_types=true"
// @generated from file page/v1/page.proto (package page.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file page/v1/page.proto.
 */
export const file_page_v1_page: GenFile = /*@__PURE__*/
  fileDesc("ChJwYWdlL3YxL3BhZ2UucHJvdG8SB3BhZ2UudjEiTQoEUGFnZRIKCgJpZBgBIAEoCRIMCgRuYW1lGAIgASgJEg0KBXJvdXRlGAMgASgJEgwKBHBhdGgYBCADKAUSDgoGZ3JvdXBzGAUgAygJQlUKC2NvbS5wYWdlLnYxQglQYWdlUHJvdG9QAaICA1BYWKoCB1BhZ2UuVjHKAgdQYWdlXFYx4gITUGFnZVxWMVxHUEJNZXRhZGF0YeoCCFBhZ2U6OlYxYgZwcm90bzM");

/**
 * @generated from message page.v1.Page
 */
export type Page = Message<"page.v1.Page"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;

  /**
   * @generated from field: string route = 3;
   */
  route: string;

  /**
   * @generated from field: repeated int32 path = 4;
   */
  path: number[];

  /**
   * @generated from field: repeated string groups = 5;
   */
  groups: string[];
};

/**
 * JSON type for the message page.v1.Page.
 */
export type PageJson = {
  /**
   * @generated from field: string id = 1;
   */
  id?: string;

  /**
   * @generated from field: string name = 2;
   */
  name?: string;

  /**
   * @generated from field: string route = 3;
   */
  route?: string;

  /**
   * @generated from field: repeated int32 path = 4;
   */
  path?: number[];

  /**
   * @generated from field: repeated string groups = 5;
   */
  groups?: string[];
};

/**
 * Describes the message page.v1.Page.
 * Use `create(PageSchema)` to create a new message.
 */
export const PageSchema: GenMessage<Page, PageJson> = /*@__PURE__*/
  messageDesc(file_page_v1_page, 0);

