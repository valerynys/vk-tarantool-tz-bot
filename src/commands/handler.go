package commands

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

func HandleCommand(db *sql.DB, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	userID := update.Message.From.ID

	var expirationTime time.Time

	switch update.Message.Command() {
	case "start":
		msg.Text = "Hello, " + update.Message.From.FirstName + "!"
	case "help":
		msg.Text = "I understand /set, /get, /delete commands"
	case "set":
		args := update.Message.CommandArguments()

		err := SetService(db, userID, args, expirationTime)
		if err != nil {
			msg.Text = err.Error()
			break
		}

		msg.Text = "Service added successfully"
	case "get":
		args := update.Message.CommandArguments()

		login, password, err := GetService(db, userID, args)
		if err != nil {
			msg.Text = err.Error()
			break
		}

		msg.Text = "login: " + login + "\npassword: " + password
	case "delete":
		args := update.Message.CommandArguments()
		err := DeleteService(db, userID, args)
		if err != nil {
			msg.Text = err.Error()
			break
		}
		msg.Text = "Service delete successfully"
	default:
		msg.Text = "I don't know that command"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
