package server

import (
	"log/slog"
	"tg-task-shell/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Server struct {
	config *config.Config
}

func New(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start() error {
	// TODO
	bot, err := tgbotapi.NewBotAPI(s.config.TG_API_TOKEN)
    if err != nil {
        return err
    }

    bot.Debug = true

    slog.Info("Authorized on account", "UserName", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil { // ignore any non-Message updates
            continue
        }

        if !update.Message.IsCommand() { // ignore any non-command Messages
            continue
        }

        // Create a new MessageConfig. We don't have text yet,
        // so we leave it empty.
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

        // Extract the command from the Message.
        switch update.Message.Command() {
        case "help":
            msg.Text = "I understand /sayhi and /status."
        case "sayhi":
            msg.Text = "Hi :)"
        case "status":
            msg.Text = "I'm ok."
        default:
            msg.Text = "I don't know that command"
        }

        if _, err := bot.Send(msg); err != nil {
            slog.Error("failed to send message", "err", err)
        }
    }

    return nil
}