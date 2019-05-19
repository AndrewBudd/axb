package axb

import (
	"sort"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

// BotCommand describes a BotCommand
type BotCommand struct {
	Function      func(*Bot, *kbchat.SubscriptionMessage, []string) error
	AdminRequired bool
}

// IsAdmin returns true if the user passed in is an admin
func (bot *Bot) IsAdmin(user string) bool {
	for _, n := range bot.admins {
		if user == n {
			return true
		}
	}
	return false
}

// IsFromAdmin returns true if the message is from an admin user, false otherwise
func (bot *Bot) IsFromAdmin(msg *kbchat.SubscriptionMessage) bool {
	return bot.IsAdmin(msg.Message.Sender.Username)
}

func (bot *Bot) interp(msg *kbchat.SubscriptionMessage, message string) error {
	bot.In.Lock()
	defer bot.In.Unlock()
	oneOnOne := true

	// use a tokenizer so that quotes and things are handled right
	args, err := shellquote.Split(message)
	if err != nil {
		return err
	}

	bot.Debug("Tokenizer found: %v", strings.Join(args, ","))

	// are you talking to me?
	if !strings.Contains(msg.Message.Channel.Name, ",") {
		if len(args) != 0 ||
			args[1] != "@"+bot.API().GetUsername() {
			return nil
		}
		oneOnOne = false
		args = args[1:]
	}

	// dumb to do every time but it's just the beginning here
	keys := make([]string, 0)
	for k := range bot.commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if strings.HasPrefix(k, strings.ToLower(args[0])) {
			if bot.commands[k].AdminRequired == true && !bot.IsFromAdmin(msg) {
				continue
			}
			return bot.commands[k].Function(bot, msg, args)
		}
	}
	if oneOnOne {
		return bot.ReplyTo(msg, "Huh???")
	}

	return nil
}
