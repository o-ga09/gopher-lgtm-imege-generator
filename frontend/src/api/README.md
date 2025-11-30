# API Client & Tanstack Query Hooks

このディレクトリには、バックエンドAPIと通信するためのAPIクライアントとTanstack Queryのカスタムフックが含まれています。

## 構成

- `client.ts` - AxiosベースのAPIクライアント設定
- `types.ts` - API型定義
- `agent.ts` - API関数の実装
- `hooks.ts` - Tanstack Queryカスタムフック
- `index.ts` - エクスポート

## 使用方法

### 1. 環境変数の設定

`.env`ファイルを作成し、API Base URLを設定します:

```env
VITE_API_BASE_URL=http://localhost:8080/v1/agent
```

### 2. AIエージェント一覧の取得

```tsx
import { useListApps } from '@/api';

function AppList() {
  const { data, isLoading, error } = useListApps();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <ul>
      {data?.apps.map((app) => (
        <li key={app.name}>{app.name}</li>
      ))}
    </ul>
  );
}
```

### 3. セッションの作成

```tsx
import { useCreateSession } from '@/api';

function CreateSessionButton() {
  const createSession = useCreateSession();

  const handleCreate = async () => {
    try {
      const result = await createSession.mutateAsync({
        appName: 'monhun_ai_agent',
        userId: 'user-123',
      });
    } catch (error) {
      console.error('Failed to create session:', error);
    }
  };

  return (
    <button onClick={handleCreate} disabled={createSession.isPending}>
      {createSession.isPending ? 'Creating...' : 'Create Session'}
    </button>
  );
}
```

### 4. AIエージェントの実行

```tsx
import { useRunAgent } from '@/api';

function ChatInterface() {
  const runAgent = useRunAgent();

  const handleSendMessage = async (message: string) => {
    try {
      const result = await runAgent.mutateAsync({
        appName: 'monhun_ai_agent',
        userId: 'user-123',
        sessionId: 'a2f792ad-cf72-4991-ad4a-2724159f0633',
        newMessage: {
          role: 'user',
          parts: [{ text: message }],
        },
      });
      console.log('Response:', result.response);
    } catch (error) {
      console.error('Failed to run agent:', error);
    }
  };

  return (
    <div>
      <button onClick={() => handleSendMessage('Hello!')} disabled={runAgent.isPending}>
        {runAgent.isPending ? 'Sending...' : 'Send Message'}
      </button>
    </div>
  );
}
```

### 5. セッション情報の取得

```tsx
import { useGetSession } from '@/api';

function SessionInfo({ appName, userId, sessionId }: {
  appName: string;
  userId: string;
  sessionId: string;
}) {
  const { data, isLoading, error } = useGetSession(appName, userId, sessionId);

  if (isLoading) return <div>Loading session...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h2>Session: {data?.sessionId}</h2>
      <div>
        {data?.messages.map((msg, idx) => (
          <div key={idx}>
            <strong>{msg.role}:</strong> {msg.parts[0]?.text}
          </div>
        ))}
      </div>
    </div>
  );
}
```

### 6. SSEでAIエージェントの実行

```tsx
import { useEffect, useState } from 'react';
import { runAgentSSE } from '@/api';

function StreamingChat() {
  const [messages, setMessages] = useState<string[]>([]);

  const handleStream = () => {
    const cleanup = runAgentSSE(
      {
        appName: 'monhun_ai_agent',
        userId: 'user-123',
        sessionId: 'session-id',
        newMessage: {
          role: 'user',
          parts: [{ text: 'Hello!' }],
        },
      },
      (data) => {
        // メッセージ受信時
        setMessages((prev) => [...prev, data]);
      },
      (error) => {
        // エラー時
        console.error('SSE Error:', error);
      },
      () => {
        // 完了時
        console.log('SSE Complete');
      }
    );

    // クリーンアップ
    return cleanup;
  };

  return (
    <div>
      <button onClick={handleStream}>Start Streaming</button>
      <div>
        {messages.map((msg, idx) => (
          <div key={idx}>{msg}</div>
        ))}
      </div>
    </div>
  );
}
```

## 型定義

すべての型は `types.ts` で定義されており、TypeScriptの型安全性を提供します。

```typescript
import type { Message, RunAgentRequest, SessionInfo } from '@/api';
```

## エラーハンドリング

APIクライアントには自動的なエラーハンドリングが組み込まれています:

- ネットワークエラー
- サーバーエラーレスポンス
- タイムアウト（30秒）

すべてのエラーはコンソールにログ出力され、呼び出し元でキャッチできます。
