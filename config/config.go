package config

import (
	"log"

	goenv "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Environment struct {
	Telegram struct {
		BotToken string `env:"TELEGRAM_BOT_TOKEN"`
	}

	Broker struct {
		Provider string `env:"BROKER_PROVIDER"`

		RabbitMQ struct {
			URL string `env:"BROKER_RABBITMQ_URL"`
		}
	}

	MongoDB struct {
		URI string `env:"MONGO_URI"`
	}

	ServiceName string `env:"SERVICE_NAME"`
	Port        string `env:"PORT"`

	Extras goenv.EnvSet
}

var Env *Environment

func init() {
	Env = &Environment{}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("problem loading env: %s", err.Error())
	}

	es, err := goenv.UnmarshalFromEnviron(Env)
	if err != nil {
		log.Fatalf("problem loading marshaling env: %s", err.Error())
	}

	Env.Extras = es
}
