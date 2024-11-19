package game

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

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

	cache map[string][]Point
}

func newMapService(client *api.Client) *mapService {
	s := &mapService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "maps"),

		cache: make(map[string][]Point),
	}

	start := time.Now()
	s.sync()
	s.logger.Info("maps sync done: " + time.Since(start).String())

	return s
}

func (s *mapService) sync() {
	page := 1

pages:
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

			s.cache[m.Content.MapContentSchema.Code] = append(s.cache[m.Content.MapContentSchema.Code], Point{
				X:    m.X,
				Y:    m.Y,
				Name: m.Content.MapContentSchema.Code,
			})

			if m.Content.MapContentSchema.Type == string(oas.GetAllMapsMapsGetContentTypeResource) {
				apiRequestCount.Inc("resources")

				resp, err := s.client.GetResourceResourcesCodeGet(context.Background(), oas.GetResourceResourcesCodeGetParams{
					Code: m.Content.MapContentSchema.Code,
				})
				if err != nil {
					s.logger.With("error", err).Error("fail get resource")
					continue pages
				}

				schema, ok := resp.(*oas.ResourceResponseSchema)
				if !ok {
					s.logger.With("error", err).Error("fail get resource")
					continue pages
				}

				for _, drop := range schema.Data.Drops {
					if drop.Rate == 1 {
						s.cache[drop.Code] = append(s.cache[drop.Code], Point{
							X:    m.X,
							Y:    m.Y,
							Name: drop.Code,
						})
					}
				}
			}
		}

		if page >= resp.Pages.Value.Int {
			break
		}
		page++
	}
}

func (s *mapService) get(code string, closestTo Point) (Point, error) {
	v, ok := s.cache[code]
	if !ok || len(v) == 0 {
		return Point{}, fmt.Errorf("not found '%s' on map", code)
	}

	if len(v) == 1 {
		return v[0], nil
	}

	closest := v[0]
	distance := s.distance(v[0], closestTo)

	for _, p := range v[1:] {
		if d := s.distance(p, closestTo); d < distance {
			closest = p
			distance = d
		}
	}

	return closest, nil
}

func (s *mapService) distance(a, b Point) int {
	return int(math.Abs(float64(a.X)-float64(b.X)) + math.Abs(float64(a.Y)-float64(b.Y)))
}
