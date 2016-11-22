package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/layeh/gopus"
	"log"
)

const (
	channels = 2
	bitrate = 320000
	maxBytes = 3840
	sampleRate = 48000
)

func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
        mu.Lock()
        if sendpcm || pcm == nil {
                mu.Unlock()
                return
        }
        sendpcm = true
        mu.Unlock()
        defer func() { sendpcm = false }()

        opusEncoder, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
        opusEncoder.SetBitrate(gopus.BitrateMaximum)
        if err != nil {
                log.Println("NewEncoder Error:", err)
                return
        }

        for {
                // read pcm from chan, exit if channel is closed.
                recv, ok := <-pcm
                if !ok {
                        log.Println("PCM Channel closed.")
                        return
                }

                // try encoding pcm frame with Opus
                opus, err := opusEncoder.Encode(recv, 960, maxBytes)
                if err != nil {
                        log.Println("Encoding Error:", err)
                        return
                }

                if v.Ready == false || v.OpusSend == nil {
                        // log.Printf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend)
                        return
                }
                // send encoded opus data to the sendOpus channel
                v.OpusSend <- opus
        }
}

