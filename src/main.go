package main

import (
	"VK_tg_bot/commands"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	smileyFaces := []string{"😔", "😪", "😬", "🙄", "🤥", "😵", "‍💫", "🤕", "🧐", "🫤", "😟", "🥺", "🥹", "😦", "😧", "☹", "😮", "😫", "😢", "😱", "😖", "😣", "😞", "😵", "🤨"}
	rand.Seed(time.Now().UnixNano())

	db, err := sql.Open("postgres", "postgres://postgres:postgres@db/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		for range time.Tick(1 * time.Minute) {
			commands.DeleteExpiredServices(db)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			randomIndex := rand.Intn(len(smileyFaces))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, smileyFaces[randomIndex])
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else {
			commands.HandleCommand(db, bot, update)
		}
	}
}
