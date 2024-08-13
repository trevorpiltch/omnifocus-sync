package main

import (
	"fmt"
  "time"
  "net/url"
	"net/http"
  "encoding/json"
  "io"
)

type ShortcutStory struct {
  Title string
  HTMLURL string
  K string
}

type Shortcut struct {
  Name string`json:"name"`
  Id int`json:"id"`
  Url string`json:"app_url"`
  Completed bool`json:"completed"`
}

type APIResponse struct {
  Next string`json:"next"`
  Total int `json:"total"`
  Data []Shortcut`json:"data"`
}

func (item ShortcutStory) String() string {
  return fmt.Sprintf("ShortcutStory: [%s] %s (%s)", item.Key(), item.Title, item.HTMLURL)
}

func (item ShortcutStory) Key() string {
  return item.K
}

func GetStories(APIKey string, Owner string) ([]ShortcutStory, error) {
  fmt.Println("Getting Shortcut Stories")

  query := fmt.Sprintf("owner:%s -is:completed", Owner)

  // Create request
  baseURL := "https://api.app.shortcut.com/api/v3/search/stories"
  params := url.Values{}
  params.Add("query", query)

  fullURL := baseURL + "?" + params.Encode()
  req, err := http.NewRequest(http.MethodGet, fullURL, nil)

  if err != nil {
    fmt.Printf("Couldn't create request")
    return nil, err
  }

  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Shortcut-Token", APIKey)

  client := http.Client{
    Timeout: 30 * time.Second, 
  }

  // Send Request
  res, err := client.Do(req)

  if err != nil {
    fmt.Printf("Couldn't make request: %s\n", err)
    return nil, err
  }

  defer res.Body.Close()

  // Decode response data
  body, err := io.ReadAll(res.Body)

  var response APIResponse
  var shortcuts []ShortcutStory

  err = json.Unmarshal(body, &response)

  if err != nil {
    fmt.Printf("Couldn't decode data\n")
    return nil, err
  }

  // Convert from API format to go script format
  for _, story := range response.Data {
    if !story.Completed {
      shortcut := ShortcutStory {
        Title: fmt.Sprintf("[sc-%d] %s", story.Id, story.Name), 
        HTMLURL: story.Url, 
        K: fmt.Sprintf("[sc-%d]", story.Id), 
      }

      shortcuts = append(shortcuts, shortcut)
    }
  }

  return shortcuts, nil
}
