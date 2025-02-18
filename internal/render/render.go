package render

import (
  "fmt"

	calendar "google.golang.org/api/calendar/v3"
)

func DisplayEvents(events []*calendar.Event) {
	fmt.Println("Upcoming events:")
	if len(events) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events {
			date := item.Start.DateTime
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}

