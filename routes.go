package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"log"
	"os/exec"
)

var (
	players = make(map[string]*musicPlayer)
)

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

func help(s *discordgo.Session, m *discordgo.Message) {
	msg :=  "Hello! I am Hatsune Miku and I am made by Max ;)\n"+
		"Here are the commands I have:\n\n" +
		"`help .............................. Show this help message`\n" +
		"`join <voice-channel> .............. Join a voice channel`\n" +
		"`play <youtube-url> ................ Add a song from YouTube to playlist`\n" +
		"`skip .............................. Skip a song in the playlist`\n" +
		"`stop .............................. Stop the entire playlist`\n" +
		"`ping .............................. Test the ping to your server ;)`\n"
	s.ChannelMessageSend(m.ChannelID,msg)
}
