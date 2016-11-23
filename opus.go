package main

const (
	channels = 2
	bitrate = 320000
	maxBytes = 3840
	sampleRate = 48000
	frameSize = 960
)

func SendPCM(mp *musicPlayer, pcm <-chan []int16) {
	v := mp.voice
        enc := mp.encoder
	mu := mp.mutex
	sendpcm := mp.sendpcm

        mu.Lock()
        if sendpcm || pcm == nil {
                mu.Unlock()
                return
        }
        sendpcm = true
        mu.Unlock()
        defer func() { sendpcm = false }()

        for {
                recv, ok := <-pcm
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

