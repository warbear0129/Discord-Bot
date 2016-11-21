package main

import (
	"github.com/bwmarrin/discordgo"
	"os/exec"
	"log"
	"fmt"
	"strings"
)

const (
	prefix = "miku"
	me = "152424821924298752"
)

var (
	cmd	*exec.Cmd
	faggot = make(map[string]string)
	players = make(map[string]*musicPlayer)
)

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author

	// if user is a bot, ignore his message
	if user.Bot {
		return
	}

	// filter out any messages less than 4 characters otherwise we will get IndexOutOfRange
	if len(m.Content) < len(prefix) {
		return
	}

	// process any message starting with "miku"
	if m.Content[:4] == "miku" {

		// some variables delcaration to store user input
		var method string
		var params string

		// filtering out "miku" since it isn't needed anymore
		content := m.Content[5:]

		// get the Server ID in which this handler is called
		channel, _ := s.Channel(m.ChannelID)
		serverID := channel.GuildID

		// save called method in method variable
		// if user input does not contain space, it must not contain any parameters
		// otherwise, save user paramets into params variable
		if !strings.Contains(content, " ") {
			method = content
		} else {
			method = strings.Split(content, " ")[0]
			params = content[len(method)+1:]
		}

		switch method {

		case "whoisafaggot":
			if faggot[serverID] == "" {
				faggot[serverID] = getRandomUserID(serverID, s)
				params = ""
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is a faggot!", faggot[serverID]))

		case "thisguyisafaggot":
			if user.ID == me {
				faggot[serverID] = params
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Miku agrees, %s is a faggot", faggot[serverID]))
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Miku disagrees, you're the faggot <@%s>!", user.ID))
			}

		case "join":
			if players[serverID] == nil {
				players[serverID] = newMusicSession(params, serverID, s)
			} else {
				s.ChannelMessageSend(m.ChannelID, "I am already in a voice channel!")
			}

		case "reboot":
			if user.ID == me {
				s.ChannelMessageSend(m.ChannelID, "Miku is rebooting ...")
				cmd = exec.Command("reboot")
				cmd.Run()
			} else {
				s.ChannelMessageSend(m.ChannelID, "You are not my husband ...!")
			}

		case "shutdown":
			if user.ID == me {
				s.ChannelMessageSend(m.ChannelID, "Miku is shutting down ...")
				cmd = exec.Command("sudo telinit 0")
				cmd.Run()
			} else {
				s.ChannelMessageSend(m.ChannelID, "You are not my husband ...!")
			}

		case "play":
			if players[serverID] != nil {
				go players[serverID].start(params)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s** added to playlist", getYoutubeTitle(params)))
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Put me in a channel first")

		case "skip":
			if players[serverID] != nil {
				players[serverID].skip = true
				s.ChannelMessageSend(m.ChannelID, "Skipping song ...")
				return
			}

		case "stop":
			if players[serverID] != nil {
				players[serverID].playing = false
				s.ChannelMessageSend(m.ChannelID, "Exiting ...")
				return
			}

		case "whoisyourhusband":
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s is my husband <3", me))

		case "doyoulovemax":
			s.ChannelMessageSend(m.ChannelID, "Yes! <3")
		default:
			return
		}
		log.Printf("----- %s -----", m.ChannelID)
		log.Printf("%s called    : %s", m.Author.Username, method)
		log.Printf("With params  : %s", params)
	}
	log.Printf("%s said : %s", user.Username, m.Content)
}

func main() {
	log.Println("----- Starting bot -----")
	discord, err := discordgo.New("Bot MjQ4MjcxODQ4MDg3ODc5Njgz.Cw1fOg.gIvJixhDUCkgQthaPeja_LmJkS0")
	if err != nil {
		log.Println("----- Error logging in -----")
		log.Println(err)
	}
	log.Println("----- Login successful -----")
	log.Println("----- Opening Discord -----")

	discord.AddHandler(messageHandler)

	err = discord.Open()
	if err != nil {
		log.Println("----- Error opening Discord -----")
		log.Println(err)
	}

	lock := make(chan int)
	<-lock
}
