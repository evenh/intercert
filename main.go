package main

import (
	"log"
	"os"

	"github.com/evenh/intercert/cmd"
)

var (
	// Version is the human readable version of intercert
	Version = "DEV-SNAPSHOT"
	// Commit is the git commit used to build this version
	Commit = "N/A"
)

func main() {
	log.SetOutput(os.Stdout)

	cmd.Version = Version
	cmd.Commit = Commit

	cmd.Execute(Version, Commit)
}
