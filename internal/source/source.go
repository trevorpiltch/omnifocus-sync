package source

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	omnifocus "github.com/trevorpiltch/omnifocus-sync/internal/OF"
)

// Header represent a header that we want to attach in a HTTP request
type Header struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// Response represents the return data from an API call
type Response struct {
	// The field that contains the data. Leave empty for a standard JSON array of responses
	DataField string `json:"DataField"`
	// The field representing the title of the item to add
	Title string `json:"Title"`
	// The url that links to the task from the original source
	URL string `json:"URL"`
	// The number of the item that
	Number string `json:"Number"`
}

// Source represents a location where we are getting the OmniFocus items from
type Source struct {
	// The name of the source
	Name string `json:"Name"`
	// The url where the source exists
	URL string `json:"URL"`
	// The headers to include in the API request
	Headers []Header `json:"Headers"`
	// The queries to attach to the API request
	Queries string `json:"Queries"`
	// The response from the API request
	Response Response `json:"Response"`
	// The tags to add to OmniFocus items when they're added
	Tags []string `json:"Tags"`
}

// MARK: Private helper methods
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

// parseResponse parses an array of bytes into an array of new OmniFocus tasks
func (source Source) parseResponse(data []byte) ([]omnifocus.NewOmniFocusItem, error) {
	if source.Response.DataField != "" {
		var m map[string]interface{}
		err := json.Unmarshal(data, &m)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sources: %s", err)
		}

		data, err = json.Marshal(m[source.Response.DataField])
		if err != nil {
			return nil, fmt.Errorf("failed to remarshal data: %s", err)
		}
	}

	var responses []map[string]interface{}

	err := json.Unmarshal(data, &responses)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sources: %s", err)
	}

	var items []omnifocus.NewOmniFocusItem
	for _, response := range responses {
		name := fmt.Sprintf("[%d] %s", int(response[source.Response.Number].(float64)), response[source.Response.Title].(string))
		item := omnifocus.NewOmniFocusItem{
			Name: name,
			Tags: source.Tags,
			Note: response[source.Response.URL].(string),
		}

		items = append(items, item)
	}

	return items, nil
}

// MARK: Public methods
// loadSources parses the `sources.json` file at the given path and returns an array of Source objects
func LoadSources(Path string) ([]Source, error) {
	log.Printf("[source] Getting sources from: %s\n", Path)

	sourcePath := path.Join(Path, "sources.json")

	var bytes []byte
	bytes, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load sources from %s", Path)
	}

	var sources []Source

	err = json.Unmarshal(bytes, &sources)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sources")
	}

	return sources, nil
}

// GetItems creates an API request to the Item Source and returns an array of items to be added to OmniFocus
func (source Source) GetItems() ([]omnifocus.NewOmniFocusItem, error) {
	log.Printf("[source] Getting items from %s", source.URL)

	url := source.createURL()
	req, err := source.createRequest(url)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	return source.parseResponse(body)
}

// GetTags returns an array of all the tags associated with the sources
func GetTags(sources []Source) []string {
	var tags []string
	for _, source := range sources {
		tags = append(tags, source.Tags...)
	}

	return tags
}
