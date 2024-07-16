package main

import (
	"log"
	"os"

	"github.com/quansolashi/message-extractor/backend/internal/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}
}
