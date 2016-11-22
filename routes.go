package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
)

const (
	me = "152424821924298752"
	hupsoonheng = "180240931893673987"
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

	if m.Author.ID == me {
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

	if players[serverID] != nil {
		go players[serverID].start(params)
	}
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
