package main

import (
	"github.com/Vallghall/VallTron/internal/handler"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	token   string
	pattern string = `^\$vt\s\w+$`
)

func init() {
	token = "Bot " + token
}

func main() {
	if token == "" {
		log.Fatalln("token is not provided")
	}

	dg, err := discordgo.New(token)
	if err != nil {
		log.Fatalf("Discord session creation failure: %v\n", err)
	}

	h, err := handler.New(pattern)
	if err != nil {
		log.Fatalln(err)
	}

	dg.AddHandler(h.Handle)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatalf("ws connection error: %v\n", err)
	}

	log.Println("Bot started")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	_ = dg.Close()
}
