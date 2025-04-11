package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_time := time.Now()
	if (_time.Weekday()) == time.Saturday || _time.Weekday() == time.Sunday {
		fmt.Println("Nothing to do at the weekend.")
		return
	}

	//load environment variables from
	//only necessary if a .env file is used
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	//get meals from mensa api
	meals, prices := getMeals(_time)
	message := fmt.Sprintf("<b>%s's menue</b>\n\n", _time.Weekday())
	for i, meal := range meals {
		if prices[i] == 0.0 {
			break
		}
		message = message + fmt.Sprintf("<b><u>Angebot %d:</u></b>\n%s\n%.2f€\n\n", i+1, meal, prices[i])
	}

	//load credentials from environmental variables
	chatID := os.Getenv("CHATID")
	botToken := os.Getenv("BOT_TOKEN")
	if chatID == "" || botToken == "" {
		log.Fatal("ChatID or Bottoken missing")
	}

	err_t := sendTelegramMessage(message, chatID, botToken)
	if err_t != nil {
		log.Fatal("failed to send Telegram Message", err_t)
	}
}
