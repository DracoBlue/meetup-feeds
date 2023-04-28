package main

import (
	"fmt"
	"github.com/eladmica/go-meetup/meetup"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/feeds"
	godotenv "gopkg.in/joho/godotenv.v1"
	"log"
	"net/http"
	"os"
	"time"
)

func ParseFeedItemsFromEvents(events []*meetup.Event) []*feeds.Item {
	var items []*feeds.Item

	for _, event := range events {
		eventTime := time.Unix(int64(event.Time/1000), 0)

		items = append(items, &feeds.Item{
			Title:       event.Name,
			Link:        &feeds.Link{Href: event.Link},
			Description: event.Description,
			Author:      &feeds.Author{Name: event.Group.Name},
			Created:     eventTime,
		})

	}

	return items
}

func feedHandler(w http.ResponseWriter, req *http.Request) {
	feedType := chi.URLParam(req, "type")
	groupUrlName := chi.URLParam(req, "groupUrlName")

	client := meetup.NewClient(nil)

	events, err := client.GetEvents(groupUrlName, nil)
	if err != nil {
		panic(err)
	}

	now := time.Now()

	feed := &feeds.Feed{
		Title:       groupUrlName,
		Link:        &feeds.Link{Href: "https://meetup.com/" + groupUrlName},
		Description: "",
		Author:      nil,
		Created:     now,
		Items:       ParseFeedItemsFromEvents(events),
	}

	if feedType == "rss" {
		rss, err := feed.ToRss()
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("%s", rss)))
		return
	}

	if feedType == "atom" {
		atom, err := feed.ToAtom()
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("%s", atom)))
		return
	}

	panic("Invalid format")
}

func ensureEnvironmentVariableIsSet(key string) {
	_, exists := os.LookupEnv(key)
	if !exists {
		panic("Environment variable " + key + " is missing!")
	}
}

func main() {
	godotenv.Load()

	ensureEnvironmentVariableIsSet("USERNAME")
	ensureEnvironmentVariableIsSet("PASSWORD")

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.BasicAuth("MyRealm", map[string]string{
		os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
	}))

	router.Get("/{type:rss}/{groupUrlName:[a-z-]+}", feedHandler)
	router.Get("/{type:atom}/{groupUrlName:[a-z-]+}", feedHandler)

	port, portIsSet := os.LookupEnv("PORT")

	if !portIsSet {
		port = "8090"
	}

	http.ListenAndServe(":"+port, router)
}
