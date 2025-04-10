package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

//telegram send messages: https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage

func main() {
	//loading environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var t time.Time = time.Now()
	meals, prices := getMeals(t)
	var message string = fmt.Sprintf("<b>%s's menue</b>\n\n", t.Weekday())
	for i, meal := range meals {
		if prices[i] == 0.0 {
			break
		}
		message = message + fmt.Sprintf("<b><u>Angebot %d:</u></b>\n%s\n%.2f€\n\n", i+1, meal, prices[i])
	}

	var chatID string = os.Getenv("CHATID")
	var botToken string = os.Getenv("BOT_TOKEN")
	if chatID == "" || botToken == "" {
		log.Fatal("ChatID or Bottoken missing")
	}
	sendTelegramMessage(message, chatID, botToken)
}
