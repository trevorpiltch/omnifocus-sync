package source

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/trevorpiltch/omnifocus-sync/internal/OF"
)

// Header represent a header that we want to attach in a HTTP request
type Header struct {
  Key string `json: "Key"`
  Value string `json:"Value"`
}

// Source represents a location where we are getting the OmniFocus items from
type Source struct {
  // The url where the source exists 
  URL string `json:"URL"`
  // The token for the url
  APIToken string `json:"APIToken"`
  // The headers to include in the API request
  Headers []Header `json:"headers"`
  // The queries to attach to the API request
  Queries string ""
}

// Creates the full URL from the source
func (source Source) createURL() string {
  params := url.Values{}

  params.Add("query", source.Queries)
  return source.URL + "?" + params.Encode()
}

// Creates the request from the source using the given url
func (source Source) createRequest(url string) (*http.Request, error) {
  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
    return nil, err
  }

  for i := range source.Headers {
    req.Header.Set(source.Headers[i].Key, source.Headers[i].Value)
  }
  
  return req, nil
}

func (source Source) GetItems() ([]omnifocus.Item, error) {
  log.Printf("[Source] Getting items from %s",source.URL)

  url := source.createURL()
  req, err := source.createRequest(url)

  client := http.Client {
    Timeout: 30 * time.Second,
  }

  res, err := client.Do(req)
  if err != nil {
    log.Fatal("[Source] Failed to send request")
  }

  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)

  var items []omnifocus.Item



  err = json.Marshal(body, &items)

} 

