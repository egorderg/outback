package com

import (
	"fmt"
	"os"
	"os/exec"
)

type App struct {
	descriptions AppDescriptions
	title        string
	process      *os.Process
	running      bool
}

func NewApp(config Config) *App {
	return &App{
		descriptions: config.apps,
		running:      false,
	}
}

func (app *App) Descriptions() AppDescriptions {
	return app.descriptions
}

func (app *App) Title() string {
	return app.title
}

func (app *App) Running() bool {
	return app.running
}

func (app *App) Start(id string) error {
	d, ok := app.descriptions[id]
	if !ok {
		return fmt.Errorf("App error: unknown app '%s'", id)
	}

	cmd := exec.Command(d.cmd, d.args...)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("App error: %s", err)
	}

	app.title = d.Title
	app.process = cmd.Process
	app.running = true

	return nil
}

func (app *App) Close() error {
	if !app.running {
		return nil
	}

	if err := app.process.Kill(); err != nil {
		return fmt.Errorf("App error: %s", err)
	}

	app.process.Wait()

	app.running = false
	return nil
}
