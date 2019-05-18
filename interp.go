package axb

import (
	"sort"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

type BotCommand struct {
	function      func(*Bot, *kbchat.SubscriptionMessage, []string) error
	adminRequired bool
}

func (bot *Bot) isAdmin(user string) bool {
	for _, n := range bot.admins {
		if user == n {
			return true
		}
	}
	return false
}

func (bot *Bot) interp(msg *kbchat.SubscriptionMessage, message string) error {
	user := msg.Message.Sender.Username
	oneonone := true

	args := strings.Split(message, " ")
	// are you talking to me?
	if !strings.Contains(msg.Message.Channel.Name, ",") {
		if len(args) == 0 || args[0] != "@"+bot.API().GetUsername() {
			return nil
		}
		oneonone = false
		args = args[1:]
	}

	if len(args) == 0 && oneonone {
		return bot.ReplyTo(msg, "Huh???")
	}

	// dumb to do every time but it's just the beginning here
	keys := make([]string, 0)
	for k := range bot.commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if strings.HasPrefix(k, strings.ToLower(args[0])) {
			if bot.commands[k].adminRequired == true {
				if bot.isAdmin(user) {
					return bot.commands[k].function(bot, msg, args)
				}
				break
			}
			return bot.commands[k].function(bot, msg, args)
		}
	}
	if oneonone {
		return bot.ReplyTo(msg, "Huh???")
	}

	return nil
}
