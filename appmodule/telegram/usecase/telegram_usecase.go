//go:generate mockgen -package=usecase -destination=telegram_usecase_mock.go . TelegramUseCase

package usecase

// TelegramUseCase defines telegram send message use case.
type TelegramUseCase interface {
	SendAsLiveOutput(chatID int64, feed <-chan string) error
}
