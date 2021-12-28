package main

import (
	"encoding/json"

	"fmt"

	"os"

	"net/http"

	"github.com/gocolly/colly"
)

type Info struct {
	StatusCode  int
	Name        string
	Title       string
	Description string
	Image       string
	Link        string
}

func handler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		fmt.Println("ERROR: Missing URL argument")
		return
	}
	fmt.Println("INFO: Visiting", URL)

	c := colly.NewCollector()

	data := Info{}

	// count links
	c.OnHTML("meta[property]", func(e *colly.HTMLElement) {
		property := e.Attr("property")
		content := e.Attr("content")
		switch property {
		case "og:site_name":
			data.Name = content

		case "og:title":
			data.Title = content

		case "og:description":
			data.Description = content

		case "og:image":
			data.Image = content

		case "og:url":
			data.Link = content
		}
		if data.Link == "" {
			data.Link = URL
		}
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("INFO: Response received", r.StatusCode)
		data.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("ERROR:", r.StatusCode, err)
		data.StatusCode = r.StatusCode
	})

	c.Visit(URL)

	// dump results
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("ERROR: Failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func GetOrigins() string {
	var origin = os.Getenv("ORIGIN_ALLOWED")
	if origin == "" {
		origin = "*"
		fmt.Println("INFO: No ORIGIN_ALLOWED environment variable detected, defaulting to " + origin)
	}
	return origin
}
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	// example usage: curl -s 'http://127.0.0.1:8080/?url=http://go-colly.org/'

	http.HandleFunc("/", handler)

	fmt.Println("INFO: Listening on port", GetPort())
	http.ListenAndServe(GetPort(), nil)
}
