package main

import (
	"fmt"
	"log"
	"os"

	storage "github.com/devamaz/clipshistory/internal/store"
	"github.com/devamaz/clipshistory/internal/tui"
	"golang.design/x/clipboard"
)

func main() {
	// Initialize clipboard library before any clipboard operations
	if err := clipboard.Init(); err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	store := new(storage.Store)
	if err := store.Init(); err != nil {
		log.Fatalf("unable to init store: %v", err)
	}

	return tui.Start(store)
}
