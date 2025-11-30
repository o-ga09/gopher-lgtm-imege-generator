import { apiClient } from "./client";
import type {
  RunAgentRequest,
  RunAgentResponse,
  CreateSessionResponse,
  SessionInfo,
  ListAppsResponse,
  ListImagesResponse,
} from "./types";

/**
 * AIエージェント一覧を取得
 */
export const listApps = async (): Promise<ListAppsResponse> => {
  const response = await apiClient.get<ListAppsResponse>("/list-apps");
  return response.data;
};

/**
 * AIエージェントを実行
 */
export const runAgent = async (
  request: RunAgentRequest
): Promise<RunAgentResponse> => {
  const response = await apiClient.post<RunAgentResponse>("/run", request);
  return response.data;
};

/**
 * AIエージェントをSSEで実行
 * @param request - エージェント実行リクエスト
 * @param onMessage - メッセージ受信時のコールバック
 * @param onError - エラー時のコールバック
 * @param onComplete - 完了時のコールバック
 */
export const runAgentSSE = (
  request: RunAgentRequest,
  onMessage: (data: string) => void,
  onError: (error: Error) => void,
  _onComplete: () => void
): (() => void) => {
  const baseURL =
    import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/v1/agent";
  const url = `${baseURL}/run_sse`;

  const eventSource = new EventSource(url);

  eventSource.onmessage = (event) => {
    onMessage(event.data);
  };

  eventSource.onerror = () => {
    eventSource.close();
    onError(new Error("SSE connection error"));
  };

  // SSE接続を開始する前にリクエストを送信
  fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  }).catch((error) => {
    eventSource.close();
    onError(error);
  });

  // クリーンアップ関数を返す
  return () => {
    eventSource.close();
  };
};

/**
 * セッションを作成
 */
export const createSession = async (
  appName: string,
  userId: string
): Promise<CreateSessionResponse> => {
  const response = await apiClient.post<CreateSessionResponse>(
    `/apps/${appName}/users/${userId}/sessions`
  );
  return response.data;
};

/**
 * 任意のIDでセッションを作成
 */
export const createSessionWithId = async (
  appName: string,
  userId: string,
  sessionId: string
): Promise<CreateSessionResponse> => {
  const response = await apiClient.post<CreateSessionResponse>(
    `/apps/${appName}/users/${userId}/sessions/${sessionId}`
  );
  return response.data;
};

/**
 * セッション情報を取得
 */
export const getSession = async (
  appName: string,
  userId: string,
  sessionId: string
): Promise<SessionInfo> => {
  const response = await apiClient.get<SessionInfo>(
    `/apps/${appName}/users/${userId}/sessions/${sessionId}`
  );
  return response.data;
};

/**
 * 生成済み画像一覧を取得
 */
export const listImages = async (): Promise<ListImagesResponse> => {
  const root = import.meta.env.VITE_API_BASE_ROOT || "http://localhost:8080";
  const response = await fetch(`${root}/v1/images`);
  if (!response.ok) {
    throw new Error("Failed to fetch images history");
  }
  return (await response.json()) as ListImagesResponse;
};
