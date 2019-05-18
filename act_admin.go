package axb

import (
	"os"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func do_shutdown(bot *Bot, msg *kbchat.SubscriptionMessage, args []string) error {
	bot.ReplyTo(msg, "Goodbye cruel world!")
	os.Exit(0)
	return nil
}
