package data

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type MusicPlayer struct {
	CurrentMusic *audio.Player
	Muted        bool
	Volume       float64
}

// Play constructs a new ebiten audio player, starts playing, and returns it. Volume is 0-1.
func (mp *MusicPlayer) Play(s *Sound) {
	if mp.CurrentMusic != nil {
		mp.CurrentMusic.Close()
	}

	var player *audio.Player
	if mp.Muted {
		player = s.Play(0)
	} else {
		player = s.Play(mp.Volume)
	}
	mp.CurrentMusic = player
}

func (mp *MusicPlayer) Stop() {
	if mp.CurrentMusic != nil {
		mp.CurrentMusic.Close()
	}
}

func (mp *MusicPlayer) Set(p string) {
	bgm, err := GetMusic(p)
	if err != nil {
	}
	mp.Play(bgm)
}

func (mp *MusicPlayer) SetVolume(v float64) {
	mp.Volume = v
	mp.CurrentMusic.SetVolume(mp.Volume)
}

func (mp *MusicPlayer) ToggleMuted() {
	if mp.Muted {
		mp.Muted = false
		mp.CurrentMusic.SetVolume(mp.Volume)
	} else {
		mp.Muted = true
		mp.CurrentMusic.SetVolume(0)
	}
}

func (mp *MusicPlayer) Update() {
	if mp.CurrentMusic != nil && !mp.CurrentMusic.IsPlaying() {
		mp.CurrentMusic.Rewind()
	}
}

var BGM MusicPlayer = MusicPlayer{
	Volume: 0.25,
}
