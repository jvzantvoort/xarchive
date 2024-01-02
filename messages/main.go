package messages

import (
	"embed"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Content missing godoc.
//
//go:embed usage/* long/*
var Content embed.FS

func GetUsage(name string) string {
	filename := fmt.Sprintf("usage/%s", name)
	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return strings.TrimRight(string(msgstr), "\n")
}

func GetLong(name string) string {
	filename := fmt.Sprintf("long/%s", name)
	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return string(msgstr)
}
