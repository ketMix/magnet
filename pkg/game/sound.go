package game

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Sound struct {
	bytes []byte
}

func NewSound(data []byte) (*Sound, error) {
	// Attempt to read the vorbis file in 44100 sample rate.
	stream, err := vorbis.DecodeWithSampleRate(44100, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Read the decoded stream into a buffer.
	var b bytes.Buffer
	_, err = b.ReadFrom(stream)
	if err != nil {
		return nil, err
	}

	// Make our sound.
	s := &Sound{
		bytes: b.Bytes(),
	}
	return s, nil
}

// Play constructs a new ebiten audio player, starts playing, and returns it. Volume is 0-1.
func (s *Sound) Play(volume float64) *audio.Player {
	player := audio.CurrentContext().NewPlayerFromBytes(s.bytes)
	player.SetVolume(volume)
	player.Play()
	return player
}
