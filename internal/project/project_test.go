package project

import (
	"fmt"
	"log"
	"testing"
  "io"
  "os"
)

// MARK: SETUP
func TestMain(m *testing.M) {
  log.SetOutput(io.Discard)
  os.Exit(m.Run())
}

// testDir is the source for where tests look the config file
const testDir = "../../testData/"
// project1 is the first example project used in testing
var project1 = Project {
    URL: "https://www.example.com",
    OFName: "project1",
}
// project2 is the second example project used in testing
var project2 = Project {
  URL: "https://www.example2.com",
  OFName: "project2",
}


// MARK: LoadProjects tests
// Tests the `LoadProjects` function with nominal data
func TestLoadProjectsSuccess(t *testing.T) {
  projects, err := LoadProjects(testDir)
  if err != nil {
    t.Fatalf("Unexpected err: %s", err)
  }

  if len(projects) != 2 {

    for project := range projects {
      fmt.Println(project)
    }

    t.Fatalf("Expected size 2, was: %d\n", len(projects))
  }


  if projects[0] != project1 {
    t.Fatal("Expected first project")
  }

  if projects[1] != project2 {
    t.Fatal("Expected second project")
  }
}

// Tests the `LoadProjects` function when the given path doesn't contain a projects.json file 
func TestLoadProjectsNoFile(t *testing.T) {
  _, err := LoadProjects("notadirectory")

  if err.Error() != "Failed to load projects from notadirectory" {
    t.Fatal("Error not thrown")
  }
}

// Tests the `LoadProjects` function when the project.json file doesn't contain project data
func TestLoadProjectsInvalidData(t *testing.T) {
  _, err := LoadProjects("../../testData/invalid")

  if err.Error() != "Failed to decode projects" {
    t.Fatal("Error not thrown")
  }
}

// MARK: GetProject tests
// Tests the `GetProject` function with valid data
func TestGetProjectSuccess(t *testing.T) {
  projects, err := LoadProjects(testDir)
  if err != nil {
    t.Fatal("Failed to load test project data")
  }

  project, err := GetProject("project1", projects)
  if err != nil {
    t.Fatalf("Unexpected error: %s", err)
  }

  if project != project1 {
    t.Fatalf("Expected project1, got: %s", project)
  }

  project, err = GetProject("project2", projects)
  if err != nil {
    t.Fatalf("Unexpected error: %s", err)
  }

  if project != project2 {
    t.Fatalf("Expected project2, got: %s", project)
  }
}

// Tests the `GetProject` function where the project doesn't exist
func TestGetProjectDoesNotExist(t *testing.T) {
  projects, err := LoadProjects(testDir)
  if err != nil {
    t.Fatal("Failed to load test project data")
  }

  project, err := GetProject("non existing project", projects)

  if err.Error() != "Project `non existing project` does not exist" {
    t.Fatal("Error not thrown")
  }

  if project != (Project{}) {
    t.Fatal("Project is non empty")
  }
}
