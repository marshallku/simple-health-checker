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

type NotificationParams struct {
	Title       string
	Description string
	Color       string
	Fields      map[string]string
	Footer      string
}

func SendNotification(cfg *config.Config, params NotificationParams) {
	if cfg.WebhookURL == "" {
		fmt.Println("Webhook URL is not set")
		return
	}

	title := params.Title
	if title == "" {
		title = "Health check failed"
	}

	sendDiscordNotification(cfg.WebhookURL, NotificationParams{
		Title:       title,
		Description: params.Description,
		Color:       params.Color,
		Fields:      params.Fields,
		Footer:      time.Now().Format(time.RFC3339),
	})
}

func sendDiscordNotification(webhookURI string, params NotificationParams) {
	colorInt, _ := strconv.Atoi(params.Color)

	discordFields := make([]DiscordField, 0, len(params.Fields))
	for name, value := range params.Fields {
		discordFields = append(discordFields, DiscordField{
			Name:  name,
			Value: value,
		})
	}

	payload := DiscordPayload{
		Embeds: []DiscordEmbed{
			{
				Type:        "rich",
				Title:       params.Title,
				Description: params.Description,
				Color:       colorInt,
				Fields:      discordFields,
			},
		},
	}

	if params.Footer != "" {
		payload.Embeds[0].Footer = DiscordFooter{
			Text: params.Footer,
		}
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
