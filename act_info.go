package axb

import (
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func do_help(bot *Bot, msg *kbchat.SubscriptionMessage, args []string) error {
	isAdmin := bot.isAdmin(msg.Message.Sender.Username)
	var sb strings.Builder
	sb.WriteString("You have access to the following commands: ")
	isFirst := true
	for k, v := range bot.commands {
		if v.AdminRequired == true && !isAdmin {
			continue
		}
		if !isFirst {
			sb.WriteString(",")
		}
		sb.WriteString("'")
		sb.WriteString(k)
		sb.WriteString("' ")
		isFirst = false
	}
	return bot.ReplyTo(msg, sb.String())
}
