package chatgpt

import (
	"context"
)

type Client struct {
	useCase UseCase
}

type UseCase interface {
	Stream(context.Context, string) (chan string, chan error)
}

func NewClient(useCase UseCase) *Client {
	return &Client{
		useCase: useCase,
	}
}

func (c *Client) Stream(ctx context.Context, question string) (chan string, chan error) {
	return c.useCase.Stream(ctx, question)
}
