package com

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const SOCKET_PATH = "/tmp/mpv-socket"

type Player struct {
	config    Config
	isPlaying bool
	process   *os.Process
	track     *Track
	sub       int
	audio     int
}

func NewPlayer(config Config) *Player {
	return &Player{
		config:    config,
		isPlaying: false,
		sub:       0,
		audio:     1,
	}
}

func (p *Player) IsPlaying() bool {
	return p.isPlaying
}

func (p *Player) Close() error {
	if !p.isPlaying {
		return nil
	}

	if err := p.process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("Player error: %s", err)
	}

	p.process.Wait()

	p.isPlaying = false
	return nil
}

func (p *Player) Play(t *Track) error {
	stopMusic()

	p.Close()

	cmd := exec.Command("mpv", "--save-position-on-quit", "--input-ipc-server="+SOCKET_PATH, "--fs", t.Path)
	err := cmd.Start()
	if err != nil {
		return err
	}

	p.isPlaying = true
	p.process = cmd.Process
	p.track = t

	return nil
}

func (p *Player) TogglePause() {
	sendCommand("cycle pause")
}

func (p *Player) ShowProgress() {
	sendCommand("show-progress")
}

func (p *Player) ChangeAudio() {
	p.audio = p.audio + 1
	if p.audio > p.track.Metadata.AudioCount {
		p.audio = 1
	}

	sendCommand(fmt.Sprintf("set audio %v", p.audio))
}

func (p *Player) ChangeSubtitles() {
	p.sub = p.sub + 1
	if p.sub > p.track.Metadata.SubtitlesCount {
		p.sub = 0
	}

	sendCommand(fmt.Sprintf("set sub %v", p.sub))
}

func (p *Player) SeekToStart() {
	sendCommand("seek 00:00:00 absolute")
}

func (p *Player) Seek(value int) {
	sendCommand(fmt.Sprintf("seek %v absolute-percent", value))
}

func (p *Player) BigBack() {
	sendCommand("seek -30")
}

func (p *Player) Back() {
	sendCommand("seek -10")
}

func (p *Player) Skip() {
	sendCommand("seek 30")
}

func (p *Player) BigSkip() {
	sendCommand("seek 30")
}

func sendCommand(cmd string) {
	cmd = "echo " + cmd + " | socat - " + SOCKET_PATH
	err := exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		log.Printf("Player error: %s\n", err)
	}
}
