package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/egorderg/outback/com"

	"github.com/gin-gonic/gin"
)

var app *com.App
var player *com.Player
var device *com.Device

func main() {
	log.SetPrefix("[Outback] ")

	addr := getArgs()
	config := com.NewConfig(addr)
	template, err := loadTemplate()
	if err != nil {
		log.Fatalln(err)
	}

	app = com.NewApp(config)
	player = com.NewPlayer(config)
	device = com.NewDevice(config)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	{
		router.SetTrustedProxies(nil)
		// router.LoadHTMLFiles("templates/index.html")
		// go-assets-builder templates/index.html -s /templates/ -o assets.go
		router.SetHTMLTemplate(template)
	}

	router.Use(
		func(c *gin.Context) {
			c.Writer.Header().Set("Cache-Control", "no-cache")
		},
	)

	commands := router.Group("/commands")
	{
		commands.GET("/mount", Mount)
		commands.GET("/start-app", StartApp)
		commands.GET("/close-app", CloseApp)
		commands.GET("/unmount", Unmount)
		commands.GET("/play", Play)
		commands.GET("/stop", Stop)

		commands.POST("/toggle", TogglePause)
		commands.POST("/progress", ShowProgres)
		commands.POST("/audio", ChangeAudio)
		commands.POST("/subtitles", ChangeSubtitles)
		commands.POST("/seek-start", SeekToStart)
		commands.POST("/seek", Seek)
		commands.POST("/big-back", BigBack)
		commands.POST("/back", Back)
		commands.POST("/skip", Skip)
		commands.POST("/big-skip", BigSkip)
	}

	router.GET("/", Home)

	log.Printf("Starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalln(err)
	}
}

func getArgs() string {
	switch len(os.Args) {
	case 2:
		return os.Args[1]
	default:
		return "127.0.0.1:3000"
	}
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".html") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
