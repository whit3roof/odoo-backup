package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/byt3roof/odoo-backup/internal/conf"
	"github.com/byt3roof/odoo-backup/internal/db/mongo"
	"github.com/byt3roof/odoo-backup/internal/services"
	"golang.org/x/crypto/argon2"
)

func main() {
	ctx := context.Background()

	cfg, err := conf.LoadConfig()
	if err != nil {
		panic(err)
	}

	conn, err := mongo.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Disconnect(ctx)

	results, err := mongo.FetchCollection(ctx, conn, "odoo", "backup_config", 10)
	if err != nil {
		panic(err)
	}

	key := argon2.IDKey([]byte(cfg.Password), []byte(cfg.Salt), 1, 164*10244, 4, 32)

	for _, doc := range results {
		encryptedPwd, ok := doc["db_password"].(string)
		if !ok {
			continue
		}

		plainPwd, err := services.DecryptText(encryptedPwd, key)
		if err != nil {
			continue
		}

		filePath, err := services.BackupOdoo(ctx, services.BackupConfig{
			Instance:   doc["instance"].(string),
			DomainURL:  doc["domain_url"].(string),
			DBPassword: plainPwd,
			DBName:     doc["db_name"].(string),
		}, "./backups")
		if err != nil {
			continue
		}

		s3Key := filepath.Base(filePath)

		if err := services.UploadToS3(ctx, filePath, s3Key); err != nil {
			fmt.Printf("failed to upload %s to S3: %v\n", filePath, err)
			continue
		}

		_ = os.Remove(filePath)
		fmt.Printf("successfully backed up and uploaded %s\n", s3Key)
	}
}
