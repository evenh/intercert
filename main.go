package main

import (
	"log"
	"os"

	"github.com/evenh/intercert/cmd"
)

var (
	Version = "DEV-SNAPSHOT"
	Commit  = "N/A"
)

func main() {
	log.SetOutput(os.Stdout)

	cmd.Version = Version
	cmd.Commit = Commit

	cmd.Execute(Version, Commit)
}
