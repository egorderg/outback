package com

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const META_EXT = ".meta"
const IMAGE_EXT = ".png"

type Track struct {
	Id       string
	Title    string
	Path     string
	Preview  string
	Metadata Metadata
}

type Metadata struct {
	Time           string
	Audio          string
	AudioCount     int
	Subtitles      string
	SubtitlesCount int
}

func loadTracks(dir string) (map[string]*Track, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Track error: %s", err)
	}

	tracks := make(map[string]*Track, 0)

	for _, e := range entries {
		if e.IsDir() || path.Ext(e.Name()) != ".mkv" {
			continue
		}

		title := strings.TrimSuffix(e.Name(), path.Ext(e.Name()))
		preview, err := getPreview(dir, title)
		if err != nil {
			log.Printf("Track error: couldn't read preview '%s'\n", err)
		}

		metadata, nil := getMetadata(dir, title)
		if err != nil {
			log.Printf("Track error: couldn't read metadata '%s'", err)
		}

		tracks[title] = &Track{
			Id:       title,
			Title:    title,
			Path:     path.Join(dir, e.Name()),
			Preview:  preview,
			Metadata: metadata,
		}
	}

	return tracks, nil
}

func getPreview(dir string, name string) (string, error) {
	data, err := os.ReadFile(path.Join(dir, name+IMAGE_EXT))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func getMetadata(dir string, name string) (Metadata, error) {
	file, err := os.Open(path.Join(dir, name+META_EXT))
	if err != nil {
		return Metadata{}, err
	}

	metadata := Metadata{}
	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		switch i {
		case 0:
			ms, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err == nil {
				metadata.Time = time.Unix(0, ms*int64(time.Millisecond)).Format("04:05")
			}
		case 1:
			metadata.Audio = strings.ToUpper(strings.TrimRight(scanner.Text(), ","))
			metadata.Audio = strings.Replace(metadata.Audio, ",", ", ", -1)
			metadata.AudioCount = len(strings.Split(metadata.Audio, ","))
		case 2:
			metadata.Subtitles = strings.ToUpper(strings.TrimRight(scanner.Text(), ","))
			metadata.Subtitles = strings.Replace(metadata.Subtitles, ",", ", ", -1)
			metadata.SubtitlesCount = len(strings.Split(metadata.Subtitles, ","))
		}
	}

	return metadata, nil
}
