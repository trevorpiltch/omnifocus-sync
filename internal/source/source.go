package source

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

// Header represent a header that we want to attach in a HTTP request
type Header struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// Response represents the return data from an API call
type Response struct {
	// The field representing the title of the item to add
	Title string `json:"Title"`
	// The url that links to the task from the original source
	URL string `json:"URL"`
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
	// The response from the API request
	Response Response `json:"response"`
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

// loadSources parses the `sources.json` file at the given path and returns an array of Source objects
func loadSources(Path string) ([]Source, error) {
	log.Printf("[source] Getting sources from: %s\n", Path)

	sourcePath := path.Join(Path, "sources.json")

	var bytes []byte
	bytes, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load sources from %s", Path)
	}

	var sources []Source

	err = json.Unmarshal(bytes, &sources)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode sources")
	}

	return sources, nil
}

/*
func (source Source) GetItems() ([]omnifocus.Item, error) {
	log.Printf("[source] Getting items from %s", source.URL)

	url := source.createURL()
	req, err := source.createRequest(url)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("[Source] Failed to send request")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	var items []omnifocus.Item

}
*/
