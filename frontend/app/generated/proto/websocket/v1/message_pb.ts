// @generated by protoc-gen-es v2.0.0 with parameter "target=ts,import_extension=js,json_types=true"
// @generated from file websocket/v1/message.proto (package websocket.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { Exception, ExceptionJson } from "../../exception/v1/exception_pb.js";
import { file_exception_v1_exception } from "../../exception/v1/exception_pb.js";
import type { Page, PageJson } from "../../page/v1/page_pb.js";
import { file_page_v1_page } from "../../page/v1/page_pb.js";
import type { Widget, WidgetJson } from "../../widget/v1/widget_pb.js";
import { file_widget_v1_widget } from "../../widget/v1/widget_pb.js";
import type { Message as Message$1 } from "@bufbuild/protobuf";

/**
 * Describes the file websocket/v1/message.proto.
 */
export const file_websocket_v1_message: GenFile = /*@__PURE__*/
  fileDesc("Chp3ZWJzb2NrZXQvdjEvbWVzc2FnZS5wcm90bxIMd2Vic29ja2V0LnYxIq8ECgdNZXNzYWdlEgoKAmlkGAEgASgJEiwKCWV4Y2VwdGlvbhgCIAEoCzIXLmV4Y2VwdGlvbi52MS5FeGNlcHRpb25IABI3Cg9pbml0aWFsaXplX2hvc3QYAyABKAsyHC53ZWJzb2NrZXQudjEuSW5pdGlhbGl6ZUhvc3RIABJKChlpbml0aWFsaXplX2hvc3RfY29tcGxldGVkGAQgASgLMiUud2Vic29ja2V0LnYxLkluaXRpYWxpemVIb3N0Q29tcGxldGVkSAASOwoRaW5pdGlhbGl6ZV9jbGllbnQYBSABKAsyHi53ZWJzb2NrZXQudjEuSW5pdGlhbGl6ZUNsaWVudEgAEk4KG2luaXRpYWxpemVfY2xpZW50X2NvbXBsZXRlZBgGIAEoCzInLndlYnNvY2tldC52MS5Jbml0aWFsaXplQ2xpZW50Q29tcGxldGVkSAASMwoNcmVuZGVyX3dpZGdldBgHIAEoCzIaLndlYnNvY2tldC52MS5SZW5kZXJXaWRnZXRIABItCgpyZXJ1bl9wYWdlGAggASgLMhcud2Vic29ja2V0LnYxLlJlcnVuUGFnZUgAEjMKDWNsb3NlX3Nlc3Npb24YCSABKAsyGi53ZWJzb2NrZXQudjEuQ2xvc2VTZXNzaW9uSAASNwoPc2NyaXB0X2ZpbmlzaGVkGAogASgLMhwud2Vic29ja2V0LnYxLlNjcmlwdEZpbmlzaGVkSABCBgoEdHlwZSJmCg5Jbml0aWFsaXplSG9zdBIPCgdhcGlfa2V5GAEgASgJEhAKCHNka19uYW1lGAIgASgJEhMKC3Nka192ZXJzaW9uGAMgASgJEhwKBXBhZ2VzGAQgAygLMg0ucGFnZS52MS5QYWdlIjMKF0luaXRpYWxpemVIb3N0Q29tcGxldGVkEhgKEGhvc3RfaW5zdGFuY2VfaWQYASABKAkiSwoQSW5pdGlhbGl6ZUNsaWVudBIXCgpzZXNzaW9uX2lkGAEgASgJSACIAQESDwoHcGFnZV9pZBgCIAEoCUINCgtfc2Vzc2lvbl9pZCIvChlJbml0aWFsaXplQ2xpZW50Q29tcGxldGVkEhIKCnNlc3Npb25faWQYASABKAkiZAoMUmVuZGVyV2lkZ2V0EhIKCnNlc3Npb25faWQYASABKAkSDwoHcGFnZV9pZBgCIAEoCRIMCgRwYXRoGAMgAygFEiEKBndpZGdldBgEIAEoCzIRLndpZGdldC52MS5XaWRnZXQiUwoJUmVydW5QYWdlEhIKCnNlc3Npb25faWQYASABKAkSDwoHcGFnZV9pZBgCIAEoCRIhCgZzdGF0ZXMYAyADKAsyES53aWRnZXQudjEuV2lkZ2V0IiIKDENsb3NlU2Vzc2lvbhISCgpzZXNzaW9uX2lkGAEgASgJIqMBCg5TY3JpcHRGaW5pc2hlZBISCgpzZXNzaW9uX2lkGAEgASgJEjMKBnN0YXR1cxgCIAEoDjIjLndlYnNvY2tldC52MS5TY3JpcHRGaW5pc2hlZC5TdGF0dXMiSAoGU3RhdHVzEhYKElNUQVRVU19VTlNQRUNJRklFRBAAEhIKDlNUQVRVU19TVUNDRVNTEAESEgoOU1RBVFVTX0ZBSUxVUkUQAkJxChBjb20ud2Vic29ja2V0LnYxQgxNZXNzYWdlUHJvdG9QAaICA1dYWKoCDFdlYnNvY2tldC5WMcoCDFdlYnNvY2tldFxWMeICGFdlYnNvY2tldFxWMVxHUEJNZXRhZGF0YeoCDVdlYnNvY2tldDo6VjFiBnByb3RvMw", [file_exception_v1_exception, file_page_v1_page, file_widget_v1_widget]);

