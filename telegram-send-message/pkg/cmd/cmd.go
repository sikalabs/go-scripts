package cmd

import (
	"fmt"

	"github.com/sikalabs/go-scripts/telegram-send-message/pkg/telegram_send_message"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "telegram-send-message <message>",
	Short: "Send a message to Telegram",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := telegram_send_message.TelegramSendMessage(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Message sent to Telegram: " + args[0])
	},
}
