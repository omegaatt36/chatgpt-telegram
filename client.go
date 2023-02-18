package main

import (
	"context"
	"log"

	"github.com/PullRequestInc/go-gpt3"
)

type gptClient struct {
	gpt3.Client
}

func (c *gptClient) Stream(ctx context.Context, question string) (chan string, chan error) {
	res := make(chan string)
	errCh := make(chan error)
	go func() {
		defer close(res)
		defer close(errCh)
		err := c.CompletionStreamWithEngine(ctx,
			gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
				Prompt:      []string{question},
				MaxTokens:   gpt3.IntPtr(300),
				Temperature: gpt3.Float32Ptr(0),
			}, func(resp *gpt3.CompletionResponse) {
				res <- resp.Choices[0].Text
			})
		if err != nil {
			errCh <- err
		}

		log.Println("gpt client stream over")
	}()

	return res, errCh
}
