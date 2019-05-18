// Package axb is derived almost entirely from an example provided by @xgess,
// many thanks!
package axb

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

type interpfunc func(*Bot, *kbchat.SubscriptionMessage, string) error

type Bot struct {
	chatAPI       *kbchat.API
	debugTeamName string
	interp        interpfunc
}

func (b *Bot) API() *kbchat.API {
	return b.chatAPI
}

func (b *Bot) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)

	if err := b.API().SendMessageByTeamName(b.debugTeamName, msg, nil); err != nil {
		fmt.Printf("Error sending message; %s", err.Error())
	}
}

func (b *Bot) ReplyTo(msg *kbchat.SubscriptionMessage, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	err := b.API().SendMessage(msg.Message.Channel, message)
	if err != nil {
		b.Debug(err.Error())
	}
	return err
}

func (b *Bot) SendToUser(user string, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	tlfName := fmt.Sprintf("%s,%s", user, b.API().GetUsername())
	err := b.API().SendMessageByTlfName(tlfName, msg)
	if err != nil {
		b.Debug(err.Error())
	}
	return err
}

func NewBot(debugTeamName string, keybaseLocation string, interp interpfunc) (*Bot, error) {
	chatAPI, err := kbchat.Start(kbchat.RunOptions{KeybaseLocation: keybaseLocation})

	if err != nil {
		return nil, err

	}

	b := Bot{
		chatAPI:       chatAPI,
		debugTeamName: debugTeamName,
		interp:        interp,
	}

	if err = chatAPI.SendMessageByTeamName(debugTeamName, "Starting up...", nil); err != nil {
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

			if msg.Message.Content.Type != "text" {
				continue
			}

			if msg.Message.Sender.Username == b.API().GetUsername() {
				continue
			}

			b.Debug("Received message, channel: %s, username: %s, message: %s, type: %s", msg.Message.Channel.Name, msg.Message.Sender.Username, msg.Message.Content.Text.Body, msg.Message.Channel.MembersType)

			err = b.interp(&b, &msg, msg.Message.Content.Text.Body)
			if err != nil {
				b.Debug("Error calling interp: %s", err.Error())
			}
		}
	}()

	return &b, nil
}
