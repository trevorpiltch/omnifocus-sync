package project

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

// Project represents the connection between a source of items and an OmniFocus project.
type Project struct {
	// The URL from which all issues should be added
	URL string
	// The name of the OF project for all the issues
	OFName string
}

// String prints out a formatted description of the project
func (p *Project) String() string {
	return fmt.Sprintf("%s: %s", p.OFName, p.URL)
}

// Loads projects from the given path
func LoadProjects(Path string) ([]Project, error) {
	log.Printf("[project] Getting projects from: %s\n", Path)

	projectPath := path.Join(Path, "projects.json")

	var bytes []byte
	bytes, err := os.ReadFile(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load projects from %s", Path)
	}

	var projects []Project

	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		return nil, fmt.Errorf("failed to decode projects %s", err)
	}

	return projects, nil
}

// GetProject returns the project with the name `key` in OmniFocus
func GetProject(key string, projects []Project) (Project, error) {
	for i := range projects {
		if projects[i].OFName == key {
			return projects[i], nil
		}
	}

	return Project{}, fmt.Errorf("Project `%s` does not exist", key)
}

func ProjectFor(url string, projects []Project) (Project, error) {
	for _, project := range projects {
		if strings.Contains(url, project.URL) {
			return project, nil
		}
	}

	return Project{}, fmt.Errorf("URL %s does not match any projects", url)
}
