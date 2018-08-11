package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

const (
	hackerNewsItemID = "https://news.ycombinator.com/item?id="
	algoliaSearchURL = "https://hn.algolia.com/api/v1/search_by_date?"
	algoliaItemURL   = "https://hn.algolia.com/api/v1/items/"
)

type AlgoliaSearchResponse struct {
	Hits []AlgoliaSearchHit
}

type AlgoliaItemResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type AlgoliaSearchHit struct {
	Tags        []string `json:"_tags"`
	ObjectID    string
	Title       string
	URL         string
	Author      string
	CreatedAt   string `json:"created_at"`
	StoryTitle  string `json:"story_title"`
	CommentText string `json:"comment_text"`
	NumComments int    `json:"num_comments"`
	Points      int    `json:"points"`
}

func (hit AlgoliaSearchHit) isComment() bool {
	for _, tag := range hit.Tags {
		if tag == "comment" {
			return true
		}
	}
	return false
}

func (hit AlgoliaSearchHit) GetTitle() string {
	if hit.isComment() {
		return fmt.Sprintf("New comment by %s in \"%s\"", hit.Author, hit.StoryTitle)
	} else {
		return hit.Title
	}
}

func (hit AlgoliaSearchHit) GetPermalink() string {
	return hackerNewsItemID + hit.ObjectID
}

func (hit AlgoliaSearchHit) GetURL(linkTo string) string {
	if linkTo == "" {
		linkTo = "url"
	}

	switch {
	case linkTo == "url" && hit.URL != "":
		return hit.URL
	default:
		return hit.GetPermalink()
	}
}

func (hit AlgoliaSearchHit) GetDescription() string {
	if hit.isComment() {
		return hit.CommentText
	} else {
		var b bytes.Buffer
		// TODO(ejd): Hide article URL if not available
		t := template.Must(template.New("description").Parse(`
<p>Article URL: <a href="{{ .URL }}">{{ .URL }}</a></p>
<p>Comments URL: <a href="{{ .GetPermalink }}">{{ .GetPermalink }}</a></p>
<p>Points: {{ .Points }}</p>
<p># Comments: {{ .NumComments }}</p>
`))
		t.Execute(&b, hit)
		return b.String()
	}
}

func (hit AlgoliaSearchHit) GetCreatedAt() time.Time {
	rv, err := time.Parse("2006-01-02T15:04:05.000Z", hit.CreatedAt)
	if err != nil {
		return time.Now().UTC()
	} else {
		return rv
	}
}

func GetResults(params url.Values) (*AlgoliaSearchResponse, error) {
	resp, err := http.Get(algoliaSearchURL + params.Encode())
	if err != nil {
		return nil, errors.New("Error getting search results from Algolia")
	}
	defer resp.Body.Close()

	var parsed AlgoliaSearchResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&parsed)
	if err != nil {
		return nil, errors.New("Invalid JSON received from Algolia")
	}

	return &parsed, nil
}

func GetItem(ID string) (*AlgoliaItemResponse, error) {
	resp, err := http.Get(algoliaItemURL + ID)
	if err != nil {
		return nil, errors.New("Error getting item results from Algolia")
	}
	defer resp.Body.Close()

	var parsed AlgoliaItemResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&parsed)
	if err != nil {
		return nil, errors.New("Invalid JSON received from Algolia")
	}

	return &parsed, nil
}
