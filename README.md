# Go Gopher LGTM Image Generator

Go Programming LanguageのマスコットキャラクターであるGo GopherくんのLGTM画像を生成するAIエージェントです。

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
GEMINI_API_KEY=your_gemini_api_key
R2_ACCOUNT_ID=your_cloudflare_account_id
R2_ACCESS_KEY_ID=your_r2_access_key_id
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
R2_BUCKET_NAME=your_r2_bucket_name
R2_PUBLIC_URL=https://your-bucket-public-url

# Server Configuration
ENV=local
PORT=8080

# CORS Configuration (本番環境では必ず設定してください)
ALLOWED_ORIGINS=*
```

**CORS設定について:**

- **local/dev環境**: `ALLOWED_ORIGINS=*` で全てのオリジンを許可
- **prod環境**: 特定のオリジンのみを許可
  ```env
  ENV=prod
  ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
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

### POST /api/generate

LGTM画像を生成します。

**リクエスト:**

```json
{
  "prompt": "Go Gopher doing something"
}
```

**レスポンス:**

```json
{
  "imageUrl": "https://your-r2-public-url/lgtm-1234567890.png"
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
