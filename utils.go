package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getServerID(s *discordgo.Session, m *discordgo.Message) string {
	channel, _ := s.Channel(m.ChannelID)
	return channel.GuildID
}

func getParams(m *discordgo.Message) string {
	content := strings.Split(m.Content, " ")

	if len(content) < 3 {
		return ""
	}

	return content[2]
}

func getParamsAll(m *discordgo.Message) (string, []string) {
	content := strings.Split(m.Content, " ")

	if len(content) < 3 {
		return "", nil
	}

	return content[2], content[3:]
}

func getMethod(m *discordgo.Message) (method string) {
	return strings.Split(m.Content, " ")[1]
}

func getUserID(target string, serverID string, s *discordgo.Session) string {
	members, _ := s.GuildMembers(serverID, "0", 100)

	for _, member := range members {
		if member.User.Username == target {
			return fmt.Sprintf("<@%s>", member.User.ID)
		}
	}

	return ""
}

func getRandomUserID(serverID string, s *discordgo.Session) string {
	rand.Seed(time.Now().Unix())
	members, _ := s.GuildMembers(serverID, "0", 100)
	return fmt.Sprintf("<@%s>", members[rand.Intn(len(members))].User.ID)
}
