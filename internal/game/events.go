package game

import (
	"context"
	"fmt"
	"log/slog"
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
	mu     sync.Mutex
}

func newEventService(client *api.Client) *eventService {
	s := &eventService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "events"),
	}

	s.update()
	go s.update()

	return s
}

func (s *eventService) update() {
	for range time.Tick(30 * time.Second) {
		resp, err := s.client.GetAllActiveEventsEventsActiveGet(context.Background(), oas.GetAllActiveEventsEventsActiveGetParams{})
		if err != nil {
			slog.With("error", err).Error("fail to list events")
			continue
		}

		s.mu.Lock()
		s.events = resp.Data
		s.mu.Unlock()
	}
}

func (s *eventService) get(code string) (Point, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.events {
		if event.Code == code {
			return Point{
				Name: event.Code,
				X:    event.Map.X,
				Y:    event.Map.Y,
			}, nil
		}
	}

	return Point{}, fmt.Errorf("not found")
}
