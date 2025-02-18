package gcal

import (
  "context"

  "google.golang.org/api/option"
	calendar "google.golang.org/api/calendar/v3"

  "github.com/kott/cal-term/internal/auth"
)

type Service struct {
  calendarService *calendar.Service
}

func New() (*Service, error) {
  config, err := auth.GetConfigFromFile()
  if err != nil {
    return nil, err
  }
  token, err := auth.GetTokenFromFile()
  if err != nil {
    return nil, err
  }

  client := config.Client(context.Background(), token)
  service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))

  return &Service{
    calendarService: service,
  }, nil
}  

func (s *Service) FetchEvents() ([]*calendar.Event, error) {

  tf, err := localTimeFrame()
  if err != nil {
    return nil, err
  }

  events, err := s.calendarService.Events.List("primary").
      TimeMin(tf.Start).
      TimeMax(tf.End).
      SingleEvents(true).
      OrderBy("startTime").
      Do()

  if err != nil {
    return nil, err
  }

  filterFullDayEvents(events.Items)

  return events.Items, nil
}

