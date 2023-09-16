package app

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ExpenseCommand = "expense"
)

type Service struct {
	notion NotionHandler
	tgBot  *tgbotapi.BotAPI
	logger *slog.Logger
}

type NotionHandler interface {
	Add(amount float64, category, title string) error
}

func NewService(notion NotionHandler, tgBot *tgbotapi.BotAPI, logger *slog.Logger) *Service {
	return &Service{
		notion: notion,
		tgBot:  tgBot,
		logger: logger,
	}
}

func (s *Service) HandleWebhook(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}
	if update.Message != nil {
		text := update.Message.Text
		cmd := WhichCommand(text)
		switch cmd {
		case ExpenseCommand:
			{
				amount, category, title := s.extract(text)
				if amount != 0 && category != "" {
					err := s.notion.Add(amount, category, title)
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
					messageText := fmt.Sprintf("B %.2f in %s added.", amount, category)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
					s.tgBot.Send(msg)
					s.logger.Info("value is added")
				}
			}
		}
	}
	return nil
}

func WhichCommand(s string) string {
	if checkIsExpenseCommand(s) {
		return ExpenseCommand
	}
	return "unknown"
}

func checkIsExpenseCommand(s string) bool {
	reg := regexp.MustCompile(`^(\d+(?:\.\d{1,2})?)([ftcgbm])\s*(.+)?$`)
	subMatchs := reg.FindStringSubmatch(s)
	return len(subMatchs) > 2
}

func (srv *Service) extract(s string) (amount float64, category, title string) {
	reg := regexp.MustCompile(`^(\d+(?:\.\d{1,2})?)([ftcgbm])\s*(.+)?$`)
	subMatchs := reg.FindStringSubmatch(s)
	if len(subMatchs) < 3 {
		srv.logger.Info("Cannot extract amount, category and title from command message",
			slog.String("subMathches", strings.Join(subMatchs, ",")))
		return 0, "", ""
	}
	amountStr := subMatchs[1]
	category = getCategory(subMatchs[2])
	if len(subMatchs) > 3 && subMatchs[3] != "" {
		title = subMatchs[3]
	}

	amount, _ = strconv.ParseFloat(amountStr, 64)
	return amount, category, title
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
