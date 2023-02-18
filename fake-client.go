package main

import (
	"context"
	"log"
	"time"
)

type fakeClient struct {
}

func (c *fakeClient) Stream(ctx context.Context, question string) (chan string, chan error) {
	res := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(res)
		defer close(errCh)

		for _, r := range "\nabcdefghijklmnopqustuvwxyz" {
			res <- string(r)
			time.Sleep(time.Millisecond * 50)
		}

		log.Println("fake client stream over")
	}()

	return res, errCh
}
