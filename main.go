package main

import (
	"context"
	"log"
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

	service := chatgpttelegram.NewService(
		bot,
		telegram.NewTelegramBot(bot),
		chatgpt.NewChatGPTClient(gpt3.NewClient(config.apiKey), config.maxToken),
	)

	service.Start(ctx)
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
	)

	app.Run()
}
