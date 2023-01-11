package handler

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strings"
)

type Command func(s *discordgo.Session, m *discordgo.MessageCreate) error

type CommandHandler struct {
	pattern *regexp.Regexp
	actions map[string]Command
}

func New(p string) (*CommandHandler, error) {
	pattern, err := regexp.Compile(p)
	if err != nil {
		return nil, err
	}

	return &CommandHandler{
		pattern: pattern,
		actions: map[string]Command{
			"grayer": handleGrayer,
		},
	}, nil
}

func (inst *CommandHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !inst.pattern.MatchString(m.Content) {
		return
	}

	indices := strings.SplitN(m.Content, " ", 3)
	if len(indices) < 2 {
		return
	}

	command := indices[1]
	action, ok := inst.actions[command]
	if !ok {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("action for command %s not found", command))
		if err != nil {
			log.Println("failed to send message")
		}
		return
	}

	err := action(s, m)
	if err != nil {
		log.Printf("error in action %s: %v", command, err)
	}
}
