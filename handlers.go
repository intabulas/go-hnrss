package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// newestPostHandler ..
func newestPostHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "(story,poll)"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Newest: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Newest"
	}
	op.Link = "https://news.ycombinator.com/newest"

	renderResults(c, &sp, &op)
}

// frontpageHandler
func frontpageHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "front_page"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Front Page: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Front Page"
	}
	op.Link = "https://news.ycombinator.com/"

	renderResults(c, &sp, &op)
}

// newCommentsHandler
func newCommentsHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment"
	if sp.Query != "" {
		sp.SearchAttributes = "default"
		op.Title = fmt.Sprintf("Hacker News - New Comments: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: New Comments"
	}
	op.Link = "https://news.ycombinator.com/newcomments"

	renderResults(c, &sp, &op)
}

// askHNHandler
func askHNHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "ask_hn"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Ask HN: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Ask HN"
	}
	op.Link = "https://news.ycombinator.com/ask"

	renderResults(c, &sp, &op)
}

// showHNHandler
func showHNHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "show_hn"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Show HN: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Show HN"
	}
	op.Link = "https://news.ycombinator.com/shownew"

	renderResults(c, &sp, &op)
}

// pollsHandler
func pollsHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "poll"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Polls: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Polls"
	}
	op.Link = "https://news.ycombinator.com/"

	renderResults(c, &sp, &op)
}

// jobsHandler
func jobsHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "job"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Jobs: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Jobs"
	}
	op.Link = "https://news.ycombinator.com/jobs"

	renderResults(c, &sp, &op)
}

// userAllHandler
func userAllHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	tags := []string{"(story,comment,poll)", "author_" + sp.ID}
	sp.Tags = strings.Join(tags, ",")

	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - %s: \"%s\"", sp.ID, sp.Query)
	} else {
		op.Title = fmt.Sprintf("Hacker News: %s", sp.ID)
	}
	op.Link = "https://news.ycombinator.com/user?id=" + sp.ID

	renderResults(c, &sp, &op)
}

// userThreadsHandler
func userThreadsHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	tags := []string{"comment", "author_" + sp.ID}
	sp.Tags = strings.Join(tags, ",")

	if sp.Query != "" {
		sp.SearchAttributes = "default"
		op.Title = fmt.Sprintf("Hacker News - %s threads: \"%s\"", sp.ID, sp.Query)
	} else {
		op.Title = fmt.Sprintf("Hacker News: %s threads", sp.ID)
	}
	op.Link = "https://news.ycombinator.com/threads?id=" + sp.ID

	renderResults(c, &sp, &op)
}

// userSubmittedHandler
func userSubmittedHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	tags := []string{"(story,poll)", "author_" + sp.ID}
	sp.Tags = strings.Join(tags, ",")

	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - %s submitted: \"%s\"", sp.ID, sp.Query)
	} else {
		op.Title = fmt.Sprintf("Hacker News: %s submitted", sp.ID)
	}
	op.Link = "https://news.ycombinator.com/submitted?id=" + sp.ID

	renderResults(c, &sp, &op)
}

// repliesHandler
func repliesHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment"
	sp.SearchAttributes = "default"

	// If ID is a number, look for comments with a parent_id equal to the ID.
	// If ID is not a number, assume it is an author username and grab replies to their comments.
	_, err := strconv.Atoi(sp.ID)
	if err == nil {
		sp.Filters = "parent_id=" + sp.ID
		op.Title = "Hacker News: Replies to item #" + sp.ID
		op.Link = "https://news.ycombinator.com/item?id=" + sp.ID
	} else {
		values := make(url.Values)
		values.Set("tags", "comment,author_"+sp.ID)
		results, err := GetResults(values)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadGateway, err.Error())
			return
		}

		filters := make([]string, len(results.Hits))
		for i, hit := range results.Hits {
			filters[i] = "parent_id=" + hit.ObjectID
		}

		sp.Filters = strings.Join(filters, " OR ")
		op.Title = "Hacker News: Replies to " + sp.ID
		op.Link = "https://news.ycombinator.com/threads?id=" + sp.ID
	}

	renderResults(c, &sp, &op)
}

// itemHandler
func itemHandler(c *gin.Context) {
	var sp searchParams
	var op outputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment,story_" + sp.ID
	sp.SearchAttributes = "default"

	// op.Title is set inside renderResults to avoid the overhead of a
	// separate HTTP request to obtain the title.
	op.Link = "https://news.ycombinator.com/item?id=" + sp.ID

	renderResults(c, &sp, &op)
}
