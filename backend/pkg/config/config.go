package config

import (
	"context"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/o-ga09/gopher-lgtm-image-generator/pkg/errors"
)

type Env string

const CtxEnvKey Env = "env"

type Config struct {
	Env                       string `env:"ENV" envDefault:"local"`
	Port                      string `env:"PORT" envDefault:"8080"`
	AllowedOrigins            string `env:"ALLOWED_ORIGINS" envDefault:"*"`
	ProjectID                 string `env:"PROJECT_ID" envDefault:""`
	CLOUDFLARE_R2_ACCOUNT_ID  string `env:"CLOUDFLARE_R2_ACCOUNT_ID" envDefault:""`
	CLOUDFLARE_R2_ACCESSKEY   string `env:"CLOUDFLARE_R2_ACCESSKEY" envDefault:""`
	CLOUDFLARE_R2_SECRETKEY   string `env:"CLOUDFLARE_R2_SECRETKEY" envDefault:""`
	CLOUDFLARE_R2_BUCKET_NAME string `env:"CLOUDFLARE_R2_BUCKET_NAME" envDefault:""`
	CLOUDFLARE_R2_ENDPOINT    string `env:"CLOUDFLARE_R2_ENDPOINT" envDefault:""`
	CLOUDFLARE_R2_PUBLIC_URL  string `env:"CLOUDFLARE_R2_PUBLIC_URL" envDefault:"http://localhost:9001"`
	CLOUDFLARE_R2_REGION      string `env:"CLOUDFLARE_R2_REGION" envDefault:"auto"`
	GeminiAPIKey              string `env:"GEMINI_API_KEY" envDefault:""`
	GeminiModel               string `env:"GEMINI_MODEL" envDefault:"gemini-2.5-flash"`
}

func New(ctx context.Context) (context.Context, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(ctx, err)
	}

	return context.WithValue(ctx, CtxEnvKey, cfg), nil
}

func GetCtxEnv(ctx context.Context) *Config {
	var cfg *Config
	var ok bool
	if cfg, ok = ctx.Value(CtxEnvKey).(*Config); !ok {
		log.Fatal("config not found")
	}
	return cfg
}
