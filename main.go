package main

import (
	"context"
	"log"
	"time"
	"unicode/utf8"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/omegaatt36/chatgpt-telegram/app"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var config struct {
	telegramBotToken string
	apiKey           string
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

	client := gpt3.NewClient(config.apiKey)

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		if utf8.RuneCountInString(c.Message().Text) == 0 {
			return nil
		}

		var res string
		err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
			Prompt:      []string{c.Message().Text},
			MaxTokens:   gpt3.IntPtr(3000),
			Temperature: gpt3.Float32Ptr(0),
		}, func(resp *gpt3.CompletionResponse) {
			res += resp.Choices[0].Text
		})
		if err != nil {
			if ierr := c.Send(err.Error()); err != nil {
				log.Println(ierr)
			}
		}

		return c.Send(res)
	})

	go func() {
		bot.Start()
	}()
	<-ctx.Done()
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
	)

	app.Run()
}
