package server

import (
	"log/slog"
	"tg-task-shell/config"
	"tg-task-shell/shell"

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

func (s *Server) Start(c chan shell.Command) error {
	// TODO
	bot, err := tgbotapi.NewBotAPI(s.config.TG_API_TOKEN)
    if err != nil {
        return err
    }

    tasks := make(map[string]config.Task)

    for _, task := range s.config.Tasks {
        tasks[task.Name] = task
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
        cmd := update.Message.Command()

        if task, ok := tasks[cmd]; ok {
            values, err := task.ParseParamValues(cmd[len(task.Name):])
            if err != nil {
                slog.Error("failed to parse parameters", "err", err)
            }
            
            args := func () []string {
                result := make([]string, len(values))
                for _, value := range values {
                    result = append(result, value.Value)
                }
                return result
            }()

            c <- shell.Command{Name: task.Command, Args: args} 

            msg.Text = "OK"

        } else {
            slog.Warn("unknown command", "command", cmd)
            msg.Text = "I don't know that command"
        }

        if _, err := bot.Send(msg); err != nil {
            slog.Error("failed to send message", "err", err)
        }
    }

    return nil
}