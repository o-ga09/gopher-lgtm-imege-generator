// API型定義

// メッセージの部分
export interface MessagePart {
  text: string;
}

// メッセージ
export interface Message {
  role: "user" | "model";
  parts: MessagePart[];
}

// AIエージェント実行リクエスト
export interface RunAgentRequest {
  appName: string;
  userId: string;
  sessionId: string;
  newMessage: Message;
}

// AIエージェント実行レスポンス
export interface RunAgentResponse {
  sessionId: string;
  response: Message;
}

// セッション作成レスポンス
export interface CreateSessionResponse {
  id: string;
  appName: string;
  userId: string;
  lastUpdateTime: number;
  events: any[];
}

// セッション情報
export interface SessionInfo {
  sessionId: string;
  appName: string;
  userId: string;
  messages: Message[];
  createdAt: string;
  updatedAt: string;
  // ADK raw events (tool calls, function responses, final model message, etc.)
  events?: any[];
}

// AIエージェント情報
export interface AppInfo {
  name: string;
  description?: string;
}

// AIエージェント一覧レスポンス
export interface ListAppsResponse {
  apps: AppInfo[];
}

// 画像履歴 API
export interface ImageInfo {
  key: string;
  url: string;
  size: number;
  lastModified: string;
}

export interface ListImagesResponse {
  images: ImageInfo[];
}
