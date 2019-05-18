package axb

import (
	"os"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func do_shutdown(bot *Bot, msg *kbchat.SubscriptionMessage, args []string) error {
	bot.ReplyTo(msg, "Goodbye cruel world!")
	os.Exit(0)
	return nil
}

func do_print_admins(bot *Bot, msg *kbchat.SubscriptionMessage, args []string) error {
	var sb strings.Builder
	sb.WriteString("Admins are: ")
	isFirst := true
	for _, v := range bot.admins {
		if !isFirst {
			sb.WriteString(", ")
		}
		sb.WriteString(v)
		isFirst = false
	}
	return bot.ReplyTo(msg, sb.String())
}

func do_add_admin(bot *Bot, msg *kbchat.SubscriptionMessage, args []string) error {
	if len(args) != 2 {
		return bot.ReplyTo(msg, "syntax 'add_admin <username>'")
	}
	bot.admins = append(bot.admins, args[1])
	return bot.ReplyTo(msg, "Added %s as an admin", args[1])
}

func do_remove_admin(bot *Bot, msg *kbchat.SubscriptionMessage, args []string) error {
	var newadmins []string
	if len(args) != 2 {
		return bot.ReplyTo(msg, "syntax 'remove_admin <username>'")
	}

	for _, v := range bot.admins {
		if v == args[1] {
			continue
		}
		newadmins = append(newadmins, v)
	}

	bot.admins = newadmins
	return bot.ReplyTo(msg, "Removed %s as an admin", args[1])
}
