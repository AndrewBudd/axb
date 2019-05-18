// Package axb is derived almost entirely from an example provided by @xgess,
// many thanks!
package axb

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

type interpfunc func(*bot, string, string) error

type bot struct {
	chatAPI       *kbchat.API
	debugTeamName string
	interp        interpfunc
}

func (b *bot) API() *kbchat.API {
	return b.chatAPI
}

func (b *bot) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)

	if err := b.API().SendMessageByTeamName(b.debugTeamName, msg, nil); err != nil {
		fmt.Printf("Error sending message; %s", err.Error())
	}
}

func (b *bot) SendToUser(user string, message string) error {
	tlfName := fmt.Sprintf("%s,%s", user, b.API().GetUsername())
	err := b.API().SendMessageByTlfName(tlfName, message)
	if err != nil {
		b.Debug(err.Error())
	}
	return err
}

func NewBot(debugTeamName string, keybaseLocation string, interp interpfunc) (*bot, error) {
	chatAPI, err := kbchat.Start(kbchat.RunOptions{KeybaseLocation: keybaseLocation})

	if err != nil {
		return nil, err

	}

	b := bot{
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

			err = b.interp(&b, msg.Message.Sender.Username, msg.Message.Content.Text.Body)
			if err != nil {
				b.Debug("Error calling interp: %s", err.Error())
			}
		}
	}()

	return &b, nil
}
