package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/cultbaus/bot/internal/config"
)

type Gpt struct {
	url     string
	apiKey  string
	client  *http.Client
	history []message
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type prompt struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
}

func New(url string) *Gpt {
	g := &Gpt{
		url:     url,
		client:  &http.Client{},
		apiKey:  config.GetOpenAiKey(),
		history: []message{},
	}
	go g.cleanup()
	return g
}

func (g *Gpt) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			refMsg := m.ReferencedMessage
			if refMsg != nil && refMsg.Author.ID == s.State.User.ID && strings.Contains(refMsg.Content, "vxtwitter") {
				continue
			}
			response, err := g.send(m.Content)
			if err != nil {
				log.Println(err)
			}
			if _, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference()); err != nil {
				log.Println(err)
			}
			break
		}
	}
}

func (g *Gpt) send(p string) (string, error) {
	g.history = append(g.history, message{
		Role:    "user",
		Content: p,
	})
	pr := prompt{
		Model:    "gpt-3.5-turbo",
		Messages: g.history,
	}
	bs, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", g.url, bytes.NewBuffer(bs))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("internal server error: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var jsonResponse chatResponse
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return "", err
	}
	if len(jsonResponse.Choices) == 0 {
		return "", fmt.Errorf("did not get a response")
	}
	if len(jsonResponse.Choices[0].Message.Content) > config.DiscordCharacterLimit {
		return "Sorry, that message is too long", nil
	}
	g.history = append(g.history, message{
		Role:    "assistant",
		Content: jsonResponse.Choices[0].Message.Content,
	})
	return jsonResponse.Choices[0].Message.Content, nil
}

func (g *Gpt) cleanup() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		n := len(g.history) / 2
		g.history = g.history[n:]
	}
}
