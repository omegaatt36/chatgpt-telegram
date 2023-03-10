//go:generate mockgen -package=usecase -destination=chatgpt_usecase_mock.go . ChatGPTUseCase

package usecase

import "context"

// ChatGPTUseCase defines ChatGPT send question use case.
type ChatGPTUseCase interface {
	Stream(ctx context.Context, question string) (<-chan string, <-chan error)
}
