# Introduction
AXB is a lightweight framework which simplifies the process of building a keybase bot.  See the example below for how to get started.

## Example
```
import (
	"os"
	"github.com/AndrewBudd/axb"
)

var commands = map[string]axb.BotCommand{
	"hello": {
		doHello, // function that handles the hello command
		false,    // does this function require an admin user?
	},
}

var admins = []string{"andrewbudd"}

func doHello(bot *axb.Bot, msg *kbchat.SubscriptionMessage, args []string) error {
	return bot.ReplyTo(msg, "Hello yourself!")
}

func main() {
	axb.NewBot(os.Getenv("KEYBASE_DEBUG_TEAM"), os.Getenv("KEYBASE_LOCATION"), commands, admins)
	select {}
}
```

## built-in commands
The bot includes a number of built in commands, specifically
* help - prints all of the installed commands that you have access to 
* printadmins - prints all of the admin users
* addadmin - adds an admin user
* removeadmin - removes an admin user
* shutdown - shuts down the bot
