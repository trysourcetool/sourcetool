import { ENVIRONMENTS } from '@/environments';
import { useAuth } from '@/hooks/use-auth';
import { useDispatch, useSelector } from '@/store';
import { widgetsStore } from '@/store/modules/widgets';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useParams, useLocation, useNavigate } from 'react-router';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { v4 as uuidv4 } from 'uuid';
import {
  MessageSchema,
  type CloseSessionJson,
  type InitializeClientJson,
  type RerunPageJson,
} from '@/pb/ts/websocket/v1/message_pb';
import {
  create,
  fromBinary,
  fromJson,
  toBinary,
  toJson,
} from '@bufbuild/protobuf';
import {
  WidgetSchema,
  type Widget,
} from '@/pb/ts/widget/v1/widget_pb';
import { pagesStore } from '@/store/modules/pages';
import type { WidgetType } from '@/store/modules/widgets';
import { hostInstancesStore } from '@/store/modules/hostInstances';
import { $path } from 'safe-routes';

const WebSocketBlock = ({ onDisable }: { onDisable: () => void }) => {
  const dispatch = useDispatch();
  const { '*': path } = useParams();
  const currentPageId = useRef('');
  const currentSessionId = useRef('');
  const prevVisibilityStatus = useRef(!document.hidden);
  const navigate = useNavigate();
  const [isVisibilityStatus, setIsVisibilityStatus] = useState(
    !document.hidden,
  );
  const isInitialLoading = useRef(false);
  const inactiveTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const INACTIVITY_TIMEOUT_MS = 10 * 60 * 1000; // 10 minutes

  const socketUrl = useMemo(
    () =>
      `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/ws`,
    [],
  );

  const pageId = useSelector(
    (state) =>
      pagesStore.selector.getPageFromPath(state, path ? `/${path}` : '')?.id,
  );
  const widgetEntities = useSelector((state) =>
    widgetsStore.selector.getWidgetEntities(state),
  );
  const widgetUpdateAt = useSelector((state) => state.widgets.updateAt);
  const exception = useSelector((state) => state.pages.exception);
  const isHostInstancePingError = useSelector(
    (state) => state.hostInstances.isHostInstancePingError,
  );

  const { sendMessage, readyState } = useWebSocket<Uint8Array>(socketUrl, {
    onMessage: (event) => {
      event.data.arrayBuffer().then((arrayBuffer: ArrayBuffer) => {
        const message = toJson(
          MessageSchema,
          fromBinary(MessageSchema, new Uint8Array(arrayBuffer)),
        );
        console.table({ message });
        if (message.initializeClientCompleted) {
          isInitialLoading.current = false;
          if (message.initializeClientCompleted.sessionId) {
            currentSessionId.current =
              message.initializeClientCompleted.sessionId;
          }
        }
        if (message.renderWidget) {
          dispatch(widgetsStore.actions.setWidgetData(message.renderWidget));
        }
        if (message.scriptFinished) {
          console.log('scriptFinished', message.scriptFinished);
          dispatch(widgetsStore.actions.renderWidgetCompleted());
        }
        if (message.exception) {
          dispatch(pagesStore.actions.setException(message.exception));
        }
      });
    },
    shouldReconnect: (event) => {
      console.log('shouldReconnect', event);
      if (isHostInstancePingError) {
        return false;
      }
      return true;
    },
  });

  const connectionStatus = {
    [ReadyState.CONNECTING]: ReadyState.CONNECTING,
    [ReadyState.OPEN]: ReadyState.OPEN,
    [ReadyState.CLOSING]: ReadyState.CLOSING,
    [ReadyState.CLOSED]: ReadyState.CLOSED,
    [ReadyState.UNINSTANTIATED]: ReadyState.UNINSTANTIATED,
  }[readyState];

  console.log({
    connectionStatus,
    currentSessionId: currentSessionId.current,
    isVisibilityStatus,
  });

  const handleCloseSession = useCallback(() => {
    if (inactiveTimeoutRef.current) {
      clearTimeout(inactiveTimeoutRef.current);
      inactiveTimeoutRef.current = null;
    }
    console.log('handleCloseSession');
    sendMessage(
      toBinary(
        MessageSchema,
        create(MessageSchema, {
          id: uuidv4(),
          type: {
            case: 'closeSession',
            value: {
              sessionId: currentSessionId.current,
            } satisfies CloseSessionJson,
          },
        }),
      ),
    );
    currentSessionId.current = '';
  }, [sendMessage]);

  useEffect(() => {
    const handleVisibilityChange = () => {
      prevVisibilityStatus.current = isVisibilityStatus;
      setIsVisibilityStatus(!document.hidden);
    };
    window.addEventListener('beforeunload', handleCloseSession);
    window.addEventListener('visibilitychange', handleVisibilityChange);
    return () => {
      window.removeEventListener('beforeunload', handleCloseSession);
      window.removeEventListener('visibilitychange', handleVisibilityChange);
      if (inactiveTimeoutRef.current) {
        clearTimeout(inactiveTimeoutRef.current);
      }
    };
  }, [isVisibilityStatus, handleCloseSession]);

  useEffect(() => {
    (async () => {
      if (!pageId && !currentPageId.current) {
        return;
      }
      if (isInitialLoading.current && currentPageId.current === pageId) {
        return;
      }
      isInitialLoading.current = true;
      if (pageId) {
        const resultAction = await dispatch(
          hostInstancesStore.asyncActions.getHostInstancePing({
            pageId,
          }),
        );
        if (
          hostInstancesStore.asyncActions.getHostInstancePing.rejected.match(
            resultAction,
          )
        ) {
          currentSessionId.current = '';
          navigate($path('/error/hostInstancePingError'));
          return;
        }
        if (pageId && currentPageId.current === '') {
          sendMessage(
            toBinary(
              MessageSchema,
              create(MessageSchema, {
                id: uuidv4(),
                type: {
                  case: 'initializeClient',
                  value: {
                    pageId: pageId,
                  } satisfies InitializeClientJson,
                },
              }),
            ),
          );

          currentPageId.current = pageId;
        }

        if (
          pageId &&
          currentPageId.current &&
          currentPageId.current !== pageId
        ) {
          dispatch(widgetsStore.actions.clearWidgets());
          dispatch(pagesStore.actions.clearException());
          sendMessage(
            toBinary(
              MessageSchema,
              create(MessageSchema, {
                id: uuidv4(),
                type: {
                  case: 'rerunPage',
                  value: {
                    sessionId: currentSessionId.current,
                    pageId: pageId,
                    states: [],
                  } satisfies RerunPageJson,
                },
              }),
            ),
          );
          currentPageId.current = pageId;
        }
      } else if (currentPageId.current) {
        console.log('CLEAR WIDGETS AND DISABLE CONTROLLER');
        dispatch(widgetsStore.actions.clearWidgets());
        dispatch(pagesStore.actions.clearException());
        handleCloseSession();
        currentPageId.current = '';
        setTimeout(() => {
          onDisable();
        }, 1000);
      }
      isInitialLoading.current = false;
    })();
  }, [pageId]);

  const handleRerunPage = useCallback(() => {
    const states: Widget[] = [];
    Object.values(widgetEntities).forEach((widget) => {
      if (widget.widget) {
        const widgetData: any = {
          id: widget.widget.id,
        };

        Object.keys(widget.widget).forEach((key) => {
          if (widget.widget && key !== 'id') {
            widgetData[key] = {};

            if (widget.widget[key as WidgetType]) {
              Object.keys((widget.widget as any)[key]).forEach((subKey) => {
                if ((widget?.widget as any)?.[key]?.[subKey] !== undefined) {
                  widgetData[key][subKey] = (widget.widget as any)[key][subKey];
                }
              });
            }
          }
        });

        states.push(fromJson(WidgetSchema, widgetData));
      }
    });

    console.log({ states });

    sendMessage(
      toBinary(
        MessageSchema,
        create(MessageSchema, {
          id: uuidv4(),
          type: {
            case: 'rerunPage',
            value: {
              pageId: currentPageId.current,
              sessionId: currentSessionId.current,
              states: states,
            },
          },
        }),
      ),
    );
  }, [widgetUpdateAt]);

  useEffect(() => {
    if (widgetUpdateAt && currentPageId.current) {
      handleRerunPage();
    }
  }, [widgetUpdateAt]);

  useEffect(() => {
    (async () => {
      if (!pageId || !currentPageId.current) {
        return;
      }
      if (isVisibilityStatus && !exception) {
        console.log('Tab became active');
        if (inactiveTimeoutRef.current) {
          console.log('Clearing inactivity timer');
          clearTimeout(inactiveTimeoutRef.current);
          inactiveTimeoutRef.current = null;
        }

        if (!currentSessionId.current) {
          console.log('Re-initializing client after becoming active');
          const resultAction = await dispatch(
            hostInstancesStore.asyncActions.getHostInstancePing({
              pageId,
            }),
          );
          if (
            hostInstancesStore.asyncActions.getHostInstancePing.rejected.match(
              resultAction,
            )
          ) {
            currentSessionId.current = '';
            navigate($path('/error/hostInstancePingError'));
            return;
          }

          sendMessage(
            toBinary(
              MessageSchema,
              create(MessageSchema, {
                id: uuidv4(),
                type: {
                  case: 'initializeClient',
                  value: {
                    pageId: pageId,
                    sessionId: currentSessionId.current,
                  } satisfies InitializeClientJson,
                },
              }),
            ),
          );
        } else {
          console.log('Session still active, no need to re-initialize');
        }
      } else {
        console.log('Tab became inactive');
        if (!inactiveTimeoutRef.current && currentSessionId.current) {
          console.log(`Starting inactivity timer (${INACTIVITY_TIMEOUT_MS}ms)`);
          inactiveTimeoutRef.current = setTimeout(() => {
            console.log('Inactivity timer expired, closing session');
            handleCloseSession();
            inactiveTimeoutRef.current = null;
          }, INACTIVITY_TIMEOUT_MS);
        }
      }
    })();
  }, [isVisibilityStatus, pageId, exception]);

  return <></>;
};

export const WebSocketController = () => {
  const location = useLocation();
  const [isSocketReady, setIsSocketReady] = useState(false);
  const { isAuthChecked, isSubDomainMatched } = useAuth();

  useEffect(() => {
    if (
      isAuthChecked === 'checked' &&
      ((ENVIRONMENTS.IS_CLOUD_EDITION && isSubDomainMatched) || !ENVIRONMENTS.IS_CLOUD_EDITION) &&
      location.pathname.match(/^\/pages\/.*$/)
    ) {
      setIsSocketReady(true);
    }
  }, [
    location.pathname,
    isSubDomainMatched,
    isAuthChecked,
  ]);

  return (
    <>
      {isSocketReady && (
        <WebSocketBlock onDisable={() => setIsSocketReady(false)} />
      )}
    </>
  );
};