/**
 * @generated from message websocket.v1.Message
 */
export type Message = Message$1<"websocket.v1.Message"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from oneof websocket.v1.Message.type
   */
  type: {
    /**
     * @generated from field: exception.v1.Exception exception = 2;
     */
    value: Exception;
    case: "exception";
  } | {
    /**
     * @generated from field: websocket.v1.InitializeHost initialize_host = 3;
     */
    value: InitializeHost;
    case: "initializeHost";
  } | {
    /**
     * @generated from field: websocket.v1.InitializeHostCompleted initialize_host_completed = 4;
     */
    value: InitializeHostCompleted;
    case: "initializeHostCompleted";
  } | {
    /**
     * @generated from field: websocket.v1.InitializeClient initialize_client = 5;
     */
    value: InitializeClient;
    case: "initializeClient";
  } | {
    /**
     * @generated from field: websocket.v1.InitializeClientCompleted initialize_client_completed = 6;
     */
    value: InitializeClientCompleted;
    case: "initializeClientCompleted";
  } | {
    /**
     * @generated from field: websocket.v1.RenderWidget render_widget = 7;
     */
    value: RenderWidget;
    case: "renderWidget";
  } | {
    /**
     * @generated from field: websocket.v1.RerunPage rerun_page = 8;
     */
    value: RerunPage;
    case: "rerunPage";
  } | {
    /**
     * @generated from field: websocket.v1.CloseSession close_session = 9;
     */
    value: CloseSession;
    case: "closeSession";
  } | {
    /**
     * @generated from field: websocket.v1.ScriptFinished script_finished = 10;
     */
    value: ScriptFinished;
    case: "scriptFinished";
  } | { case: undefined; value?: undefined };
};

/**
 * JSON type for the message websocket.v1.Message.
 */
export type MessageJson = {
  /**
   * @generated from field: string id = 1;
   */
  id?: string;

  /**
   * @generated from field: exception.v1.Exception exception = 2;
   */
  exception?: ExceptionJson;

  /**
   * @generated from field: websocket.v1.InitializeHost initialize_host = 3;
   */
  initializeHost?: InitializeHostJson;

  /**
   * @generated from field: websocket.v1.InitializeHostCompleted initialize_host_completed = 4;
   */
  initializeHostCompleted?: InitializeHostCompletedJson;

  /**
   * @generated from field: websocket.v1.InitializeClient initialize_client = 5;
   */
  initializeClient?: InitializeClientJson;

  /**
   * @generated from field: websocket.v1.InitializeClientCompleted initialize_client_completed = 6;
   */
  initializeClientCompleted?: InitializeClientCompletedJson;

  /**
   * @generated from field: websocket.v1.RenderWidget render_widget = 7;
   */
  renderWidget?: RenderWidgetJson;

  /**
   * @generated from field: websocket.v1.RerunPage rerun_page = 8;
   */
  rerunPage?: RerunPageJson;

  /**
   * @generated from field: websocket.v1.CloseSession close_session = 9;
   */
  closeSession?: CloseSessionJson;

  /**
   * @generated from field: websocket.v1.ScriptFinished script_finished = 10;
   */
  scriptFinished?: ScriptFinishedJson;
};

