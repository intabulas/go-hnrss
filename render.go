package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderResults(c *gin.Context, sp *searchParams, op *outputParams) {
	if op.Format == "" {
		op.Format = "rss"
	}

	results, err := GetResults(sp.Values())
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.Header("X-Algolia-URL", algoliaSearchURL+sp.Values().Encode())

	if len(results.Hits) > 0 {
		item := results.Hits[0]

		recent := item.GetCreatedAt()
		c.Header("Last-Modified", Timestamp("http", recent))

		if c.Request.URL.Path == "/item" {
			if sp.Query != "" {
				op.Title = fmt.Sprintf("Hacker News - \"%s\": \"%s\"", item.StoryTitle, sp.Query)
			} else {
				op.Title = fmt.Sprintf("Hacker News: New comments on \"%s\"", item.StoryTitle)
			}
		}
	}

	switch op.Format {
	case "rss":
		c.XML(http.StatusOK, resultsAsRSS(results, op))
	case "atom":
		c.XML(http.StatusOK, resultsAsAtom(results, op))
	case "jsonfeed":
		c.JSON(http.StatusOK, resultsAsJSONFeed(results, op))
	}
}
