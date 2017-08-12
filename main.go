package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

const (
	myID      = "152424821924298752"                                              // change to your discord user ID
	myChannel = "280009528966250498"                                              // change to your server's ID
	prefix    = "miku"                                                            // change to your preferred prefix
	token     = "Bot MzIwMTI3NzM0NTEwNzgwNDE5.DFYb3A.Xhcb3oWrcfxSk3N8Pq1XRzU1tMY" // replace with your bot's token
)

var (
	r = newRouter()
)

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// if user is a bot, ignore its message
	if m.Author.Bot {
		return
	}

	// ignore any messages less than the length of our prefix
	if len(m.Content) <= len(prefix) {
		return
	}

	// ignore any messages that does not start with our prefix
	if m.Content[:4] != prefix {
		return
	}

	// call the corresponding route
	r.getRoute(s, m.Message)
}

func main() {
	log.Printf("----- Logging in with token : %s", token)

	discord, err := discordgo.New(token)
	if err != nil {
		log.Printf("----- Error logging in: %s", err)
		os.Exit(0)
	}

	err = discord.Open()
	if err != nil {
		log.Printf("----- Error opening Discord: %s", err)
		os.Exit(0)
	}

	discord.AddHandler(messageCreateHandler)

	r = newRouter()

	go r.addRoute("ping", ping)
	go r.addRoute("run", run)
	go r.addRoute("help", help)

	c := make(chan os.Signal, 1)
	done := make(chan bool)

	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Println("--- Handling CTRL-C ")
			discord.Logout()
			discord.Close()
			done <- true
		}
	}()
	<-done
}
