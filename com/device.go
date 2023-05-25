package com

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

const BG_FILE = "background.png"
const MUSIC_FILE = "background.mp3"

var AlreadyMountedError = errors.New("Device error: already mounted")
var AlreadyUnmountedError = errors.New("Device error: already unmounted")
var NotMountedError = errors.New("Device error: not mounted")

type Device struct {
	config  Config
	mounted bool
	tracks  map[string]*Track
}

func NewDevice(config Config) *Device {
	return &Device{
		config:  config,
		mounted: false,
	}
}

func (d *Device) Mounted() bool {
	return d.mounted
}

func (d *Device) Tracks() map[string]*Track {
	return d.tracks
}

func (d *Device) TrackById(id string) (t *Track, ok bool) {
	t, ok = d.tracks[id]
	return
}

func (d *Device) Mount(id string) error {
	if d.mounted {
		return AlreadyMountedError
	}

	dev, ok := DEVICES[id]
	if !ok {
		return fmt.Errorf("Device error: unknown device '%s'", id)
	}

	out, _ := exec.Command("sh", "-c", fmt.Sprintf("mount | grep %v", d.config.mountDir)).Output()
	if len(out) > 0 {
		tracks, err := loadTracks(d.config.mountDir)
		if err != nil {
			return err
		}

		d.mounted = true
		d.tracks = tracks
		return nil
	}

	if _, err := os.Stat(d.config.mountDir); os.IsNotExist(err) {
		os.Mkdir(d.config.mountDir, os.ModePerm)
	}

	err := syscall.Mount(dev.dev, d.config.mountDir, dev.fstype, syscall.MS_RDONLY, "")
	if err != nil {
		return fmt.Errorf("Device error: %s for %s to %s", err, dev.dev, d.config.mountDir)
	}

	tracks, err := loadTracks(d.config.mountDir)
	if err != nil {
		return err
	}

	if USE_BG {
		bg := path.Join(d.config.mountDir, BG_FILE)
		if _, err := os.Stat(bg); os.IsNotExist(err) {
			resetBackground()
		} else {
			setBackground(bg)
		}
	}

	if USE_MUSIC {
		music := path.Join(d.config.mountDir, MUSIC_FILE)
		if _, err := os.Stat(music); os.IsNotExist(err) {
			stopMusic()
		} else {
			playMusic(music)
		}
	}

	d.tracks = tracks
	d.mounted = true

	return nil
}

func (d *Device) Unmount() error {
	if !d.mounted {
		return AlreadyUnmountedError
	}

	if err := syscall.Unmount(d.config.mountDir, 0); err != nil {
		return fmt.Errorf("Device error: %s for %s", err, d.config.mountDir)
	}

	d.mounted = false
	d.tracks = nil

	stopMusic()
	resetBackground()

	return nil
}

func setBackground(bg string) {
	exec.Command("feh", "--no-fehbg", "--bg-fill", bg).Start()
}

func resetBackground() {
	exec.Command("xsetroot", "-solid", "#000000").Start()
}

func playMusic(music string) {
	exec.Command("paplay", music).Start()
}

func stopMusic() {
	exec.Command("pkill", "paplay").Run()
}
