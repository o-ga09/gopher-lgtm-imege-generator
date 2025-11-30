# Go Gopher LGTM Image Generator

Go Programming LanguageのマスコットキャラクターであるGo GopherくんのLGTM画像を生成するAIエージェントです。

## 動作デモ

https://github.com/user-attachments/assets/055034cc-dfc6-4b2e-964b-e956b783a8a6

## 技術スタック

### バックエンド

- Go
- Google Gemini API (Imagen 3)
- Cloudflare R2 (画像ストレージ)
- AWS SDK for Go v2 (S3互換クライアント)

### フロントエンド

- React + TypeScript
- Vite
- Tailwind CSS v4
- Tanstack Query
- Tanstack Router
- Lucide React (アイコン)

## セットアップ

### 前提条件

- Go 1.23以上
- Node.js 20以上
- pnpm
- Gemini API キー
- Cloudflare R2 バケットとクレデンシャル

### バックエンド

1. 環境変数の設定

```bash
cd backend
cp .env.example .env
```

`.env` ファイルを編集して、以下の値を設定してください：

```env
ENV=DEV
PORT=8080
GEMINI_API_KEY=your_gemini_api_key_here
CLOUDFLARE_R2_ACCOUNT_ID=your_account_id_here
CLOUDFLARE_R2_ACCESSKEY=your_access_key_here
CLOUDFLARE_R2_SECRETKEY=your_secret_key_here
CLOUDFLARE_R2_BUCKET_NAME=your_bucket_name_here
CLOUDFLARE_R2_ENDPOINT=https://your_account_id.r2.cloudflarestorage.com
CLOUDFLARE_R2_PUBLIC_URL=https://pub-your_public_url.r2.dev
CLOUDFLARE_R2_REGION=auto
```

1. 依存関係のインストール

```bash
go mod download
```

3. サーバーの起動

```bash
go run cmd/agent/main.go
```

サーバーは `http://localhost:8080` で起動します。

### フロントエンド

1. 依存関係のインストール

```bash
cd frontend
pnpm install
```

2. 開発サーバーの起動

```bash
pnpm dev
```

フロントエンドは `http://localhost:5173` で起動します。

## 使い方

1. ブラウザで `http://localhost:5173` にアクセス
2. プロンプト入力欄にGo Gopherの画像生成指示を入力（例: "Go Gopher giving thumbs up"）
3. "Generate LGTM" ボタンをクリック
4. 生成された画像が表示されます

## API エンドポイント

### GET /v1/agent/list-apps

AIエージェント一覧を取得します。

**レスポンス:**

```json
{
  "apps": ["lgtm_image_generator"]
}
```

### POST /v1/agent/apps/{app_name}/users/{user_id}/sessions

セッションを作成します。

**パラメータ:**
- `app_name`: アプリケーション名（例: `lgtm_image_generator`）
- `user_id`: ユーザーID（任意の文字列）

**レスポンス:**

```json
{
  "sessionId": "a2f792ad-cf72-4991-ad4a-2724159f0633"
}
```

### POST /v1/agent/apps/{app_name}/users/{user_id}/sessions/{session_id}

任意のIDでセッションを作成します。

**パラメータ:**
- `app_name`: アプリケーション名
- `user_id`: ユーザーID
- `session_id`: セッションID（指定したIDで作成）

### GET /v1/agent/apps/{app_name}/users/{user_id}/sessions/{session_id}

セッション情報を取得します。

**レスポンス:**

```json
{
  "sessionId": "a2f792ad-cf72-4991-ad4a-2724159f0633",
  "appName": "lgtm_image_generator",
  "userId": "user-123",
  "createdAt": "2025-01-01T00:00:00Z"
}
```

### POST /v1/agent/run

AIエージェントを実行してLGTM画像を生成します。

**リクエスト:**

```json
{
  "appName": "lgtm_image_generator",
  "userId": "user-123",
  "sessionId": "session-id",
  "newMessage": {
    "role": "user",
    "parts": [
      {
        "text": "Go Gopher giving thumbs up"
      }
    ]
  }
}
```

**レスポンス:**

```json
{
  "message": {
    "role": "model",
    "parts": [
      {
        "text": "画像を生成しました"
      }
    ]
  }
}
```

### POST /v1/agent/run_sse

SSE(Server-Sent Events)でAIエージェントを実行します。

リクエストボディは `/v1/agent/run` と同じです。

### GET /v1/images

アップロードされた画像の一覧を取得します。

**レスポンス:**

```json
{
  "images": [
    {
      "key": "lgtm-1234567890.png",
      "url": "https://your-r2-public-url/lgtm-1234567890.png",
      "lastModified": "2025-01-01T00:00:00Z",
      "size": 123456
    }
  ]
}
```

### GET /health

ヘルスチェックエンドポイント。

## ディレクトリ構成

```
.
├── backend/
│   ├── cmd/
│   │   └── agent/
│   │       └── main.go          # エントリーポイント
│   ├── internal/
│   │   ├── agent/
│   │   │   └── agent.go         # 画像生成ロジック
│   │   └── server/
│   │       └── server.go        # HTTPサーバー
│   ├── go.mod
│   └── .env.example
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   └── ImageGenerator.tsx
│   │   ├── routes/
│   │   │   ├── __root.tsx
│   │   │   └── index.tsx
│   │   ├── main.tsx
│   │   └── index.css
│   ├── package.json
│   └── vite.config.ts
└── .github/
    └── workflows/
        └── test.yml             # CI設定
```

## ライセンス

MIT
