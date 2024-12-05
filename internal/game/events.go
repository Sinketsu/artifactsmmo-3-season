package game

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type eventService struct {
	client *api.Client
	logger *slog.Logger

	events []oas.ActiveEventSchema
	mu     sync.Mutex // protects events
}

func newEventService(client *api.Client) *eventService {
	s := &eventService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "events"),
	}

	start := time.Now()
	s.sync()
	s.logger.Info("events sync done: " + time.Since(start).String())

	go s.update()

	return s
}

func (s *eventService) sync() {
	resp, err := s.client.GetAllActiveEventsEventsActiveGet(context.Background(), oas.GetAllActiveEventsEventsActiveGetParams{
		Size: oas.NewOptInt(50), // assume that active events count will be < 50 always (now limit is 8)
	})
	if err != nil {
		slog.With("error", err).Error("fail to list events")
		return
	}

	s.mu.Lock()
	for _, ev := range s.events {
		if !slices.ContainsFunc(resp.Data, func(e oas.ActiveEventSchema) bool {
			return e.Code == ev.Code
		}) {
			s.logger.Info(fmt.Sprintf("event ended: %s (%s)", ev.Name, ev.Code))
		}
	}

	eventsActive.ResetAll()
	for _, ev := range resp.Data {
		eventsActive.Set(1, ev.Code)

		if !slices.ContainsFunc(s.events, func(e oas.ActiveEventSchema) bool {
			return e.Code == ev.Code
		}) {
			s.logger.Info(fmt.Sprintf("event started: %s (%s)", ev.Name, ev.Code))
		}
	}

	s.events = resp.Data
	s.mu.Unlock()
}

func (s *eventService) update() {
	for range time.Tick(30 * time.Second) {
		s.sync()
	}
}

func (s *eventService) get(code string) (Point, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.events {
		if event.Code == code {
			return Point{
				Name: event.Map.Content.MapContentSchema.Code,
				X:    event.Map.X,
				Y:    event.Map.Y,
			}, nil
		}
	}

	return Point{}, fmt.Errorf("not found")
}

func (s *eventService) all() []oas.ActiveEventSchema {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.Clone(s.events)
}
