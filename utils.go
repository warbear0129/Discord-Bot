package main

import (
	"os/exec"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
	"fmt"
	"log"
)

func getUserID(target string, serverID string, s *discordgo.Session) (string) {
	log.Printf("**** Finding : %s ... ****", target)
	members, _ := s.GuildMembers(serverID, 0, 100)

	for _, member := range members {
		if member.User.Username == target {
			log.Printf("**** Found user @ %s ****", member.User.ID)
			return fmt.Sprintf("<@%s>", member.User.ID)
		}
	}
	log.Printf("**** Member not found ****")
	return ""
}

func getRandomUserID(serverID string, s *discordgo.Session) (string) {
	rand.Seed(time.Now().Unix())
	members, _ := s.GuildMembers(serverID, 0, 100)
	return fmt.Sprintf("<@%s>", members[rand.Intn(len(members))].User.ID)
}

func getYoutubeTitle(url interface{}) (title []byte) {
	youtubedl := exec.Command("youtube-dl", "-e", url.(string))
	title, _ = youtubedl.Output()
	return
}
