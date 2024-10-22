package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/marshallku/simple_health_checker/config"
)

type DiscordEmbed struct {
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Color       int            `json:"color,string"`
	Fields      []DiscordField `json:"fields"`
	Footer      DiscordFooter  `json:"footer"`
}

type DiscordField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DiscordFooter struct {
	Text string `json:"text"`
}

type DiscordPayload struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

func SendNotification(cfg *config.Config, description, color string, fields map[string]string) {
	if cfg.WebhookURL == "" {
		fmt.Println("Webhook URL is not set")
		return
	}

	sendDiscordNotification(cfg.WebhookURL, "Health check failed", description, color, fields)
}

func sendDiscordNotification(webhookURI string, title string, description string, color string, fields map[string]string) {
	colorInt, _ := strconv.Atoi(color)

	discordFields := make([]DiscordField, 0, len(fields))
	for name, value := range fields {
		discordFields = append(discordFields, DiscordField{
			Name:  name,
			Value: value,
		})
	}

	payload := DiscordPayload{
		Embeds: []DiscordEmbed{
			{
				Type:        "rich",
				Title:       title,
				Description: description,
				Color:       colorInt,
				Fields:      discordFields,
				Footer: DiscordFooter{
					Text: time.Now().Format(time.RFC3339),
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	resp, err := http.Post(webhookURI, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error sending Discord notification: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Discord API returned non-OK status: %d\n", resp.StatusCode)
	}
}
