package main

import (
	"encoding/json"
	"fmt"
	"github.com/eladmica/go-meetup/meetup"
	_ "github.com/joho/godotenv/autoload"
	"testing"
)

func generateEvents() []*meetup.Event {
	jsonString := []byte(`[
   {
    "created": 1678588826000,
    "duration": 86340000,
    "id": "123456789",
    "name": "My Event Day 1",
    "date_in_series_pattern": false,
    "status": "upcoming",
    "time": 1686261600000,
    "local_date": "2023-06-09",
    "local_time": "00:00",
    "updated": 1678588826000,
    "utc_offset": 7200000,
    "waitlist_count": 0,
    "yes_rsvp_count": 11,
    "is_online_event": true,
    "eventType": "ONLINE",
    "group": {
      "created": 1582809229000,
      "name": "Example",
      "id": 123456789,
      "join_mode": "open",
      "lat": 52.52000045776367,
      "lon": 13.380000114440918,
      "urlname": "example",
      "who": "Examplers",
      "localized_location": "Example, Country",
      "state": "",
      "country": "en",
      "region": "en_US",
      "timezone": "Europe/Berlin"
    },
    "link": "https://www.meetup.com/example/events/123456789/",
    "description": "<p>html descriptions are the best</p>",
    "visibility": "public",
    "member_pay_fee": false
  }
]`)
	var events []*meetup.Event
	_ = json.Unmarshal(jsonString, &events)
	return events
}

func TestParseFeedItemsFromEvents(t *testing.T) {
	events := ParseFeedItemsFromEvents(generateEvents())

	for _, event := range events {
		fmt.Printf("%s\n", event.Title)

	}

	if len(events) < 1 {
		t.Errorf("Expected at least 1 event, but got %d", len(events))
	}
}
