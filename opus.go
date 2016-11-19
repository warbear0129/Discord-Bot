package main

import (
	"log"
)

const (
	sampleRate = 48000
	channels = 2
	bufferSize = 1000
)

func (mp *musicPlayer) processPCM(pcm []int16) {

	opusSend := mp.voice.OpusSend

	for {

		if !mp.voice.Ready {
			log.Println("**** Voice channel is not ready ****")
			return
		}

		frameSize := len(pcm)
		frameSizeMs := float32(frameSize) / channels * 1000 / sampleRate
		switch frameSizeMs {
		case 2.5, 5, 10, 20, 40, 60:
			// Good.
		default:
			log.Println("**** Invalid PCM frame size ****")
			return
		}

		data := make([]byte, bufferSize)
		n, err := mp.encoder.Encode(pcm, data)

		if err != nil {
			log.Println("**** Encoding error ****")
			log.Println(err)
		}

		data = data[:n]

		opusSend <- data
	}
}
