package data

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

var BGM MusicPlayer = MusicPlayer{
	volume: 0.50,
}

type MusicPlayer struct {
	currentTrack string
	currentMusic *audio.Player
	Muted        bool
	volume       float64
}

func (mp *MusicPlayer) Play(s *Sound) {
	if mp.currentMusic != nil {
		mp.currentMusic.Close()
	}

	var player *audio.Player
	if mp.Muted {
		player = s.Play(0)
	} else {
		player = s.Play(mp.volume)
	}
	mp.currentMusic = player
}

func (mp *MusicPlayer) Stop() {
	if mp.currentMusic != nil {
		mp.currentMusic.Close()
	}
}

func (mp *MusicPlayer) Set(p string) {
	if mp.currentTrack == p {
		return
	}
	bgm, err := GetMusic(p)
	if err != nil {
	}
	mp.currentTrack = p
	mp.Play(bgm)
}

func (mp *MusicPlayer) SetVolume(v float64) {
	mp.volume = v
	mp.currentMusic.SetVolume(mp.volume)
}

func (mp *MusicPlayer) ToggleMute() {
	mp.Muted = !mp.Muted
	if mp.Muted {
		mp.currentMusic.SetVolume(0)
	} else {
		mp.currentMusic.SetVolume(mp.volume)
	}
}

func (mp *MusicPlayer) Update() {
	if mp.currentMusic != nil && !mp.currentMusic.IsPlaying() {
		mp.currentMusic.Rewind()
		mp.currentMusic.Play()
	}
}

func (mp *MusicPlayer) GetCurrentTrack() string {
	return mp.currentTrack
}

func (mp *MusicPlayer) GetAllTracks() []string {
	tracks, _ := GetPathFiles("music")
	return tracks
}
func FormatTrackName(p string) string {
	trackName := strings.Split(p, ".")[0]
	trackName = strings.ToUpper(string(trackName[0])) + trackName[1:]
	return trackName
}
