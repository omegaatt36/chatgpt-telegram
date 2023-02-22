package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/omegaatt36/chatgpt-telegram/app"
	chatgpttelegram "github.com/omegaatt36/chatgpt-telegram/app/chatgpt-telegram"
	chatgpt "github.com/omegaatt36/chatgpt-telegram/appmodule/chatgpt/repository"
	telegram "github.com/omegaatt36/chatgpt-telegram/appmodule/telegram/repository"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var config struct {
	telegramBotToken string
	apiKey           string
	maxToken         int
	timeout          int
	allowedUsers     []int64
}

// Main starts process in cli.
func Main(ctx context.Context) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.telegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	client := gpt3.NewClient(
		config.apiKey,
		gpt3.WithHTTPClient(
			&http.Client{
				Timeout: time.Duration(time.Duration(config.timeout) * time.Second),
			},
		))

	service := chatgpttelegram.NewService(
		bot,
		telegram.NewTelegramBot(bot),
		chatgpt.NewChatGPTClient(client, config.maxToken),
	)

	service.Start(ctx,
		chatgpttelegram.UseAllowedUsers{AllowedUsers: config.allowedUsers},
	)
	<-ctx.Done()
	log.Println("app stopping")
}

func main() {
	app := app.App{
		Main:  Main,
		Flags: []cli.Flag{},
	}

	app.Flags = append(app.Flags,
		&cli.StringFlag{
			Name:        "chatgpt-api-key",
			EnvVars:     []string{"CHATGPT_API_KEY"},
			Destination: &config.apiKey,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "telegram-bot-token",
			EnvVars:     []string{"TELEGRAM_BOT_TOKEN"},
			Destination: &config.telegramBotToken,
			Required:    true,
		},
		&cli.IntFlag{
			Name:        "chatgpt-max-token",
			EnvVars:     []string{"CHATGPT_MAX_TOKEN"},
			Destination: &config.maxToken,
			DefaultText: "3000",
			Value:       3000,
		},
		&cli.IntFlag{
			Name:        "chatgpt-timeout",
			EnvVars:     []string{"CHATGPT_TIMEOUT"},
			Destination: &config.timeout,
			DefaultText: "60",
			Value:       60,
		},
		&cli.MultiInt64Flag{
			Target: &cli.Int64SliceFlag{
				Name:    "allowed-users",
				EnvVars: []string{"ALLOWED_USERS"},
			},
			Value:       []int64{},
			Destination: &config.allowedUsers,
		},
	)

	app.Run()
}
