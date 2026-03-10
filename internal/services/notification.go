package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/byt3roof/odoo-backup/internal/conf"
)

type Embed struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Color       int                    `json:"color"`
	Timestamp   string                 `json:"timestamp"`
	Footer      map[string]interface{} `json:"footer"`
}

type DiscordPayload struct {
	Embeds []Embed `json:"embeds"`
}

var (
	successColor = 0x00FF00
	failureColor = 0xFF0000
)

func discordNotification(url string, payload DiscordPayload, succeed bool) error {
	color := failureColor
	if succeed {
		color = successColor
	}

	for i := range payload.Embeds {
		payload.Embeds[i].Color = color
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func SendBackupNotification(cfg *conf.Config, instanceName string, sizeMB float64, err error, success bool) error {
	description := fmt.Sprintf("%s - %.2f MB", instanceName, sizeMB)

	if err != nil {
		description = fmt.Sprintf("%s - %.2f MB - Intance only stored on local \n\n```\n%v\n```", instanceName, sizeMB, err)
	}

	payload := DiscordPayload{
		Embeds: []Embed{
			{
				Title:       "Backup Summary",
				Description: description,
				Timestamp:   time.Now().Format(time.RFC3339),
				Footer: map[string]interface{}{
					"text": "Odoo backup service",
				},
			},
		},
	}

	return discordNotification(cfg.DiscordURL, payload, success)
}
