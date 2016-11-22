package main

import (
        "log"
	"fmt"
	"sync"
        "encoding/binary"
	"strings"
        "os/exec"
        "gopkg.in/oleiade/lane.v1"
        "github.com/bwmarrin/discordgo"
	"github.com/layeh/gopus"
)

type musicPlayer struct {
	voice		*discordgo.VoiceConnection
	session		*discordgo.Session
        encoder         *gopus.Encoder
        queue           *lane.Queue
	pcmChannel	chan []int16
        playing         bool
        skip            bool
}

var (
	ffmpeg		*exec.Cmd
	youtubedl	*exec.Cmd
	sendpcm		bool
	recv		chan *discordgo.Packet
	mu		sync.Mutex
)

func (mp *musicPlayer) play(url string) {
	youtubedl = exec.Command("youtube-dl", url, "-q", "-o", "-")
	youtubedlStdout, err := youtubedl.StdoutPipe()
	if err != nil {
		log.Println("***** youtube-dl stdout error *****")
		log.Println(err)
	}

	ffmpeg = exec.Command("ffmpeg", "-i", "-", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1", "-af", "0.5")
	ffmpeg.Stdin = youtubedlStdout
	ffmpegStdout, err := ffmpeg.StdoutPipe()
	if err != nil {
		log.Println("***** ffmpeg stdout error *****")
		log.Println(err)
	}

	youtubedl.Start()
	ffmpeg.Start()

	audioBuffer := make([]int16, 1920)

	mp.voice.Speaking(true)
	defer mp.voice.Speaking(false)

	for {
		err = binary.Read(ffmpegStdout, binary.LittleEndian, &audioBuffer)
		if err != nil {
			log.Println("**** Error reading from stdout ****")
			log.Println(err)
			break
		}

		if mp.playing == false || mp.skip == true {
			break
		}
		mp.pcmChannel <- audioBuffer
	}
}

func (mp *musicPlayer) run() {
	mp.playing = true
	for mp.playing {
		mp.skip = false
		url := mp.queue.Dequeue()
		if url == nil {
			break
		}
		mp.sendMessage(fmt.Sprintf("Now playing - **%s**", getYoutubeTitle(url)))
		mp.play(url.(string))
	}
	mp.exit()
}

func newMusicSession(target string, sID string, s *discordgo.Session) (mp *musicPlayer) {
        enc, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
        if err != nil {
                log.Println("**** Error creating encoder ****")
                log.Println(err)
        }
	enc.SetBitrate(gopus.BitrateMaximum)

        mp = &musicPlayer {
		session: s,
                encoder: enc,
                playing: false,
                skip: false,
                pcmChannel: make(chan []int16, 2),
                queue: lane.NewQueue(),
        }


	log.Printf("***** Finding channel : %s ... *****", target)
	channels, _ := s.GuildChannels(sID)
	channelID := channels[1].ID

	for _, channel := range channels {
		if channel.Name == target {
			log.Printf("**** Channel found @ %s ****", channel.ID)
			channelID = channel.ID
			break
		}
	}
	log.Printf("**** Joining channel %s ****", channelID)
	mp.voice, _ = s.ChannelVoiceJoin(sID, channelID, false, false)
	return mp
}

func (mp *musicPlayer) start(url string) {
	if url == "" {
		return
	}

	if !strings.Contains(url, "https://www.youtube.com/") {
		return
	}

	if mp.playing {
		mp.queue.Enqueue(url)
		mp.sendMessage(fmt.Sprintf("Added to queue - **%s**", getYoutubeTitle(url)))
		return
	}

	go SendPCM(mp, mp.pcmChannel)
	mp.queue.Enqueue(url)
	go mp.run()
}

func (mp *musicPlayer) stop() {
	if mp.playing {
		mp.playing = false
	} else {
		mp.exit()
	}
}

func (mp *musicPlayer) exit() {
        close(mp.pcmChannel)
        mp.voice.Close()
        mp.voice.Disconnect()
        delete(players, mp.voice.GuildID)
}

func (mp *musicPlayer) sendMessage(msg string) {
	mp.session.ChannelMessageSend(mp.voice.GuildID, msg)
}
