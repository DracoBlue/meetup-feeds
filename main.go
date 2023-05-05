package main

import (
	"fmt"
	"github.com/eladmica/go-meetup/meetup"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/feeds"
	godotenv "gopkg.in/joho/godotenv.v1"
	"html"
	"log"
	"net/http"
	"os"
	"time"
)

func ParseFeedItemsFromEvents(events []*meetup.Event) []*feeds.Item {
	var items []*feeds.Item

	for _, event := range events {
		eventTime := time.Unix(int64(event.Time/1000), 0)
		updateTime := time.Unix(int64(event.Updated/1000), 0)

		venueSuffix := ""

		if event.Venue != nil {
			venueSuffix = "<p>"
			venueSuffix += html.EscapeString(event.Venue.Name)
			if event.Venue.Address1 != "" {
				venueSuffix += html.EscapeString(", " + event.Venue.Address1)
			}
			if event.Venue.Address2 != "" {
				venueSuffix += html.EscapeString(", " + event.Venue.Address2)
			}
			if event.Venue.City != "" {
				venueSuffix += html.EscapeString(", " + event.Venue.City)
			}
			if event.Venue.LocalizedCountryName != "" {
				venueSuffix += html.EscapeString(", " + event.Venue.LocalizedCountryName)
			}
			venueSuffix += "</p>"
		} else {
			venueSuffix = "<p>No venue</p>"
		}

		items = append(items, &feeds.Item{
			Title:       event.Name + " (" + eventTime.Format("Mon, 2 Jan 2006 15:04:05 MST") + ")",
			Link:        &feeds.Link{Href: event.Link},
			Description: event.Description + venueSuffix + "<p>" + eventTime.Format("Mon, 2 Jan 2006 15:04:05 MST") + "</p>",
			Created:     updateTime,
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
		Link:        &feeds.Link{Href: "https://meetup.com/" + groupUrlName + "/events/"},
		Description: "",
		Author:      nil,
		Created:     now,
		Items:       ParseFeedItemsFromEvents(events),
	}

	if len(events) > 0 {
		feed.Title = events[0].Group.Name
		feed.Description = events[0].Group.Name
	}

	if feedType == "rss" {
		rss, err := feed.ToRss()
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(fmt.Sprintf("%s", rss)))
		return
	}

	if feedType == "atom" {
		atom, err := feed.ToAtom()
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/atom+xml")
		w.Write([]byte(fmt.Sprintf("%s", atom)))
		return
	}

	panic("Invalid format")
}

func main() {
	godotenv.Load()

	_, userNameExists := os.LookupEnv("USERNAME")

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	if userNameExists {
		router.Use(middleware.BasicAuth("MyRealm", map[string]string{
			os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
		}))
	}

	router.Get("/{type:rss}/{groupUrlName:[a-z-]+}", feedHandler)
	router.Get("/{type:atom}/{groupUrlName:[a-z-]+}", feedHandler)

	port, portIsSet := os.LookupEnv("PORT")

	if !portIsSet {
		port = "8090"
	}

	http.ListenAndServe(":"+port, router)
}
