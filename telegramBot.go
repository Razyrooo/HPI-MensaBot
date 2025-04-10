package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func sendTelegramMessage(message string, chatID string, botToken string) error {
	// Construct API URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Prepare data for http.PostForm
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)
	data.Set("parse_mode", "HTML")

	// Send the POST request
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %w", err)
	}
	// Ensure the response body is closed
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Try to read the error response body from Telegram
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			// If reading body fails, return the original status error
			return fmt.Errorf("telegram API error: %s (failed to read response body: %v)", resp.Status, readErr)
		}
		return fmt.Errorf("telegram API error: %s, response: %s", resp.Status, string(bodyBytes))
	}
	return nil
}
