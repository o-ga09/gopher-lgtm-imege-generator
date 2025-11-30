import { useState, useEffect, useMemo } from "react";
import { Loader2, Send, Images, History, Copy, Link } from "lucide-react";
import {
  useCreateSession,
  useRunAgent,
  useGetSession,
  useListImages,
} from "../api";

const APP_NAME = "gopher-lgtm-image-generator-agent"; // AIエージェント名
const USER_ID = "user-" + Math.random().toString(36).substring(7); // ランダムなユーザーID

// テキスト内のURLをリンクに変換するヘルパー関数
const linkifyText = (text: string) => {
  const urlRegex = /(https?:\/\/[^\s]+)/g;
  const parts = text.split(urlRegex);

  return parts.map((part, index) => {
    if (part.match(urlRegex)) {
      return (
        <a
          key={index}
          href={part}
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-600 hover:text-blue-800 underline break-all"
        >
          {part}
        </a>
      );
    }
    return <span key={index}>{part}</span>;
  });
};

export default function ImageGenerator() {
  const [prompt, setPrompt] = useState("");
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const [finalMessage, setFinalMessage] = useState<string>("");
  const [activeTab, setActiveTab] = useState<"generate" | "history">(
    "generate"
  );

  // 画像をクリップボードにコピー
  const copyImageToClipboard = async (url: string) => {
    try {
      // Canvas経由で画像を取得（CORS回避）
      const img = new Image();
      img.crossOrigin = "anonymous";

      await new Promise((resolve, reject) => {
        img.onload = resolve;
        img.onerror = reject;
        img.src = url;
      });

      const canvas = document.createElement("canvas");
      canvas.width = img.width;
      canvas.height = img.height;
      const ctx = canvas.getContext("2d");
      if (!ctx) throw new Error("Failed to get canvas context");

      ctx.drawImage(img, 0, 0);

      canvas.toBlob(async (blob) => {
        if (!blob) throw new Error("Failed to create blob");
        await navigator.clipboard.write([
          new ClipboardItem({ "image/png": blob }),
        ]);
      }, "image/png");
    } catch (err) {
      console.error("Failed to copy image:", err);
      // フォールバック: URLをコピー
      await navigator.clipboard.writeText(url);
    }
  };

  // セッション作成
  const createSession = useCreateSession();

  // AIエージェント実行
  const runAgent = useRunAgent();

  // 初回マウント時にセッションを作成
  useEffect(() => {
    if (!sessionId) {
      createSession.mutate(
        { appName: APP_NAME, userId: USER_ID },
        {
          onSuccess: (data) => {
            setSessionId(data.id);
          },
          onError: (error) => {
            console.error("Failed to create session:", error);
          },
        }
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sessionId]);

  // セッション詳細 (イベント含む)
  const sessionInfo = useGetSession(
    APP_NAME,
    USER_ID,
    sessionId || "",
    !!sessionId
  );
  const imagesHistory = useListImages();

  const parseEvents = useMemo(() => {
    const events: any[] = sessionInfo.data?.events || [];
    let textFromModel = "";
    // 最終 model ロールのテキスト抽出
    for (const ev of events) {
      if (ev.content?.parts) {
        for (const p of ev.content.parts) {
          if (p.text) {
            textFromModel = p.text; // 最後に出たものを保持
          }
        }
      }
    }
    return { events, textFromModel };
  }, [sessionInfo.data]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!prompt || !sessionId) return;

    try {
      await runAgent.mutateAsync({
        appName: APP_NAME,
        userId: USER_ID,
        sessionId: sessionId,
        newMessage: {
          role: "user",
          parts: [{ text: prompt }],
        },
      });

      // セッション再取得後に useEffect でパースされる想定
    } catch (error) {
      console.error("Failed to generate image:", error);
    }
  };

  // イベントから最終テキストと画像URL抽出
  useEffect(() => {
    const text = parseEvents.textFromModel;
    if (text) {
      setFinalMessage(text);
      // バッククォートやスペース等を考慮したURL抽出
      const urlMatch = text.match(/https?:\/\/[^\s`'"]+/i);
      if (urlMatch) {
        const url = urlMatch[0];
        setImageUrl(url);
      } else {
        console.error("No URL found in message");
      }
    }
  }, [parseEvents]);

  const isLoading = runAgent.isPending;
  const isCreatingSession = createSession.isPending;
  const hasError = createSession.error || runAgent.error;

  // セッション作成処理
  const handleCreateNewSession = () => {
    setImageUrl(null);
    setFinalMessage("");
    setPrompt("");
    setSessionId(null); // これで useEffect が再実行される
  };

  // セッションエラーハンドリング
  useEffect(() => {
    if (sessionInfo.error) {
      setSessionId(null); // 新しいセッションを作成
    }
  }, [sessionInfo.error]);

  return (
    <div className="max-w-3xl mx-auto p-6 bg-white rounded-xl shadow-lg max-h-screen overflow-y-auto">
      {/* タブヘッダー */}
      <div className="flex gap-2 mb-4 border-b pb-2 items-center justify-between">
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setActiveTab("generate")}
            className={`flex items-center gap-1 px-3 py-2 rounded-md text-sm font-medium transition-colors ${
              activeTab === "generate"
                ? "bg-blue-600 text-white"
                : "bg-gray-100 text-gray-700 hover:bg-gray-200"
            }`}
          >
            <Images className="h-4 w-4" /> 生成
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("history")}
            className={`flex items-center gap-1 px-3 py-2 rounded-md text-sm font-medium transition-colors ${
              activeTab === "history"
                ? "bg-blue-600 text-white"
                : "bg-gray-100 text-gray-700 hover:bg-gray-200"
            }`}
          >
            <History className="h-4 w-4" /> 履歴
          </button>
        </div>
        {activeTab === "generate" && (
          <button
            type="button"
            onClick={handleCreateNewSession}
            disabled={isCreatingSession}
            className="text-xs px-3 py-1.5 rounded bg-gray-100 hover:bg-gray-200 disabled:opacity-50 font-medium"
          >
            New Session
          </button>
        )}
      </div>

      {activeTab === "generate" && (
        <>
          <form onSubmit={handleSubmit} className="space-y-3">
            <div>
              <label
                htmlFor="prompt"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                Prompt
              </label>
              <textarea
                id="prompt"
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                className="w-full h-24 p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none text-sm"
                placeholder="Describe the Go Gopher doing something..."
                disabled={isCreatingSession || isLoading || !sessionId}
              />
            </div>

            <button
              type="submit"
              disabled={isCreatingSession || isLoading || !prompt || !sessionId}
              className="w-full flex items-center justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isCreatingSession || isLoading ? (
                <>
                  <Loader2 className="animate-spin -ml-1 mr-2 h-5 w-5" />
                  {isCreatingSession ? "Creating Session..." : "Generating..."}
                </>
              ) : (
                <>
                  <Send className="-ml-1 mr-2 h-5 w-5" />
                  Generate LGTM
                </>
              )}
            </button>
          </form>

          {hasError && (
            <div className="mt-4 p-4 bg-red-50 text-red-700 rounded-lg">
              Error:{" "}
              {((createSession.error || runAgent.error) as Error)?.message}
            </div>
          )}

          {imageUrl && (
            <div className="mt-4 space-y-3 animate-in fade-in duration-500">
              <h3 className="text-base font-medium text-gray-900">
                Generated Image
              </h3>
              <p className="text-xs text-gray-600 whitespace-pre-line break-words line-clamp-2">
                {linkifyText(finalMessage)}
              </p>
              <div className="relative aspect-square w-full max-w-xs overflow-hidden rounded-lg bg-gray-100 group">
                <img
                  src={imageUrl}
                  alt={prompt}
                  className="object-cover w-full h-full"
                />
                {/* 左上に小さく配置されたボタン */}
                <div className="absolute top-2 left-2 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                  <button
                    onClick={() => copyImageToClipboard(imageUrl)}
                    className="p-2 bg-white rounded-md shadow-lg hover:bg-gray-100 transition-colors"
                    title="Copy Image"
                  >
                    <Copy className="h-4 w-4" />
                  </button>
                  <button
                    onClick={() => navigator.clipboard.writeText(imageUrl)}
                    className="p-2 bg-white rounded-md shadow-lg hover:bg-gray-100 transition-colors"
                    title="Copy URL"
                  >
                    <Link className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          )}
        </>
      )}

      {activeTab === "history" && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">
              過去生成された画像
            </h3>
            <button
              onClick={() => imagesHistory.refetch()}
              className="text-xs px-2 py-1 rounded bg-gray-100 hover:bg-gray-200"
            >
              更新
            </button>
          </div>
          {imagesHistory.isLoading && (
            <div className="text-sm text-gray-500">読み込み中...</div>
          )}
          {imagesHistory.error && (
            <div className="text-sm text-red-600">
              履歴取得に失敗しました: {imagesHistory.error.message}
            </div>
          )}
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            {imagesHistory.data?.images.map((img) => (
              <div
                key={img.key}
                className="group border rounded-lg p-2 bg-gray-50 flex flex-col"
              >
                <div className="relative aspect-square w-full overflow-hidden rounded bg-white">
                  <img
                    src={img.url}
                    alt={img.key}
                    className="object-cover w-full h-full"
                  />
                  {/* 左上に小さく配置されたボタン */}
                  <div className="absolute top-2 left-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                    <button
                      onClick={() => copyImageToClipboard(img.url)}
                      className="p-1.5 bg-white rounded shadow-lg hover:bg-gray-100 transition-colors"
                      title="Copy Image"
                    >
                      <Copy className="h-3 w-3" />
                    </button>
                    <button
                      onClick={() => navigator.clipboard.writeText(img.url)}
                      className="p-1.5 bg-white rounded shadow-lg hover:bg-gray-100 transition-colors"
                      title="Copy URL"
                    >
                      <Link className="h-3 w-3" />
                    </button>
                  </div>
                </div>
                <div className="mt-2 flex flex-col gap-1">
                  <span
                    className="text-[10px] text-gray-600 truncate"
                    title={img.key}
                  >
                    {img.key}
                  </span>
                  <span className="text-[10px] text-gray-400">
                    {img.lastModified
                      ? new Date(img.lastModified).toLocaleString()
                      : ""}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
