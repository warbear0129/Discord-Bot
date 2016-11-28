package main

import (
	"github.com/layeh/gopus"
)

const (
	channels = 2
	maxBytes = 3840
	sampleRate = 48000
	frameSize = 960
)

func SendPCM(mp *musicPlayer, pcm <-chan []int16) {
	var ok bool
	v := mp.voice
	recv := mp.recv
	enc, _ := gopus.NewEncoder(sampleRate, channels, gopus.Audio)

        for {
                recv, ok = <-pcm
		if !ok {
			return
		}

                opus, err := enc.Encode(recv, frameSize, maxBytes)
                if err != nil {
                        return
                }

                if !v.Ready || v.OpusSend == nil {
                        return
                }

                v.OpusSend <- opus
        }
}
