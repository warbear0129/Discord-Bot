package main

import (
        "log"
	"fmt"
        "encoding/binary"
	"strings"
	"sync"
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
	sendpcm         bool
        recv            chan *discordgo.Packet
	volume		float32
	ffmpeg          *exec.Cmd
        youtubedl       *exec.Cmd
	mutex		sync.Mutex
	buffer		[]int16
}

var (
	enc, _ = gopus.NewEncoder(sampleRate, channels, gopus.Audio)

)

func (mp *musicPlayer) play(url string) {
	mp.sendMessage(fmt.Sprintf("Now playing - **%s**", getYoutubeTitle(url)))

	// start youtube-dl to "download" the audio
	mp.youtubedl = exec.Command("youtube-dl", url, "-q", "-o", "-")
	youtubedlStdout, err := mp.youtubedl.StdoutPipe()
	if err != nil {
		log.Println("***** youtube-dl stdout error *****")
		log.Println(err)
	}

	// ffmpeg to pass the audio from youtube-dl to Discord
	mp.ffmpeg = exec.Command("ffmpeg", "-i", "-", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	mp.ffmpeg.Stdin = youtubedlStdout

	ffmpegStdout, err := mp.ffmpeg.StdoutPipe()
	if err != nil {
		log.Printf("***** ffmpeg stdout error: %s *****", err)
	}

	mp.youtubedl.Start()
	mp.ffmpeg.Start()

	mp.voice.Speaking(true)

	for {
		err = binary.Read(ffmpegStdout, binary.LittleEndian, mp.buffer)
		if err != nil {
			log.Printf("**** Error reading from stdout: %s ****", err)
			break
		}

		if mp.playing == false || mp.skip == true {
			break
		}

		mp.pcmChannel <- mp.buffer
	}

	mp.voice.Speaking(false)
}

func (mp *musicPlayer) run() {
	mp.playing = true

	for mp.playing {
		mp.skip = false

		url := mp.queue.Dequeue()
		if url == nil {
			break
		}

		mp.play(url.(string))
	}

	mp.playing = false
	mp.stop()
}

func newMusicSession(target string, serverID string, s *discordgo.Session) (mp *musicPlayer) {
        mp = &musicPlayer {
		session: s,
                encoder: enc,
                playing: false,
                skip: false,
                pcmChannel: make(chan []int16, 2),
                queue: lane.NewQueue(),
		buffer: make([]int16, 1920),
        }

	log.Printf("***** Finding channel : %s ... *****", target)
	channels, _ := s.GuildChannels(serverID)
	channelID := channels[1].ID

	for _, channel := range channels {
		if channel.Name == target {
			log.Printf("**** Channel found @ %s ****", channel.ID)
			channelID = channel.ID
			break
		}
	}

	log.Printf("**** Joining channel %s ****", channelID)
	mp.voice, _ = s.ChannelVoiceJoin(serverID, channelID, false, false)

	return mp
}

func (mp *musicPlayer) start(url string) {
	// if no url is entered, do not run
	if url == "" {
		return
	}

	// if url is not a youtube url, do not run
	if !strings.Contains(url, "https://www.youtube.com/") {
		mp.sendMessage("I can only play YouTube videos for now :(")
		return
	}

	// url is valid, add it to the queue
	mp.queue.Enqueue(url)

	// don't do anything else if it is already playing
	if mp.playing {
		mp.sendMessage(fmt.Sprintf("Added to queue - **%s**", getYoutubeTitle(url)))
		return
	}

	// if it is not playing already, play the song
	go SendPCM(mp, mp.pcmChannel)
	go mp.run()
}

func (mp *musicPlayer) stop() {
	if mp.playing {
		mp.playing = false
	} else {
		close(mp.pcmChannel)
		mp.voice.Close()
		mp.voice.Disconnect()
		delete(players, mp.voice.GuildID)
	}
}

func (mp *musicPlayer) sendMessage(msg string) {
	mp.session.ChannelMessageSend(mp.voice.GuildID, msg)
}
