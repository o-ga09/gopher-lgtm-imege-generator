import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { UseMutationResult, UseQueryResult } from "@tanstack/react-query";
import {
  listApps,
  runAgent,
  createSession,
  createSessionWithId,
  getSession,
  listImages,
} from "./agent";
import type {
  RunAgentRequest,
  RunAgentResponse,
  CreateSessionResponse,
  SessionInfo,
  ListAppsResponse,
} from "./types";

// クエリキー
export const queryKeys = {
  apps: ["apps"] as const,
  session: (appName: string, userId: string, sessionId: string) =>
    ["sessions", appName, userId, sessionId] as const,
  images: ["images"] as const,
};

/**
 * AIエージェント一覧を取得するフック
 */
export const useListApps = (): UseQueryResult<ListAppsResponse, Error> => {
  return useQuery({
    queryKey: queryKeys.apps,
    queryFn: listApps,
    staleTime: 5 * 60 * 1000, // 5分
  });
};

/**
 * AIエージェントを実行するフック
 */
export const useRunAgent = (): UseMutationResult<
  RunAgentResponse,
  Error,
  RunAgentRequest,
  unknown
> => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: runAgent,
    onSuccess: (_data, variables) => {
      // セッション情報を無効化して再取得
      queryClient.invalidateQueries({
        queryKey: queryKeys.session(
          variables.appName,
          variables.userId,
          variables.sessionId
        ),
      });
    },
  });
};

/**
 * セッションを作成するフック
 */
export const useCreateSession = (): UseMutationResult<
  CreateSessionResponse,
  Error,
  { appName: string; userId: string },
  unknown
> => {
  return useMutation({
    mutationFn: ({ appName, userId }) => createSession(appName, userId),
  });
};

/**
 * 任意のIDでセッションを作成するフック
 */
export const useCreateSessionWithId = (): UseMutationResult<
  CreateSessionResponse,
  Error,
  { appName: string; userId: string; sessionId: string },
  unknown
> => {
  return useMutation({
    mutationFn: ({ appName, userId, sessionId }) =>
      createSessionWithId(appName, userId, sessionId),
  });
};

/**
 * セッション情報を取得するフック
 */
export const useGetSession = (
  appName: string,
  userId: string,
  sessionId: string,
  enabled = true
): UseQueryResult<SessionInfo, Error> => {
  return useQuery({
    queryKey: queryKeys.session(appName, userId, sessionId),
    queryFn: () => getSession(appName, userId, sessionId),
    enabled: enabled && !!appName && !!userId && !!sessionId,
    staleTime: 30 * 1000, // 30秒
  });
};

/**
 * 生成済み画像一覧
 */
export const useListImages = (): UseQueryResult<
  import("./types").ListImagesResponse,
  Error
> => {
  return useQuery({
    queryKey: queryKeys.images,
    queryFn: listImages,
    refetchInterval: 60 * 1000, // 1分毎に更新
  });
};
