package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	//load environment variables from
	//only necessary if a .env file is used
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	//get meals from mensa api
	var t time.Time = time.Now()
	meals, prices := getMeals(t)
	var message string = fmt.Sprintf("<b>%s's menue</b>\n\n", t.Weekday())
	for i, meal := range meals {
		if prices[i] == 0.0 {
			break
		}
		message = message + fmt.Sprintf("<b><u>Angebot %d:</u></b>\n%s\n%.2f€\n\n", i+1, meal, prices[i])
	}

	//load credentials from environmental variables
	var chatID string = os.Getenv("CHATID")
	var botToken string = os.Getenv("BOT_TOKEN")
	if chatID == "" || botToken == "" {
		log.Fatal("ChatID or Bottoken missing")
	}

	sendTelegramMessage(message, chatID, botToken)
}
