import axios from "axios";

// APIクライアントの設定
export const apiClient = axios.create({
  baseURL:
    import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/v1/agent",
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 30000, // 30秒
});

// リクエストインターセプター
apiClient.interceptors.request.use(
  (config) => {
    // 必要に応じて認証トークンなどを追加
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// レスポンスインターセプター
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // エラーハンドリング
    if (error.response) {
      // サーバーからのエラーレスポンス
      console.error("API Error:", error.response.data);
    } else if (error.request) {
      // リクエストが送信されたがレスポンスがない
      console.error("Network Error:", error.request);
    } else {
      // その他のエラー
      console.error("Error:", error.message);
    }
    return Promise.reject(error);
  }
);
