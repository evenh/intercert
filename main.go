package main

import (
	"github.com/evenh/intercert/cmd"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	cmd.Execute()
}
