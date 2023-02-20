package chatgpttelegram

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	chatgpt "github.com/omegaatt36/chatgpt-telegram/appmodule/chatgpt/usecase"
	telegram "github.com/omegaatt36/chatgpt-telegram/appmodule/telegram/usecase"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	s := assert.New(t)

	var (
		ctx            = context.TODO()
		question       = "test"
		chatID   int64 = 1

		ch            = make(chan string)
		readOnlyCh    = func(c chan string) <-chan string { return c }(ch)
		readOnlyErrCh = make(<-chan error)
	)

	go func() {
		for _, r := range "answer" {
			ch <- string(r)
		}
	}()

	controller := gomock.NewController(t)
	mockTelegram := telegram.NewMockTelegramUseCase(controller)
	{
		mockTelegram.EXPECT().SendAsLiveOutput(chatID, readOnlyCh).
			Times(1).Return(nil)
	}

	mockChatGPT := chatgpt.NewMockChatGPTUseCase(controller)
	{
		mockChatGPT.EXPECT().Stream(ctx, question).Times(1).Return(readOnlyCh, readOnlyErrCh)
	}
	service := NewService(nil, mockTelegram, mockChatGPT)
	service.ctx = ctx
	s.NoError(service.processChatGPTQuestion(chatID, question))
}
