package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"vassaneer-tg-bot/internal/app"
	"vassaneer-tg-bot/internal/notion"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Strting the application...")
	envs := loadEnv(logger)

	bot, err := tgbotapi.NewBotAPI(envs.TelegramBotToken)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Authorized on account %s", bot.Self.UserName)

	notionHandlers := make(map[string]*app.NotionHandlers)
	for k, v := range envs.NotionDatabaseID {
		notionHandlers[k] = app.NewHandler(notion.NewNotion(v, envs.NotionSecret, logger))
	}

	//notionHandlers["exercise"] = app.NewHandler(notion.NewNotion(envs.NotionDatabaseID, envs.NotionSecret, logger))

	service := app.NewService(notionHandlers, bot, logger)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	r.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		tgUpdate := new(tgbotapi.Update)
		if err := json.NewDecoder(r.Body).Decode(tgUpdate); err != nil {
			logger.Error("Error when decoding webhook update", slog.String("errorMessage", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		}
		if tgUpdate != nil && tgUpdate.Message != nil {
			logger.LogAttrs(context.Background(), slog.LevelInfo, "message details",
				slog.Int64("senderId", tgUpdate.Message.From.ID),
				slog.Int64("chatId", tgUpdate.Message.Chat.ID),
			)
			service.HandleWebhook(tgUpdate)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	http.ListenAndServe(":"+envs.Port, r)
}
