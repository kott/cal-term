package gcal

import (
	"time"

	calendar "google.golang.org/api/calendar/v3"
)

type timeFrame struct {
	Start string
	End   string
}

func localTimeFrame() (*timeFrame, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)

	timeMin := startOfDay.Format(time.RFC3339)
	timeMax := endOfDay.Format(time.RFC3339)

	return &timeFrame{
		Start: timeMin,
		End:   timeMax,
	}, nil
}

func filterFullDayEvents(events []*calendar.Event) {
	var filteredItems []*calendar.Event
	for _, item := range events {
		if item.Start.DateTime != "" {
			filteredItems = append(filteredItems, item)
		}
	}
	events = filteredItems
}
