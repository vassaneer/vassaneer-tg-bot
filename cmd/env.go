package main

import (
	"log/slog"
	"os"
)

type Envs struct {
	TelegramBotToken string
	NotionSecret     string
	NotionDatabaseID string
	Port             string
}

func loadEnv(logger *slog.Logger) Envs {
	tgBotToken, ok := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !ok {
		logger.Error("TELEGRAM_BOT_TOKEN env isn't set")
		os.Exit(1)
	}
	notionSecret, ok := os.LookupEnv("NOTION_SECRET")
	if !ok {
		logger.Error("NOTION_SECRET env isn't set")
		os.Exit(1)
	}
	notionDatabase, ok := os.LookupEnv("NOTION_DATABASE_ID")
	if !ok {
		logger.Error("NOTION_DATABASE_ID env isn't set")
		os.Exit(1)
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		logger.Error("PORT env isn't set")
		os.Exit(1)
	}
	return Envs{
		TelegramBotToken: tgBotToken,
		NotionSecret:     notionSecret,
		NotionDatabaseID: notionDatabase,
		Port:             port,
	}
}
