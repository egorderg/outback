package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartApp(g *gin.Context) {
	id := g.Query("id")
	if id == "" {
		return
	}

	if err := app.Start(id); err != nil {
		log.Println(err)
	}

	g.Redirect(302, "/")
}

func CloseApp(g *gin.Context) {
	if err := app.Close(); err != nil {
		log.Println(err)
	}

	g.Redirect(302, "/")
}

func Mount(g *gin.Context) {
	id := g.Query("id")
	if id == "" {
		return
	}

	if err := player.Close(); err != nil {
		log.Println(err)
	}

	if err := device.Mount(id); err != nil {
		log.Println(err)
	}

	g.Redirect(302, "/")
}

func Unmount(g *gin.Context) {
	if err := player.Close(); err != nil {
		log.Println(err)
	}

	if err := device.Unmount(); err != nil {
		log.Println(err)
	}

	g.Redirect(302, "/")
}

func Play(g *gin.Context) {
	g.Redirect(302, "/")

	id := g.Query("id")
	if id == "" {
		return
	}

	track, exists := device.TrackById(id)
	if !exists {
		log.Printf("Unknown track %s\n", id)
		return
	}

	player.Play(track)
}

func Stop(g *gin.Context) {
	player.Close()
	g.Redirect(302, "/")
}

func TogglePause(g *gin.Context) {
	player.TogglePause()
}

func ShowProgres(g *gin.Context) {
	player.ShowProgress()
}

func ChangeAudio(g *gin.Context) {
	player.ChangeAudio()
}

func ChangeSubtitles(g *gin.Context) {
	player.ChangeSubtitles()
}

func SeekToStart(g *gin.Context) {
	player.SeekToStart()
}

func Seek(g *gin.Context) {
	v := g.Query("value")
	if v == "" {
		return
	}

	value, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		log.Printf("Unknown value %v\n", value)
	}

	player.Seek(int(value))
}

func BigBack(g *gin.Context) {
	player.BigBack()
}

func Back(g *gin.Context) {
	player.Back()
}

func Skip(g *gin.Context) {
	player.Skip()
}

func BigSkip(g *gin.Context) {
	player.BigSkip()
}
