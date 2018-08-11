package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Dispatcher generates all the basic responses.
func Dispatcher(c *gin.Context) {
	var (
		sp SearchParams
		op OutputParams
	)

	// Set default tags, title, and link
	switch c.Request.URL.Path {
	}

	// Parse the search params
	err := c.ShouldBindQuery(&sp)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}

	// Tweak as needed if there was a query
	if sp.Query != "" {
		op.Title = fmt.Sprintf("%s (\"%s\")", op.Title, sp.Query)
	}

	// Needed to search comments
	if sp.Query != "" && c.Request.URL.Path == "/newcomments" {
		sp.SearchAttributes = "default"
	}

	// Make the request to Algolia
	results, err := GetResults(sp.Values())
	if err != nil {
		c.String(http.StatusBadRequest, err.Error()) // TODO(ejd): Bad Gateway instead?
	}
	c.Header("X-Algolia-URL", algoliaURL+sp.Values().Encode())

	// Parse the output params
	err = c.ShouldBindQuery(&op)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}

	// Generate the response
	op.Output(c, results)
}

func ParseRequest(c *gin.Context) (*SearchParams, *OutputParams) {
	var (
		sp SearchParams
		op OutputParams
	)

	err := c.ShouldBindQuery(&sp)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}

	err = c.ShouldBindQuery(&op)
	if err != nil {
		c.String(http.StatusBadRequest, "Error parsing the request")
	}

	return &sp, &op
}

func Generate(c *gin.Context, sp *SearchParams, op *OutputParams) {
	if op.Format == "" {
		op.Format = "rss"
	}

	results, err := GetResults(sp.Values())
	if err != nil {
		c.String(http.StatusBadGateway, err.Error()) // TODO(ejd): inspect error to know which HTTP type?
	}
	c.Header("X-Algolia-URL", algoliaURL+sp.Values().Encode())

	switch op.Format {
	case "rss":
		rss := NewRSS(results, op)
		c.XML(http.StatusOK, rss)
	}
}

func Newest(c *gin.Context) {
	sp, op := ParseRequest(c)

	sp.Tags = "(story,poll)"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Newest: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Newest"
	}
	op.Link = "https://news.ycombinator.com/newest"

	Generate(c, sp, op)
}

func Frontpage(c *gin.Context) {
	sp, op := ParseRequest(c)

	sp.Tags = "front_page"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Front Page: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Front Page"
	}
	op.Link = "https://news.ycombinator.com/"

	Generate(c, sp, op)
}

func Newcomments(c *gin.Context) {
	sp, op := ParseRequest(c)

	sp.Tags = "comment"
	if sp.Query != "" {
		sp.SearchAttributes = "default"
		op.Title = fmt.Sprintf("Hacker News - New Comments: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: New Comments"
	}
	op.Link = "https://news.ycombinator.com/newcomments"

	Generate(c, sp, op)
}

func AskHN(c *gin.Context) {
	sp, op := ParseRequest(c)

	sp.Tags = "ask_hn"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Ask HN: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Ask HN"
	}
	op.Link = "https://news.ycombinator.com/ask"

	Generate(c, sp, op)
}

func ShowHN(c *gin.Context) {
	sp, op := ParseRequest(c)

	sp.Tags = "show_hn"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Show HN: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Show HN"
	}
	op.Link = "https://news.ycombinator.com/shownew"

	Generate(c, sp, op)
}

func Polls(c *gin.Context) {
	sp, op := ParseRequest(c)

	sp.Tags = "poll"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Polls: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Polls"
	}
	op.Link = "https://news.ycombinator.com/"

	Generate(c, sp, op)
}

func Jobs(c *gin.Context) {
	sp, op := ParseRequest(c)

	sp.Tags = "job"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Jobs: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Jobs"
	}
	op.Link = "https://news.ycombinator.com/jobs"

	Generate(c, sp, op)
}
