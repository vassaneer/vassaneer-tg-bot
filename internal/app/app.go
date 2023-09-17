package app

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	notions map[string]*NotionHandlers
	tgBot   *tgbotapi.BotAPI
	logger  *slog.Logger
}

type NotionHandlers struct {
	notion NotionHandler
}

type NotionHandler interface {
	Add(cmd ReturnCommand) error
}

func NewHandler(notion NotionHandler) *NotionHandlers {
	return &NotionHandlers{
		notion: notion,
	}
}

func NewService(notions map[string]*NotionHandlers, tgBot *tgbotapi.BotAPI, logger *slog.Logger) *Service {
	return &Service{
		notions: notions,
		tgBot:   tgBot,
		logger:  logger,
	}
}

func (s *Service) HandleWebhook(update *tgbotapi.Update) error {
	//
	ExpenseCommand := NewCommand("expense", `^(\d+(?:\.\d{1,2})?)([ftcgbm])\s*(.+)?$`)
	ExerciseCommand := NewCommand("excercise", ``)
	RepCommand := NewCommand("rep", `^(\w+)\s*(\d+)[X](\d+)`)
	CommandHandler := NewCommandHandler([]*Command{ExpenseCommand, ExerciseCommand, RepCommand})
	if update.Message == nil {
		return nil
	}
	if update.Message != nil {
		text := update.Message.Text
		cmd := CommandHandler.WhichCommand(text)
		// case cmd
		switch cmd {
		case ExpenseCommand:
			{
				resp := ExpenseCommand.extract(text, s, ExpenseCommandExtract)
				err := s.notions["expense"].notion.Add(resp)
				if err != nil {
					s.logger.Error("Error when adding expense to Notion",
						slog.String("errorMessage", err.Error()),
						slog.String("command", "expense"),
						slog.String("commandMessage", text),
					)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "something happens with Notion. check system log in Cloud Run Logs.")
					msg.ReplyToMessageID = update.Message.MessageID
					s.tgBot.Send(msg)
					return nil
				}
				messageText := fmt.Sprintf("B %.2f in %s added.", resp.Fields["amount"].Value, resp.Fields["category"].Value)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
				s.tgBot.Send(msg)
				s.logger.Info("value is added")
			}
		case RepCommand:
			{
				resp := RepCommand.extract(text, s, RepCommandExtract)
				err := s.notions["rep"].notion.Add(resp)
				if err != nil {
					s.logger.Error("Error when adding expense to Notion",
						slog.String("errorMessage", err.Error()),
						slog.String("command", "expense"),
						slog.String("commandMessage", text),
					)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "something happens with Notion. check system log in Cloud Run Logs.")
					msg.ReplyToMessageID = update.Message.MessageID
					s.tgBot.Send(msg)
					return nil
				}
				messageText := fmt.Sprintf("Workout %s is %fX%f", resp.Fields["Name"].Value, resp.Fields["Weight"].Value, resp.Fields["Rep"].Value)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
				s.tgBot.Send(msg)
				s.logger.Info("value is added")
			}
		}
	}
	return nil
}

func getCategory(s string) string {
	switch s {
	case "b":
		return "beverage"
	case "f":
		return "food"
	case "t":
		return "transport"
	case "c":
		return "clothes"
	case "g":
		return "grocery"
	case "m":
		return "misc"
	default:
		return "unknown"
	}
}
