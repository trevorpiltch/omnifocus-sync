package source

import (
  "testing"
  "net/http"
  "reflect"
)

// MARK: SETUP
// testDir is the source for where tests look the config file
const testDir = "../../testData/"

// source1 is the first example source used for testing
var source1 =  Source {
  Name:  "Source1", 
  URL: "www.example.com",
  Headers: []Header{Header { Key: "Header", Value: "Value"}},
  Queries: "query", 
  Response: Response{
    DataField: "",
    Title: "Title",
    URL: "url",
    Number: "number",
  }, 
  Tags: []string{"tag"},
}

// source2 is the second example source used for testing
var source2 =  Source {
  Name:  "Source2", 
  URL: "www.example2.com",
  Headers: []Header{},
  Queries: "", 
  Response: Response{
    DataField: "Data",
    Title: "Title",
    URL: "url",
    Number: "number",
  }, 
  Tags: []string{"tag2"},
}

func headersEqual(h1, h2 http.Header) bool {
    if len(h1) != len(h2) {
        return false
    }
    for k, v1 := range h1 {
        v2, ok := h2[k]
        if !ok {
          return false
        }

        if len(v1) != len(v2) {
          return false
        }
        for i := range v1 {
          if v1[i] != v2[i] {
            return false
          }
        }
    }
    return true
}

func requestsEqual(req1, req2 *http.Request) bool {
    return req1.Method == req2.Method &&
           req1.URL.String() == req2.URL.String() &&
           headersEqual(req1.Header, req2.Header) &&
           req1.ContentLength == req2.ContentLength
}

func TestCreateURL(t *testing.T) {
  url := source1.createURL();

  if url != "www.example.com?query=query" {
    t.Fatalf("Unexpected URL: %s", url);
  }

  url = source2.createURL() 

  if url != "www.example2.com" {
    t.Fatalf("Unexpected URL: %s", url);
  }
}

func TestCreateRequest(t *testing.T) {
  sourceURLString := source1.createURL();
  
  request, err := source1.createRequest(sourceURLString);
  if err != nil {
    t.Fatalf("Unexpected error: %s", err);
  }

  expectedRequest, err := http.NewRequest(http.MethodGet, sourceURLString, nil);
  expectedRequest.Header.Add("Header", "Value")
  if err != nil {
    t.Fatalf("Unexpected error: %s", err);
  }

  if !requestsEqual(request, expectedRequest) {
    t.Fatalf("Unexpected request for source 1");
  }

  sourceURLString = source2.createURL();
  
  request, err = source2.createRequest(sourceURLString);
  if err != nil {
    t.Fatalf("Unexpected error: %s", err);
  }

  expectedRequest, err = http.NewRequest(http.MethodGet, sourceURLString, nil);
  if err != nil {
    t.Fatalf("Unexpected error: %s", err);
  }

  if !requestsEqual(request, expectedRequest) {
    t.Fatalf("Unexpected request for source 2");
  }
}


func TestLoadSourcesSuccess(t *testing.T) {
  sources, err := LoadSources(testDir);
  if err != nil {
    t.Fatalf("Unexpected err: %s", err);
  }

  if len(sources) != 2 {
    t.Fatalf("Expected 2 sources, was %d", len(sources))
  }

  if !reflect.DeepEqual(sources[0], source1) {
    t.Fatal("Unexpected first source");
  }

  if !reflect.DeepEqual(sources[1], source2) {
    t.Fatal("Unexpected second source");
  }
}

