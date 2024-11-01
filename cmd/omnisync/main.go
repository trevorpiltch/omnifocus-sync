package main

import (
	"log"
	"os"
	"path"

	"github.com/trevorpiltch/omnifocus-sync/internal/project"
)

const version = "1.0.0"

func main() {
  log.Printf("[main] Starting OmniSync version: %s", version)

  home, err := os.UserHomeDir()
  if err != nil {
    log.Fatal("Failed to find user home directory.")
  }

  configDir := path.Join(home, ".config", "omnisync") 

  projects, err := project.LoadProjects(configDir)
  if err != nil {
    log.Fatal(err)
  }

  log.Printf("%d", len(projects))
}
