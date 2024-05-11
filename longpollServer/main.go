package main

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"longpollServer/authorization"
	"longpollServer/chat"
	"longpollServer/web"
)

func main() {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to load config")
	}

	if err = env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("Unable to parse config")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresName, cfg.PostgresPort)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database")
	}

	tokenAuth := jwtauth.New("HS256", []byte(cfg.Secret), nil)

	chatController := chat.NewController(db)
	authController := authorization.NewController(db, tokenAuth)

	longpollChat := web.NewLongpollChat(chatController, authController, tokenAuth)
	longpollChat.Start()
}

type Config struct {
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     int    `env:"POSTGRES_PORT"`
	PostgresName     string `env:"POSTGRES_NAME"`
	Secret           string `env:"SECRET"`
}
