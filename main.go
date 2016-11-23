package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"flag"
	"os"
	"bufio"
)

const (
	myID = "152424821924298752"		// change to your discord user ID
	myChannel = "180240931893673987"	// change to your server's ID
	prefix = "miku"				// change to your preferred prefix
)

var (
	status	string
	token	string
	r	*router
)

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// if user is a bot, ignore his message
	if m.Author.Bot {
		return
	}

	// filter out any messages less than 4 characters otherwise we will get IndexOutOfRange
	if len(m.Content) <= len(prefix) {
		return
	}

	// process any message starting with "miku"
	if m.Content[:4] != prefix {
		return
	}

	// call the corresponding method
	r.getRoute(s, m.Message)
}

func init() {
	// parse flags
	tokenPtr := flag.String("token", "", "Discord Bot Token")
	statusPtr := flag.String("status", "War Thunder", "Discord Bot Status")
	flag.Parse()

	// exit if no token entered by user
	token = *tokenPtr
	if token == "" {
		log.Println("***** No token entered, use the -token flag *****")
		os.Exit(0)
	}

	// set status according to flag
	status = *statusPtr
}

func main() {
	log.Printf("----- Logging in with token : %s", token)

	discord, err := discordgo.New(token)
	if err != nil {
		log.Printf("----- Error logging in: %s -----", err)
		os.Exit(0)
	}

	err = discord.Open()
	if err != nil {
		log.Printf("----- Error opening Discord: %s -----", err)
		os.Exit(0)
	}

	log.Println("----- Adding handlers -----")
	discord.AddHandler(messageCreateHandler)

	log.Printf("----- Setting status to: %s -----", status)
	discord.UpdateStatus(0, status)

	log.Println("----- Setting up router -----")
	r = newRouter()
	go r.addRoute("thisguyisafaggot", thisguyisafaggot)
	go r.addRoute("whoisafaggot", whoisafaggot)
	go r.addRoute("join", join)
	go r.addRoute("sing", sing)
	go r.addRoute("skip", skip)
	go r.addRoute("stop", stop)
	go r.addRoute("ping", ping)
	go r.addRoute("run", run)

	log.Println("----- Starting scanner -----")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		discord.ChannelMessageSend(myChannel, scanner.Text())
	}
}
