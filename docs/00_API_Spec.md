## GopherくんLGTM画像生成AIエージェントのAPI仕様書

このドキュメントは、GopherくんLGTM画像生成AIエージェントのAPI仕様書です。エージェントは、HTTPサーバーとして動作し、LGTM画像の生成リクエストを受け付けます。

### エンドポイント

- 基本URL: `http://<サーバーのホスト>:<ポート>/v1/agent`

#### エンドポイント一覧:

app_name: `list-apps`で取得できるAIエージェント名
session_id: `POST /apps/{app_name}/users/{user_id}/sessions`で作成したID
user_id: 任意の文字列

- `GET /list-apps`
  - AIエージェント一覧

- `POST /run`:
  - AIエージェント起動
  - リクエストボディ
```json
{
    "appName": "monhun_ai_agent",
    "userId": "user-123",
    "sessionId": "a2f792ad-cf72-4991-ad4a-2724159f0633",
    "newMessage": {
        "role": "user",
        "parts": [
            {
                "text": "モンスターで「あ」から始まるモンスターを5こ教えて？"
            }
        ]
    }
}
```


- `POST /run_sse`:
  - sseでAIエージェント起動
  - リクエストボディ
```json
{
    "appName": "monhun_ai_agent",
    "userId": "user-123",
    "sessionId": "a2f792ad-cf72-4991-ad4a-2724159f0633",
    "newMessage": {
        "role": "user",
        "parts": [
            {
                "text": "モンスターで「あ」から始まるモンスターを5こ教えて？"
            }
        ]
    }
}
```
- `POST /apps/{app_name}/users/{user_id}/sessions`:
  - セッション作成
- `POST /apps/{app_name}/users/{user_id}/sessions/{session_id}`:
  - 任意のIDでセッション作成
- `GET /apps/{app_name}/users/{user_id}/sessions/{session_id}`:
  - セッション情報取得
