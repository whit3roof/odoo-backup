package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type BackupConfig struct {
	Instance   string
	DomainURL  string
	DBPassword string
	DBName     string
}

func BackupOdoo(ctx context.Context, cfg BackupConfig, outDir string) (string, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()

		writer.WriteField("master_pwd", cfg.DBPassword)
		writer.WriteField("name", cfg.DBName)
		writer.WriteField("backup_format", "zip")
	}()

	url := fmt.Sprintf("https://%s/web/database/backup", cfg.DomainURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pr)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}

	fileName := fmt.Sprintf("%s/odoo_%s_%s.zip", outDir, cfg.Instance, time.Now().Format("2006-01-02"))
	f, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("backed up %s -> %s\n", cfg.Instance, fileName)
	return fileName, nil
}
