package main

import (
	"github.com/AJAYK-01/passwordless-go/passwordless"
	"github.com/alexedwards/scs/v2"
	"log"
	"log/slog"
	"net/http"
	"os"
	"webauthn.rasc.ch/internal/config"
)

type application struct {
	config             *config.Config
	sessionManager     *scs.SessionManager
	passwordlessClient *passwordless.Client
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("reading config failed %v\n", err)
	}

	var logger *slog.Logger

	switch cfg.Environment {
	case config.Development:
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case config.Production:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)

	passwordlessClient := passwordless.Client{ApiSecret: cfg.Passwordless.SecretApiKey,
		BaseUrl: cfg.Passwordless.ApiUrl}

	sm := scs.New()
	sm.Lifetime = cfg.Session.Lifetime
	sm.Cookie.SameSite = http.SameSiteStrictMode
	if cfg.Session.CookieDomain != "" {
		sm.Cookie.Domain = cfg.Session.CookieDomain
	}
	sm.Cookie.Secure = cfg.Session.SecureCookie
	slog.Info("secure cookie", "secure", sm.Cookie.Secure)

	app := &application{
		config:             &cfg,
		sessionManager:     sm,
		passwordlessClient: &passwordlessClient,
	}

	err = app.serve()
	if err != nil {
		slog.Error("http serve failed", "error", err)
		os.Exit(1)
	}
}
