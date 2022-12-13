package main

import (
	"github.com/Vallghall/VallTron/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

var (
	token string
)

func init() {
	token = "Bot " + token
}

func main() {
	defer utils.Logger.Sync()
	if token == "" {
		utils.Sugar.Fatalln("token is not provided")
	}

	dg, err := discordgo.New(token)
	if err != nil {
		utils.Sugar.Fatalf("Discord session creation failure: %v\n", err)
	}

	dg.AddHandler(pingPong)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		utils.Sugar.Fatalf("ws connection error: %v\n", err)
	}

	utils.Sugar.Infoln("Bot started")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	_ = dg.Close()
}

func pingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			utils.Sugar.Errorf("message send error: %v\n", err)
		}
	}
}
