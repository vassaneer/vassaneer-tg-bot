package app

import (
	"errors"
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
	ArticleCommand := NewCommand("article", `^([^\s]+)\s*([fwx])\s*(.*)$`)
	CommandHandler := NewCommandHandler([]*Command{ExpenseCommand, ExerciseCommand, RepCommand, ArticleCommand})
	if update.Message == nil {
		return nil
	}
	if update.Message != nil {
		text := update.Message.Text
		cmd := CommandHandler.WhichCommand(text)
		messageText := ""
		err := errors.New("math: square root of negative number")
		// case cmd
		switch cmd {
		case ExpenseCommand:
			{
				resp := ExpenseCommand.extract(text, s, ExpenseCommandExtract)
				err = s.notions["expense"].notion.Add(resp)
				messageText = fmt.Sprintf("B %.2f in %s added.", resp.Fields["amount"].Value, resp.Fields["category"].Value)
			}
		case RepCommand:
			{
				resp := RepCommand.extract(text, s, RepCommandExtract)
				err = s.notions["rep"].notion.Add(resp)
				messageText = fmt.Sprintf("Workout %s is %.0fX%.0f", resp.Fields["Name"].Value, resp.Fields["Weight"].Value, resp.Fields["Rep"].Value)
			}
		case ArticleCommand:
			{
				resp := ArticleCommand.extract(text, s, ArticleCommandExtract)
				err = s.notions["article"].notion.Add(resp)
				messageText = fmt.Sprintf("Article %s from %s", resp.Fields["Name"].Value, resp.Fields["Link"].Value)

			}
		}
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
		s.tgBot.Send(msg)
		s.logger.Info("value is added")

	}
	return nil
}

func getSource(s string) string {
	switch s {
	case "f":
		return "Facebook"
	case "x":
		return "X"
	case "w":
		return "Website"
	default:
		return "unknown"
	}
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
