package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"runtime/debug"
)

type application struct {
	config            Config
	logger            *slog.Logger
	azureOpenAIClient *azopenai.Client
}

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	keyCredential := azcore.NewKeyCredential(cfg.AzureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(cfg.AzureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		return err
	}

	app := &application{
		config:            cfg,
		logger:            logger,
		azureOpenAIClient: client,
	}

	return app.serveHTTP()
}

type Config struct {
	HttpPort                  int
	AzureOpenAIKey            string
	AzureOpenAIEndpoint       string
	AzureOpenAIDeploymentName string
	AwsBedrockUserAccessKey   string
	AwsBedrockUserSecretKey   string
}

func loadConfig() (Config, error) {
	var cfg Config

	viper.SetDefault("httpPort", 4444)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
