package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"log"
	"os/exec"
)

var (
	ttsCooldown = make(map[string]int)
	players = make(map[string]*musicPlayer)
	faggot = make(map[string]string)
)

func whoisafaggot(s *discordgo.Session, m *discordgo.Message) {
	serverID := getServerID(s, m)

	if faggot[serverID] == "" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is a faggot!", getRandomUserID(serverID, s)))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is a faggot!", faggot))
	}
}

func thisguyisafaggot(s *discordgo.Session, m *discordgo.Message) {
	serverID := getServerID(s, m)

	params := getParams(m)
	if params == "" {
		return
	}

	if m.Author.ID == myID {
		faggot[serverID] = params
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Miku agrees, %s is a faggot", faggot[serverID]))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Miku disagrees, you're the faggot <@%s>!", m.Author.ID))
	}
}

func join(s *discordgo.Session, m *discordgo.Message) {
	serverID := getServerID(s, m)
	params := getParams(m)

	if players[serverID] == nil {
		players[serverID] = newMusicSession(params, serverID, s)
	}
}

func sing(s *discordgo.Session, m *discordgo.Message) {
	serverID := getServerID(s, m)
	params  := getParams(m)

	if players[serverID] == nil {
		players[serverID] = newMusicSession("", serverID, s)
	}

	go players[serverID].start(params)
}

func skip(s *discordgo.Session, m *discordgo.Message) {
	serverID := getServerID(s, m)

	if players[serverID] != nil {
		players[serverID].skip = true
	}
}

func stop(s *discordgo.Session, m *discordgo.Message) {
	serverID := getServerID(s, m)

	if players[serverID] != nil {
		players[serverID].stop()
	}
}

func ping(s *discordgo.Session, m *discordgo.Message) {
	params := getParams(m)

	if params == "" {
		params = "discord.gg"
	}

	ping := exec.Command("ping", "-c", "4", params)
	stdout, err := ping.Output()
	if err != nil {
		log.Println(err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", stdout))
}

func run(s *discordgo.Session, m *discordgo.Message) {
	if m.Author.ID != myID {
		s.ChannelMessageSend(m.ChannelID, "I only listen to my husband ;)")
		return
	}

	cmd, params := getParamsAll(m)

	if cmd == "" {
		return
	}

	run := exec.Command(cmd, params...)
	stdout, err := run.Output()
        if err != nil {
                s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", err))
                return
        }
        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", stdout))
}
