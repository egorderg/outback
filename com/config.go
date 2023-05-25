package com

type AppDescription struct {
	Title string
	cmd   string
	args  []string
}

type DeviceDescription struct {
	Title  string
	dev    string
	fstype string
}

type AppDescriptions map[string]AppDescription
type DeviceDescriptions map[string]DeviceDescription

const MOUNT_DIR string = "/mnt/outback"
const USE_BG bool = false
const USE_MUSIC bool = false

var DEVICES DeviceDescriptions = DeviceDescriptions{
	"optical": DeviceDescription{
		Title:  "DVD",
		dev:    "/dev/sr0",
		fstype: "iso9660",
	},
	"usb": DeviceDescription{
		Title:  "USB",
		dev:    "/dev/sdc1",
		fstype: "vfat",
	},
}
var APPS = AppDescriptions{
	// "kodi": AppDescription{
	// 	Title: "Kodi",
	// 	cmd:   "kodi",
	// 	args:  []string{},
	// },
}

type Config struct {
	addr     string
	mountDir string
	apps     AppDescriptions
}

func NewConfig(addr string) Config {
	return Config{
		addr:     addr,
		mountDir: MOUNT_DIR,
		apps:     APPS,
	}
}
