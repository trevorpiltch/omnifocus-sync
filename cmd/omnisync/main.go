package main

import (
	"log"
	"os"
	"path"

	omnifocus "github.com/trevorpiltch/omnifocus-sync/internal/OF"
	"github.com/trevorpiltch/omnifocus-sync/internal/delta"
	"github.com/trevorpiltch/omnifocus-sync/internal/project"
	"github.com/trevorpiltch/omnifocus-sync/internal/source"
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

	sources, err := source.LoadSources(configDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, source := range sources {
		log.Printf("[main] **** %s ****", source.Name)

		currentState, err := omnifocus.GetAllItems(projects, source.Tags)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[main] Current state: %d\n", len(currentState))

		items, err := source.GetItems()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[main] Desired state: %d\n", len(items))

		for i, item := range items {
			projectName, _ := project.ProjectFor(item.Note, projects)
			items[i].ProjectName = projectName.OFName
		}

		d := delta.Delta(toSetSource(items), toSet(currentState))

		log.Printf("[main] Found %d changes to apply", len(d))
		for _, d := range d {
			if d.Type == delta.Add {
				err := omnifocus.AddItem(*(d.Item.(*omnifocus.NewOmniFocusItem)))
				if err != nil {
					log.Fatal(err)
				}
			} else if d.Type == delta.Remove {
				err := omnifocus.CompleteItem(*(d.Item.(*omnifocus.Item)))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func toSet(l []omnifocus.Item) map[delta.Keyed]struct{} {
	r := map[delta.Keyed]struct{}{}
	for _, i := range l {
		// need to clone because range reuses `i` for each item!
		r[&omnifocus.Item{
			ID:   i.ID,
			Name: i.Name,
		}] = struct{}{}
	}
	return r
}

func toSetSource(l []omnifocus.NewOmniFocusItem) map[delta.Keyed]struct{} {
	r := map[delta.Keyed]struct{}{}
	for _, i := range l {
		// need to clone because range reuses `i` for each item!
		r[&omnifocus.NewOmniFocusItem{
			Name:        i.Name,
			ProjectName: i.ProjectName,
			Tags:        i.Tags,
			Note:        i.Name,
		}] = struct{}{}
	}
	return r
}
