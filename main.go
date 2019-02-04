package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	bindAddr    = flag.String("bind", "127.0.0.1:9000", "HOST:PORT")
	buildString string
)

func registerEndpoint(r *gin.Engine, url string, fn gin.HandlerFunc) {
	r.GET(url, SetFormat("rss"), fn)
	r.GET(url+".jsonfeed", SetFormat("jsonfeed"), fn)
	r.GET(url+".atom", SetFormat("atom"), fn)
}

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	registerEndpoint(r, "/newest", newestPostHandler)
	registerEndpoint(r, "/frontpage", frontpageHandler)
	registerEndpoint(r, "/newcomments", newCommentsHandler)
	registerEndpoint(r, "/ask", askHNHandler)
	registerEndpoint(r, "/show", showHNHandler)
	registerEndpoint(r, "/polls", pollsHandler)
	registerEndpoint(r, "/jobs", jobsHandler)
	registerEndpoint(r, "/user", userAllHandler)
	registerEndpoint(r, "/threads", userThreadsHandler)
	registerEndpoint(r, "/submitted", userSubmittedHandler)
	registerEndpoint(r, "/replies", repliesHandler)
	registerEndpoint(r, "/item", itemHandler)
	registerEndpoint(r, "/whoishiring/jobs", seekingEmployeesHandler)
	registerEndpoint(r, "/whoishiring/hired", seekingEmployersHandler)
	registerEndpoint(r, "/whoishiring/freelance", seekingFreelanceHandler)
	registerEndpoint(r, "/whoishiring", seekingAllHandler)

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://news.ycombinator.com/favicon.ico")
	})
	r.GET("/robots.txt", func(c *gin.Context) {
		c.String(http.StatusOK, "User-agent: *\nDisallow:\n")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://edavis.github.io/hnrss/")
	})

	flag.Parse()

	srv := &http.Server{
		Addr:    *bindAddr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %s\n", err)
	}
	log.Println("Server exiting")
}
