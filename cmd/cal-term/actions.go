package main

import (
	"fmt"

	"github.com/kott/cal-term/internal/auth"
	"github.com/kott/cal-term/internal/gcal"
	"github.com/kott/cal-term/internal/render"
)

func authAction(clientId, clientSecret string) error {
	if clientId == "" || clientSecret == "" {
		return fmt.Errorf("clientID and clientSecret cannot be empty")
	}

	a := auth.New(clientId, clientSecret)
	token, err := a.GetTokenFromWeb()
	if err != nil {
		return err
	}

	if err := a.StoreToken(token); err != nil {
		return err
	}

	if err := a.StoreCredentials(); err != nil {
		return err
	}

	return nil
}

func viewAction() error {
	s, err := gcal.New()
	if err != nil {
		return err
	}

	events, err := s.FetchEvents()
	if err != nil {
		return err
	}

	render.DisplayEvents(events)
	return nil
}
