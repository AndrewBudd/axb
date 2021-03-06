package axb

import (
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func doHelp(bot *Bot, msg *kbchat.SubscriptionMessage, args []string) error {
	isAdmin := bot.IsFromAdmin(msg)
	var sb strings.Builder
	sb.WriteString("You have access to the following commands: ")
	for k, v := range bot.commands {
		if v.AdminRequired == true && !isAdmin {
			continue
		}
		sb.WriteString("\n")
		sb.WriteString("'")
		sb.WriteString(k)
		sb.WriteString("'")
	}
	return bot.ReplyTo(msg, sb.String())
}
