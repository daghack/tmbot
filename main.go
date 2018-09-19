package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BotToken string `split_words:"true" required:"true"`
	SkillsFile string `split_words:"true" required:"true"`
}

var conf Config

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	err = envconfig.Process("tmbot", &conf)
	if err != nil {
		panic(err)
	}
}

func main() {
	dg, err := discordgo.New("Bot " + conf.BotToken)
	if err != nil {
		panic(err)
	}
	handler, err := NewMessageHandler(conf.SkillsFile)
	if err != nil {
		panic(err)
	}
	dg.AddHandler(handler.Handler)
	err = dg.Open()
	if err != nil {
		panic(err)
	}
	defer dg.Close()
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}
