// Package axb is derived almost entirely from an example provided by @xgess,
// many thanks!
package axb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

var builtInCommands = map[string]BotCommand{
	"help": {
		doHelp,
		false,
	},
	"shutdown": {
		doShutdown,
		true,
	},
	"printadmins": {
		doPrintAdmins,
		true,
	},
	"addadmin": {
		doAddAdmin,
		true,
	},
	"removeadmin": {
		doRemoveAdmin,
		true,
	},
}

// Bot is an internal data structure used by the Bot
type Bot struct {
	In            sync.Mutex
	Out           sync.Mutex
	chatAPI       *kbchat.API
	debugTeamName string
	admins        []string
	commands      map[string]BotCommand
}

// API returns the underlying kbchat API struct
func (bot *Bot) API() *kbchat.API {
	return bot.chatAPI
}

// Debug prints to stderr, as well as the debug team
func (bot *Bot) Debug(format string, args ...interface{}) {
	bot.Out.Lock()
	defer bot.Out.Unlock()
	msg := fmt.Sprintf(format, args...)
	fmt.Errorf(msg)
	if _, err := bot.API().SendMessageByTeamName(bot.debugTeamName, nil, msg); err != nil {
		fmt.Errorf("Error sending message; %s", err.Error())
	}
}

// ReplyTo sends a message to the origin of a message
func (bot *Bot) ReplyTo(msg *kbchat.SubscriptionMessage, format string, args ...interface{}) error {
	bot.Out.Lock()
	defer bot.Out.Unlock()
	message := fmt.Sprintf(format, args...)
	_, err := bot.API().SendMessage(msg.Message.Channel, message)
	if err != nil {
		bot.Debug(err.Error())
	}
	return err
}

// SendToUser sends a message to a particular user
func (bot *Bot) SendToUser(user string, format string, args ...interface{}) error {
	bot.Out.Lock()
	defer bot.Out.Unlock()
	msg := fmt.Sprintf(format, args...)
	tlfName := fmt.Sprintf("%s,%s", user, bot.API().GetUsername())
	_, err := bot.API().SendMessageByTlfName(tlfName, msg)
	if err != nil {
		bot.Debug(err.Error())
	}
	return err
}

// NewBot is used to initialize and connect a new bot
func NewBot(debugTeamName string, keybaseLocation string, commands map[string]BotCommand, admins []string) (*Bot, error) {
	chatAPI, err := kbchat.Start(kbchat.RunOptions{KeybaseLocation: keybaseLocation})
	targetCommands := make(map[string]BotCommand)

	if err != nil {
		return nil, err

	}

	for k, v := range builtInCommands {
		targetCommands[k] = v
	}

	// allow the user to clobber the built-in commands
	for k, v := range commands {
		targetCommands[k] = v
	}

	b := Bot{
		chatAPI:       chatAPI,
		debugTeamName: debugTeamName,
		commands:      targetCommands,
		admins:        admins,
	}

	if _, err = chatAPI.SendMessageByTeamName(debugTeamName, nil, "Starting up..."); err != nil {
		return nil, err
	}

	subscription, err := chatAPI.ListenForNewTextMessages()

	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg, err := subscription.Read()
			switch err := err.(type) {
			case nil:
			case *json.SyntaxError:
				b.Debug("Error reading message (fatal): %s", err.Error())
				os.Exit(1)
			default:
				b.Debug("Error reading message (nonfatal): %s", err.Error())
				continue
			}

			if msg.Message.Content.TypeName != "text" {
				continue
			}

			if msg.Message.Sender.Username == b.API().GetUsername() {
				continue
			}

			// b.Debug("Received message, channel: %s, username: %s, message: %s, type: %s", msg.Message.Channel.Name, msg.Message.Sender.Username, msg.Message.Content.Text.Body, msg.Message.Channel.MembersType)

			err = b.interp(&msg, msg.Message.Content.Text.Body)
			if err != nil {
				b.Debug("Error calling interp: %s", err.Error())
			}
		}
	}()

	return &b, nil
}
