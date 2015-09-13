package crawler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	service      *youtube.Service
	developerKey = os.Getenv("FURRYTEMPLE_YOUTUBE_DEVELOPER_KEY")
)

const (
	maxResults = 25
)

func createService() (*youtube.Service, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)

	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	return service, err
}

func init() {
	service, _ = createService()
}

var increase int

func walk(call *youtube.SearchListCall, ch chan<- []byte, pageToken string) {
	call = call.PageToken(pageToken)
	response, _ := call.Do()

	jsonResponse, _ := json.Marshal(response)
	ch <- jsonResponse

	walk(call, ch, response.NextPageToken)
}

func Crawl(keyword string, out chan<- []byte) {
	go func(ch chan<- []byte, query string) {
		call := service.Search.List("id,snippet").Q(query).MaxResults(maxResults).Order("relevance")
		response, _ := call.Do()

		jsonResponse, _ := json.Marshal(response)
		ch <- jsonResponse

		walk(call, ch, response.NextPageToken)
	}(out, keyword)
}
