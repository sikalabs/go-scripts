package telegram_send_message

import (
	"github.com/sikalabs/slu/utils/telegram_utils"
	"github.com/sikalabsx/sikalabs-encrypted-go/pkg/encrypted"
)

func TelegramSendMessage(message string) error {
	telegramConfig, err := encrypted.GetSikalabsNotificationV1TelegramConfig()
	if err != nil {
		return err
	}
	return telegram_utils.TelegramSendMessage(telegramConfig.Token, telegramConfig.ChatID, message)
}
