package main

import (
	"io"
        "log"
	"fmt"
	"sync"
        "encoding/binary"
        "os/exec"
        "net/http"
        "gopkg.in/oleiade/lane.v1"
        "github.com/bwmarrin/discordgo"
	"github.com/hraban/opus"
	"github.com/layeh/gopus"
)

const (
	maxBytes = 3840 
	channels = 2
	sampleRate = 48000
)

type musicPlayer struct {
        encoder         *opus.Encoder
        voice	        *discordgo.VoiceConnection
        queue           *lane.Queue
	pcmChannel	chan []int16
        playing         bool
        skip            bool
        stop            bool
}

var (
	ffmpeg		*exec.Cmd
	youtubedl	*exec.Cmd
	sendpcm		bool
	recv		chan *discordgo.Packet
	mu		sync.Mutex
	players	= make(map[string]*musicPlayer)
)



func (mp *musicPlayer) play(url string) {
	mp.playing = true

	youtubedl = exec.Command("youtube-dl", url, "-q", "-o", "-")
	youtubedlStdout, err := youtubedl.StdoutPipe()
	if err != nil {
		log.Println("***** youtube-dl stdout error *****")
		log.Println(err)
	}

	ffmpeg = exec.Command("ffmpeg", "-i", "-", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
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
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			log.Println("**** Error reading from stdout ****")
			log.Println(err)
			break
		}

		if mp.stop == true || mp.skip == true {
			break
		}
		mp.pcmChannel <- audioBuffer
	}

	mp.playing = false
}

func (mp *musicPlayer) run() {
	if mp.playing == false {
		for {
			mp.skip = false
			url := mp.queue.Dequeue()
			if url == nil || mp.stop == true {
				break
			}

			mp.play(url.(string))
		}
		log.Println("***** No more songs in queue, closing pcm channel *****")
		close(mp.pcmChannel)
		mp.voice.Close()
		delete(players, mp.voice.GuildID)
		log.Println("***** Music Player cleanup finished *****")
	}
}

func initializeMp(url string, voice *discordgo.VoiceConnection) (mp *musicPlayer) {
        enc, err := opus.NewEncoder(sampleRate, channels, opus.APPLICATION_VOIP)
        if err != nil {
                log.Println("**** Error creating encoder ****")
                log.Println(err)
        }

	mp = &musicPlayer {
                voice: voice,
                encoder: enc,
                playing: false,
                skip: false,
                stop: false,
                pcmChannel: make(chan []int16, 2),
                queue: lane.NewQueue(),
	}
	go SendPCM(mp.voice, mp.pcmChannel)
	mp.queue.Enqueue(url)
	mp.run()
	players[mp.voice.GuildID] = mp
	return
}

func (mp *musicPlayer) skipSong() {
	mp.skip = true
}

func (mp *musicPlayer) stopSong() {
	mp.stop = true
}

func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	mu.Lock()
	if sendpcm || pcm == nil {
		mu.Unlock()
		return
	}
	sendpcm = true
	mu.Unlock()
	defer func() { sendpcm = false }()

	opusEncoder, err := gopus.NewEncoder(48000, channels, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder Error:", err)
		return
	}

	for {
		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			fmt.Println("PCM Channel closed.")
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, 960, maxBytes)
		if err != nil {
			fmt.Println("Encoding Error:", err)
			return
		}

		if v.Ready == false || v.OpusSend == nil {
			// fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend)
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}
