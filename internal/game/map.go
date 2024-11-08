package game

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type Point struct {
	X, Y int
	Name string
}

type mapService struct {
	client *api.Client
	logger *slog.Logger

	cache map[string]Point
}

func newMapService(client *api.Client) *mapService {
	m := &mapService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "map"),

		cache: make(map[string]Point),
	}

	m.init()
	return m
}

func (s *mapService) init() {
	page := 1
	for {
		apiRequestCount.Inc("maps")

		resp, err := s.client.GetAllMapsMapsGet(context.Background(), oas.GetAllMapsMapsGetParams{
			Page: oas.NewOptInt(page),
		})
		if err != nil {
			s.logger.With("error", err).Error("fail get all maps")
			continue
		}

		for _, m := range resp.Data {
			if m.Content.IsNull() {
				continue
			}

			s.cache[m.Content.MapContentSchema.Code] = Point{
				X:    m.X,
				Y:    m.Y,
				Name: m.Content.MapContentSchema.Code,
			}
		}

		if page >= resp.Pages.Value.Int {
			break
		}
		page++
	}
}

func (s *mapService) get(ctx context.Context, code string) (Point, error) {
	v, ok := s.cache[code]
	if !ok {
		return Point{}, fmt.Errorf("not found '%s' on map", code)
	}

	return v, nil
}
