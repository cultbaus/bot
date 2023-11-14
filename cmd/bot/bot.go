package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const Phasedruid = "359008637139812373"

var (
	session    *discordgo.Session
	tweetRegex *regexp.Regexp
)

func init() {
	var err error

	session, err = discordgo.New("Bot " + getToken())
	if err != nil {
		log.Fatal(err)
	}

	tweetRegex = regexp.MustCompile(`(https://(?:twitter\.com|x\.com)|(http://(?:twitter\.com|x\.com))/\S+)`)
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	if err := session.Open(); err != nil {
		panic(err)
	}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if !tweetRegex.MatchString(m.Content) {
			return
		}

		msg := m.Author.Mention()

		if m.Author.ID == Phasedruid {
			msg += ", I think you should stay off twitter for a while"
		} else {
			msg += ", embeds: " + tweetRegex.ReplaceAllString(m.Content, "https://vxtwitter.com")
		}

		if _, err := s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference()); err != nil {
			log.Println(err)
		}
		if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
			log.Println(err)
		}
	})

	<-stop
}

func getToken() string {
	return readEnv("TOKEN")
}

func readEnv(ev string) string {
	if value := os.Getenv(ev); value != "" {
		return value
	}
	panic(fmt.Sprintf("config: %s is not set", ev))
}
