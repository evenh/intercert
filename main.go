package main

import (
	"github.com/evenh/intercert/cmd"
	"log"
	"os"
)

var (
	Version string
	Commit string
	BuiltAt string
)

func main() {
	log.SetOutput(os.Stdout)

	cmd.Execute()
}
