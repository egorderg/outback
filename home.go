package main

import (
	"egorderg/outback/com"

	"github.com/gin-gonic/gin"
)

const DEFAULT_TITLE = "Outback"

const (
	modeIdle = iota
	modeMounted
	modePlaying
	modeApp
)

func Home(c *gin.Context) {
	mode := getMode()
	title := DEFAULT_TITLE

	if mode == modeApp {
		title = app.Title()
	}

	c.HTML(200, "index.html", gin.H{
		"title":         title,
		"apps":          app.Descriptions(),
		"tracks":        device.Tracks(),
		"devices":       com.DEVICES,
		"isModeIdle":    mode == modeIdle,
		"isModeMounted": mode == modeMounted,
		"isModePlaying": mode == modePlaying,
		"isModeApp":     mode == modeApp,
	})
}

func getMode() int {
	if app.Running() {
		return modeApp
	}

	if device.Mounted() {
		if player.IsPlaying() {
			return modePlaying
		}

		return modeMounted
	}

	if device.Mounted() && player.IsPlaying() {
		return modePlaying
	}

	return modeIdle
}
