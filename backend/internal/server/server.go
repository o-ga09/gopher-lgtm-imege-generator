package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	EnvCfg "github.com/o-ga09/gopher-lgtm-image-generator/pkg/config"
	"github.com/o-ga09/gopher-lgtm-image-generator/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/server/adkrest"
	"google.golang.org/adk/session"
)

type Server struct {
	srv *http.Server
	cfg *launcher.Config
}

// NewServer creates a new agent server
func NewServer(ctx context.Context, root agent.Agent, sub ...agent.Agent) (*Server, error) {
	var agentLoader agent.Loader
	var err error
	if len(sub) > 0 {
		agents := make([]agent.Agent, 0, len(sub))
		agents = append(agents, sub...)
		agentLoader, err = agent.NewMultiLoader(root, agents...)
		if err != nil {
			log.Fatalf("failed to create multi agent loader: %v", err)
		}
	} else {
		agentLoader = agent.NewSingleLoader(root)
	}
	// Configure the ADK REST API
	config := &launcher.Config{
		AgentLoader:     agentLoader,
		SessionService:  session.InMemoryService(),
		ArtifactService: artifact.InMemoryService(),
	}

	// Create the REST API handler
	apiHandler := adkrest.NewHandler(config)

	// Create a mux for routing
	mux := http.NewServeMux()

	// Register the API handler at the /v1/agent/ path and inject DB session via middleware
	mux.Handle("/v1/agent/", http.StripPrefix("/v1/agent", apiHandler))

	// Images history endpoint (list all uploaded images in the bucket)
	mux.HandleFunc("/v1/images", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		env := EnvCfg.GetCtxEnv(ctx)
		// build s3 client
		accessKey := env.CLOUDFLARE_R2_ACCESSKEY
		secretKey := env.CLOUDFLARE_R2_SECRETKEY
		endpoint := env.CLOUDFLARE_R2_ENDPOINT
		region := env.CLOUDFLARE_R2_REGION
		if region == "" {
			region = "auto"
		}

		creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
		awsCfg := aws.Config{Region: region, Credentials: creds}
		s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			if endpoint != "" {
				o.BaseEndpoint = aws.String(endpoint)
			}
			o.UsePathStyle = true
		})

		bucket := env.CLOUDFLARE_R2_BUCKET_NAME
		if bucket == "" {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"bucket not configured"}`))
			return
		}

		// list objects (limit to first 200 for now)
		resp, err := s3Client.ListObjectsV2(r.Context(), &s3.ListObjectsV2Input{
			Bucket:  aws.String(bucket),
			MaxKeys: aws.Int32(200),
		})
		if err != nil {
			log.Printf("failed to list objects: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to list objects"}`))
			return
		}

		// Use public URL from config or fallback to default format
		baseURL := env.CLOUDFLARE_R2_PUBLIC_URL
		if baseURL == "" {
			// Default format: https://[ACCOUNT_ID].r2.cloudflarestorage.com/[BUCKET_NAME]
			accountID := env.CLOUDFLARE_R2_ACCOUNT_ID
			if accountID != "" {
				baseURL = fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", accountID, bucket)
			} else {
				baseURL = fmt.Sprintf("https://r2.cloudflarestorage.com/%s", bucket)
			}
		}

		type imageInfo struct {
			Key          string `json:"key"`
			URL          string `json:"url"`
			Size         int64  `json:"size"`
			LastModified string `json:"lastModified"`
		}
		out := struct {
			Images []imageInfo `json:"images"`
		}{Images: []imageInfo{}}

		for _, obj := range resp.Contents {
			if obj.Key == nil {
				continue
			}
			if !strings.HasSuffix(*obj.Key, ".png") {
				continue
			}
			lm := ""
			if obj.LastModified != nil {
				lm = obj.LastModified.Format(time.RFC3339)
			}
			size := int64(0)
			if obj.Size != nil {
				size = *obj.Size
			}
			out.Images = append(out.Images, imageInfo{
				Key:          *obj.Key,
				URL:          fmt.Sprintf("%s/%s", baseURL, *obj.Key),
				Size:         size,
				LastModified: lm,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		_ = enc.Encode(out)
	})

	// Add a health check endpoint
	mux.HandleFunc("/v1/agent/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	env := EnvCfg.GetCtxEnv(ctx)
	port := fmt.Sprintf(":%s", env.Port)

	// CORSミドルウェアを適用
	handler := corsMiddleware(mux, env)

	return &Server{
		srv: &http.Server{
			Handler: handler,
			Addr:    port,
		},
		cfg: config,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	// サーバーの起動
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, fmt.Sprintf("Failed to listen and serve: %v", err))
		}
	}()

	logger.Info(ctx, fmt.Sprintf("Server is running on %s", s.srv.Addr))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "graceful shutdown")

	// サーバーのタイムアウト設定
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// サーバーのシャットダウン
	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to shutdown server: %v", err))
		return err
	}
	return nil
}

func (s *Server) DebugServer(ctx context.Context) error {
	l := full.NewLauncher()
	err := l.Execute(ctx, s.cfg, []string{"web", "webui", "api"})
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
	return nil
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler, env *EnvCfg.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// 許可するオリジンを取得
		allowedOrigins := strings.Split(env.AllowedOrigins, ",")

		// オリジンチェック
		var isAllowed bool
		for _, allowed := range allowedOrigins {
			allowed = strings.TrimSpace(allowed)
			if allowed == "*" || allowed == origin {
				isAllowed = true
				break
			}
		}

		// prod環境では厳密にチェック、それ以外は許可
		if env.Env == "prod" {
			if isAllowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else {
			// local/dev環境ではすべて許可
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// プリフライトリクエストの処理
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
