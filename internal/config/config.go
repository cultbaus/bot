package config

import (
	"fmt"
	"os"
)

type DiscordUserID = string

const (
	DiscordCharacterLimit int           = 2000
	Phasedruid            DiscordUserID = "359008637139812373"
)

func GetToken() string {
	return readEnv("TOKEN")
}

func GetOpenAiKey() string {
	return readEnv("OPEN_AI_KEY")
}

func readEnv(ev string) string {
	if value := os.Getenv(ev); value != "" {
		return value
	}
	panic(fmt.Sprintf("config: %s is not set", ev))
}
