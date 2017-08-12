package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/bwmarrin/discordgo"
)

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
		s.ChannelMessageSend(m.ChannelID, "Don't speak to me please.")
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
	msg := "```help .............................. Show this help message\n" +
		"ping .............................. Test the ping to your server ;)```\n\n" +
		"You can view my source code at:\n"
	s.ChannelMessageSend(m.ChannelID, msg)
}
