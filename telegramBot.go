package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

//telegram send messages: https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage

func sendTelegramMessage(message string, chatID string, botToken string) error {
	//construct API URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	//prepare data for http.PostForm
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)
	data.Set("parse_mode", "HTML")

	//send the POST request
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %w", err)
	}
	//ensure the response body is closed
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("failed to clode response body: %w", err)
		}
	}()

	//checks for errors if the response fails
	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("telegram API error: %s (failed to read response body: %v)", resp.Status, readErr)
		}
		return fmt.Errorf("telegram API error: %s, response: %s", resp.Status, string(bodyBytes))
	}
	return nil //response ok
}
