package omnifocus

import (
  "fmt"
  "log"

  "github.com/trevorpiltch/omnifocus-sync/internal/project"
)

// Item is an existing item in OmniFocus
type Item struct {
  ID string `json:"id"`
  Name string `json:"name"`
} 

func (i Item) String() string {
  return fmt.Sprintf("OmniFocus Item: [%s] %s", i.ID, i.Name)
}

// Key returns the string representation of the item for conformance to Delta's key interface 
func (i Item) Key() string {
  return i.String()
}

// ItemQuery defines a query to find OmniFocus items with the given criteria
type ItemQuery struct {
  ProjectName string `json:"projectName"`
  Tags []string `json:"tags"`
}

// Tag represents a tag in OmniFocus 
type Tag struct {
  Name string
}

// NewOmniFocusItem defines a request to create a new Item in OmniFocus
type NewOmniFocusItem struct {
  ProjectName string `json:"projectName"`
  Name string `json:"name"`
  Tags []string `json:"tags"`
  Note string `json:"note"`
  DueDateMS int64 `json:"dueDateMS"`
}

func (i NewOmniFocusItem) String() string {
  return fmt.Sprintf("[NewItem] %s: %s", i.Name, i.ProjectName)
}

// GetAllItems returns an array containing all of the items with the given tags for the given list of projects from OmniFocus
func GetAllItems(projects []project.Project, tags []string) ([]Item, error) {
  log.Print("[OF] Getting all items")
  items := []Item{}

  for i := 0; i < len(projects); i++ {
    query := ItemQuery {
      ProjectName: projects[i].OFName,
      Tags: tags,
    }
    projectItems , err := ItemsForQuery(query)

    if err != nil {
      return nil, err
    }

    items = append(items, projectItems[:]...) 
  }

  return items, nil
}

// GetItems returns n array containing all of the items with the tags from the given project in OmniFocus 
func GetItems(project project.Project, tags []string) ([]Item, error) {
  log.Printf("[OF] Getting items from %s", project)
  query := ItemQuery {
    ProjectName: project.OFName,
    Tags: tags,
  }
  items, err := ItemsForQuery(query)

  if err != nil {
    return nil, err

  }
  return items, nil
}

// AddItems adds the item to the OmniFocus application
func AddItem(i NewOmniFocusItem) error {
  log.Printf("[OF] Adding item %s", i)
  log.Printf("AddItem: %s", i)

  _, err := AddNewOmnifocusItem(i)
  if err != nil {
    return fmt.Errorf("Failed to add item: %v", err)
  }

  return nil
}

// CompleteItem completes the item in the OmniFocus application
func CompleteItem(i Item) error {
  log.Printf("[OF] Complete item %s", i)
  err := MarkOmnifocusItemComplete(i)
  if err != nil {
    return fmt.Errorf("Failed to complete item: %v", err)
  }

  return nil
}
