package conf

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUri   string `env:"MONGO_URI,required"`
	Password   string `env:"PASSWORD,required"`
	Salt       string `env:"SALT,required"`
	AccessKey  string `env:"ACCESS_KEY,required"`
	SecretKey  string `env:"SECRET_KEY,required"`
	Bucket     string `env:"BUCKET,required"`
	S3Endpoint string `env:"S3_ENDPOINT"`
	DiscordURL string `env:"DISCORD_WEBHOOK"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	var missingVars []string

	t := reflect.TypeOf(*cfg)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")

		if strings.Contains(tag, "required") {
			envVar := strings.Split(tag, ",")[0]
			if strings.TrimSpace(os.Getenv(envVar)) == "" {
				missingVars = append(missingVars, envVar)
			}
		}

		envVar := strings.Split(tag, ",")[0]
		reflect.ValueOf(cfg).Elem().Field(i).SetString(os.Getenv(envVar))
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	return cfg, nil
}
