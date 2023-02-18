package main

import (
	"log"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

type telegramBot struct {
	bot          *telebot.Bot
	editInterval time.Duration
}

func ensureFormatting(text string) string {
	numDelimiters := strings.Count(text, "```")
	numSingleDelimiters := strings.Count(strings.Replace(text, "```", "", -1), "`")

	if (numDelimiters % 2) == 1 {
		text += "```"
	}
	if (numSingleDelimiters % 2) == 1 {
		text += "`"
	}

	return text
}

func (b *telegramBot) SendAsLiveOutput(chatID int64, feed chan string) error {
	var message *telebot.Message
	var lastResp, tmpResp string
	var done bool

	// aggregate message
	go func() {
		for {
			s, ok := <-feed
			tmpResp += s
			if !ok {
				done = true
				return
			}
		}
	}()

	defer func() { log.Println("SendAsLiveOutput done") }()

	chat := &telebot.Chat{ID: chatID}

	send := func(register string) error {
		defer func() {
			lastResp = register
		}()

		if len(strings.Trim(register, "\n")) == 0 {
			time.Sleep(time.Microsecond * 50)
			return nil
		}

		if message == nil {
			var err error
			if message, err = b.bot.Send(chat, register); err != nil {
				return err
			}
		} else {
			if register == lastResp {
				return nil
			}

			text := ensureFormatting(register)
			if _, err := b.bot.Edit(message, text); err != nil {
				return err
			}
		}

		return nil
	}

	for {
		if done {
			return send(tmpResp)
		}

		if tmpResp == "" {
			continue
		}

		if err := send(tmpResp); err != nil {
			return err
		}
	}
}
