package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func fetchHiring(c *gin.Context, query string) {
	params := make(url.Values)
	if query != "" {
		params.Set("query", fmt.Sprintf("\"%s\"", query))
		params.Set("hitsPerPage", "1")
	}
	params.Set("tags", "story,author_whoishiring")

	results, err := GetResults(params)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	if len(results.Hits) < 1 {
		e := errors.New("no whoishiring stories found")
		c.Error(e)
		c.String(http.StatusBadGateway, e.Error())
		return
	}

	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment"
	if query != "" {
		sp.Filters = "parent_id=" + results.Hits[0].ObjectID
	} else {
		filters := make([]string, len(results.Hits))
		for i, hit := range results.Hits {
			filters[i] = "parent_id=" + hit.ObjectID
		}
		sp.Filters = strings.Join(filters, " OR ")
	}
	sp.SearchAttributes = "default"
	op.Title = results.Hits[0].Title
	op.Link = hackerNewsItemID + results.Hits[0].ObjectID

	Generate(c, &sp, &op)
}

// seekingEmployeesHandler ...
func seekingEmployeesHandler(c *gin.Context) {
	fetchHiring(c, "Ask HN: Who is hiring?")
}

// seekingEmployersHandler
func seekingEmployersHandler(c *gin.Context) {
	fetchHiring(c, "Ask HN: Who wants to be hired?")
}

// seekingFreelanceHandler
func seekingFreelanceHandler(c *gin.Context) {
	fetchHiring(c, "Ask HN: Freelancer? Seeking freelancer?")
}

// seekingAllHandler
func seekingAllHandler(c *gin.Context) {
	fetchHiring(c, "")
}