/**
 * Describes the message websocket.v1.Message.
 * Use `create(MessageSchema)` to create a new message.
 */
export const MessageSchema: GenMessage<Message, MessageJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 0);

/**
 * @generated from message websocket.v1.InitializeHost
 */
export type InitializeHost = Message$1<"websocket.v1.InitializeHost"> & {
  /**
   * @generated from field: string api_key = 1;
   */
  apiKey: string;

  /**
   * @generated from field: string sdk_name = 2;
   */
  sdkName: string;

  /**
   * @generated from field: string sdk_version = 3;
   */
  sdkVersion: string;

  /**
   * @generated from field: repeated page.v1.Page pages = 4;
   */
  pages: Page[];
};

/**
 * JSON type for the message websocket.v1.InitializeHost.
 */
export type InitializeHostJson = {
  /**
   * @generated from field: string api_key = 1;
   */
  apiKey?: string;

  /**
   * @generated from field: string sdk_name = 2;
   */
  sdkName?: string;

  /**
   * @generated from field: string sdk_version = 3;
   */
  sdkVersion?: string;

  /**
   * @generated from field: repeated page.v1.Page pages = 4;
   */
  pages?: PageJson[];
};

/**
 * Describes the message websocket.v1.InitializeHost.
 * Use `create(InitializeHostSchema)` to create a new message.
 */
export const InitializeHostSchema: GenMessage<InitializeHost, InitializeHostJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 1);

/**
 * @generated from message websocket.v1.InitializeHostCompleted
 */
export type InitializeHostCompleted = Message$1<"websocket.v1.InitializeHostCompleted"> & {
  /**
   * @generated from field: string host_instance_id = 1;
   */
  hostInstanceId: string;
};

/**
 * JSON type for the message websocket.v1.InitializeHostCompleted.
 */
export type InitializeHostCompletedJson = {
  /**
   * @generated from field: string host_instance_id = 1;
   */
  hostInstanceId?: string;
};

/**
 * Describes the message websocket.v1.InitializeHostCompleted.
 * Use `create(InitializeHostCompletedSchema)` to create a new message.
 */
export const InitializeHostCompletedSchema: GenMessage<InitializeHostCompleted, InitializeHostCompletedJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 2);

/**
 * @generated from message websocket.v1.InitializeClient
 */
export type InitializeClient = Message$1<"websocket.v1.InitializeClient"> & {
  /**
   * @generated from field: optional string session_id = 1;
   */
  sessionId?: string;

  /**
   * @generated from field: string page_id = 2;
   */
  pageId: string;
};

/**
 * JSON type for the message websocket.v1.InitializeClient.
 */
export type InitializeClientJson = {
  /**
   * @generated from field: optional string session_id = 1;
   */
  sessionId?: string;

  /**
   * @generated from field: string page_id = 2;
   */
  pageId?: string;
};

/**
 * Describes the message websocket.v1.InitializeClient.
 * Use `create(InitializeClientSchema)` to create a new message.
 */
export const InitializeClientSchema: GenMessage<InitializeClient, InitializeClientJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 3);

/**
 * @generated from message websocket.v1.InitializeClientCompleted
 */
export type InitializeClientCompleted = Message$1<"websocket.v1.InitializeClientCompleted"> & {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId: string;
};

/**
 * JSON type for the message websocket.v1.InitializeClientCompleted.
 */
export type InitializeClientCompletedJson = {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId?: string;
};

/**
 * Describes the message websocket.v1.InitializeClientCompleted.
 * Use `create(InitializeClientCompletedSchema)` to create a new message.
 */
export const InitializeClientCompletedSchema: GenMessage<InitializeClientCompleted, InitializeClientCompletedJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 4);

/**
 * @generated from message websocket.v1.RenderWidget
 */
