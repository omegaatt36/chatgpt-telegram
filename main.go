package main

import (
	"context"
	"log"
	"time"
	"unicode/utf8"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/omegaatt36/chatgpt-telegram/app"
	"github.com/omegaatt36/chatgpt-telegram/appmodule/chatgpt"
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

	b := telegramBot{bot: bot, editInterval: 500 * time.Millisecond}

	source := gptClient{gpt3.NewClient(config.apiKey)}
	// source := fakeClient{}
	client := chatgpt.NewClient(&source)

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		if utf8.RuneCountInString(c.Message().Text) == 0 {
			return nil
		}

		log.Printf("start user(%d) prompt(%s)\n", c.Message().Chat.ID, c.Message().Text)
		defer func() { log.Println("done", c.Message().Text) }()

		messageCh, errCh := client.Stream(ctx, c.Message().Text)

		done := make(chan struct{}, 1)
		go func() {
			if err := b.SendAsLiveOutput(c.Message().Chat.ID, messageCh); err != nil {
				log.Println(err)
			}
			done <- struct{}{}
		}()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			case err, ok := <-errCh:
				if !ok {
					return nil
				}
				if err == nil {
					continue
				}

				if e := c.Send(err); e != nil {
					log.Fatal(err)
				}

				return nil
			}
		}
	})

	log.Println("starting telegram bot")
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
