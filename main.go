package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/farbodsalimi/promptctl/cmd"
)

var Version = "development"

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
}

func main() {
	if err := cmd.Execute(Version); err != nil {
		log.Fatal(err)
	}
}
