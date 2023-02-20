package repository

import (
	"context"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/omegaatt36/chatgpt-telegram/appmodule/chatgpt/usecase"
)

type chatGPTClient struct {
	client gpt3.Client

	maxToken int
}

var _ usecase.ChatGPTUseCase = &chatGPTClient{}

// NewChatGPTClient returns implement of usecase.ChatGPTUseCase.
func NewChatGPTClient(client gpt3.Client, maxToken int) *chatGPTClient {
	return &chatGPTClient{
		client:   client,
		maxToken: maxToken,
	}
}

func (c *chatGPTClient) Stream(ctx context.Context, question string) (<-chan string, <-chan error) {
	res := make(chan string)
	errCh := make(chan error)
	go func() {
		defer close(res)
		defer close(errCh)
		err := c.client.CompletionStreamWithEngine(ctx,
			gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
				Prompt:      []string{question},
				MaxTokens:   gpt3.IntPtr(c.maxToken),
				Temperature: gpt3.Float32Ptr(0),
			}, func(resp *gpt3.CompletionResponse) {
				res <- resp.Choices[0].Text
			})
		if err != nil {
			errCh <- err
		}
	}()

	return res, errCh
}
