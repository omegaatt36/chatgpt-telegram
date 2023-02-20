package chatgpttelegram

import (
	"context"
	"log"
	"unicode/utf8"

	chatgpt "github.com/omegaatt36/chatgpt-telegram/appmodule/chatgpt/usecase"
	telegram "github.com/omegaatt36/chatgpt-telegram/appmodule/telegram/usecase"
	"gopkg.in/telebot.v3"
)

type Service struct {
	ctx context.Context

	bot      *telebot.Bot
	telegram telegram.TelegramUseCase
	gpt      chatgpt.ChatGPTUseCase
}

type Confit struct {
}

func NewService(bot *telebot.Bot, tu telegram.TelegramUseCase, cu chatgpt.ChatGPTUseCase) *Service {
	return &Service{
		bot:      bot,
		telegram: tu,
		gpt:      cu,
	}
}

func (s *Service) registerEndpoint() {
	s.bot.Handle(telebot.OnText, s.chatGPTQuestion)
}

func (s *Service) Start(ctx context.Context) {
	s.ctx = ctx

	s.registerEndpoint()

	go func() {
		log.Println("starting telegram bot")
		s.bot.Start()
	}()

	go func() {
		<-ctx.Done()
		log.Println("stopping telegram bot")
		s.bot.Stop()
		log.Println("telegram bot is stopped")
	}()
}

func (s Service) processChatGPTQuestion(chatID int64, question string) error {
	messageCh, errChatGPTCh := s.gpt.Stream(s.ctx, question)

	errCh := make(chan error, 1)
	defer close(errCh)

	done := make(chan struct{}, 1)
	go func() {
		if err := s.telegram.SendAsLiveOutput(chatID, messageCh); err != nil {
			errCh <- err
		}
		done <- struct{}{}
	}()

	for {
		select {
		case <-s.ctx.Done():
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

			return err
		case err, ok := <-errChatGPTCh:
			if !ok {
				return nil
			}
			if err == nil {
				continue
			}

			return err
		}
	}
}

func (s *Service) chatGPTQuestion(c telebot.Context) error {
	if utf8.RuneCountInString(c.Message().Text) == 0 {
		return nil
	}

	log.Printf("start(%d) user(%d) prompt(%s)\n",
		c.Message().ID, c.Message().Chat.ID, c.Message().Text)
	defer func() { log.Printf("done(%d)\n", c.Message().ID) }()

	err := s.processChatGPTQuestion(c.Message().Chat.ID, c.Message().Text)
	if err != nil {
		if ierr := c.Send(err.Error()); ierr != nil {
			log.Fatalln(ierr)
		}
	}

	return err
}
