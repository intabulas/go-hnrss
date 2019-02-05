package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	NSDublinCore = "http://purl.org/dc/elements/1.1/"
	NSAtom       = "http://www.w3.org/2005/Atom"
	SiteURL      = "https://hnrss.org"
)

type CDATA struct {
	Value string `xml:",cdata"`
}

func Timestamp(fmt string, input time.Time) string {
	switch fmt {
	case "rss":
		return input.Format(time.RFC1123Z)
	case "atom", "jsonfeed":
		return input.Format(time.RFC3339)
	case "http":
		return input.Format(http.TimeFormat)
	default:
		return input.Format(time.RFC1123Z)
	}
}

func UTCNow() time.Time {
	return time.Now().UTC()
}

func ParseRequest(c *gin.Context, sp *searchParams, op *outputParams) {
	err := c.ShouldBindQuery(sp)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if strings.Contains(sp.Query, " OR ") {
		sp.Query = strings.Replace(sp.Query, " OR ", " ", -1)

		fields := strings.Fields(sp.Query)
		q := make([]string, len(fields))
		for i, f := range fields {
			q[i] = fmt.Sprintf("\"%s\"", f)
		}
		sp.Query = strings.Join(q, " ")
		sp.OptionalWords = strings.Join(q, " ")
	}

	err = c.ShouldBindQuery(op)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	op.Format = c.GetString("format")
	op.SelfLink = SiteURL + c.Request.URL.String()
}
