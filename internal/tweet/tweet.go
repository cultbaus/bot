package tweet

import (
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"

	"github.com/cultbaus/bot/internal/config"
)

var tweetRegex = regexp.MustCompile(`(https://(?:twitter\.com|x\.com)|(http://(?:twitter\.com|x\.com))/\S+)`)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !tweetRegex.MatchString(m.Content) {
		return
	}

	msg := m.Author.Mention()

	if m.Author.ID == config.Phasedruid {
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
}
