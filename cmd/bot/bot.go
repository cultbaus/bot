package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/cultbaus/bot/internal/config"
	"github.com/cultbaus/bot/internal/gpt"
	"github.com/cultbaus/bot/internal/tweet"
)

var (
	chat    *gpt.Gpt
	session *discordgo.Session
)

func init() {
	var err error
	session, err = discordgo.New("Bot " + config.GetToken())
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	chat = gpt.New("https://api.openai.com/v1/chat/completions")
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	session.AddHandler(chat.Handler)
	session.AddHandler(tweet.Handler)

	if err := session.Open(); err != nil {
		panic(err)
	}
	defer session.Close()

	<-stop
}