export type RenderWidget = Message$1<"websocket.v1.RenderWidget"> & {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId: string;

  /**
   * @generated from field: string page_id = 2;
   */
  pageId: string;

  /**
   * @generated from field: repeated int32 path = 3;
   */
  path: number[];

  /**
   * @generated from field: widget.v1.Widget widget = 4;
   */
  widget?: Widget;
};

/**
 * JSON type for the message websocket.v1.RenderWidget.
 */
export type RenderWidgetJson = {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId?: string;

  /**
   * @generated from field: string page_id = 2;
   */
  pageId?: string;

  /**
   * @generated from field: repeated int32 path = 3;
   */
  path?: number[];

  /**
   * @generated from field: widget.v1.Widget widget = 4;
   */
  widget?: WidgetJson;
};

/**
 * Describes the message websocket.v1.RenderWidget.
 * Use `create(RenderWidgetSchema)` to create a new message.
 */
export const RenderWidgetSchema: GenMessage<RenderWidget, RenderWidgetJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 5);

/**
 * @generated from message websocket.v1.RerunPage
 */
export type RerunPage = Message$1<"websocket.v1.RerunPage"> & {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId: string;

  /**
   * @generated from field: string page_id = 2;
   */
  pageId: string;

  /**
   * @generated from field: repeated widget.v1.Widget states = 3;
   */
  states: Widget[];
};

/**
 * JSON type for the message websocket.v1.RerunPage.
 */
export type RerunPageJson = {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId?: string;

  /**
   * @generated from field: string page_id = 2;
   */
  pageId?: string;

  /**
   * @generated from field: repeated widget.v1.Widget states = 3;
   */
  states?: WidgetJson[];
};

/**
 * Describes the message websocket.v1.RerunPage.
 * Use `create(RerunPageSchema)` to create a new message.
 */
export const RerunPageSchema: GenMessage<RerunPage, RerunPageJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 6);

/**
 * @generated from message websocket.v1.CloseSession
 */
export type CloseSession = Message$1<"websocket.v1.CloseSession"> & {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId: string;
};

/**
 * JSON type for the message websocket.v1.CloseSession.
 */
export type CloseSessionJson = {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId?: string;
};

/**
 * Describes the message websocket.v1.CloseSession.
 * Use `create(CloseSessionSchema)` to create a new message.
 */
export const CloseSessionSchema: GenMessage<CloseSession, CloseSessionJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 7);

/**
 * @generated from message websocket.v1.ScriptFinished
 */
export type ScriptFinished = Message$1<"websocket.v1.ScriptFinished"> & {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId: string;

  /**
   * @generated from field: websocket.v1.ScriptFinished.Status status = 2;
   */
  status: ScriptFinished_Status;
};

/**
 * JSON type for the message websocket.v1.ScriptFinished.
 */
export type ScriptFinishedJson = {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId?: string;

  /**
   * @generated from field: websocket.v1.ScriptFinished.Status status = 2;
   */
  status?: ScriptFinished_StatusJson;
};

/**
 * Describes the message websocket.v1.ScriptFinished.
 * Use `create(ScriptFinishedSchema)` to create a new message.
 */
export const ScriptFinishedSchema: GenMessage<ScriptFinished, ScriptFinishedJson> = /*@__PURE__*/
  messageDesc(file_websocket_v1_message, 8);

/**
 * @generated from enum websocket.v1.ScriptFinished.Status
 */
export enum ScriptFinished_Status {
  /**
   * @generated from enum value: STATUS_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: STATUS_SUCCESS = 1;
   */
  SUCCESS = 1,

  /**
   * @generated from enum value: STATUS_FAILURE = 2;
   */
  FAILURE = 2,
}

/**
 * JSON type for the enum websocket.v1.ScriptFinished.Status.
 */
export type ScriptFinished_StatusJson = "STATUS_UNSPECIFIED" | "STATUS_SUCCESS" | "STATUS_FAILURE";

/**
 * Describes the enum websocket.v1.ScriptFinished.Status.
 */
export const ScriptFinished_StatusSchema: GenEnum<ScriptFinished_Status, ScriptFinished_StatusJson> = /*@__PURE__*/
  enumDesc(file_websocket_v1_message, 8, 0);

