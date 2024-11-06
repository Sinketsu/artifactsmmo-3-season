package game

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

var (
	// TODO add more common macros
	GrandExchange = Point{Name: "grand_exchange"}
	Bank          = Point{Name: "bank"}
)

type Point struct {
	X, Y int
	Name string
}

type maps struct {
	client *api.Client
	logger *slog.Logger

	cache map[string]Point
}

func newMaps(client *api.Client) *maps {
	m := &maps{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "map"),

		cache: make(map[string]Point),
	}

	m.init()
	return m
}

func (s *maps) init() {
	page := 1
	for {
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

			switch m.Content.MapContentSchema.Code {
			case GrandExchange.Name:
				GrandExchange.X, GrandExchange.Y = m.X, m.Y
			case Bank.Name:
				Bank.X, Bank.Y = m.X, m.Y
			}
		}

		if page >= resp.Pages.Value.Int {
			break
		}
		page++
	}
}

func (m *maps) Get(ctx context.Context, code string) (Point, error) {
	mapRequestRate.Inc()

	v, ok := m.cache[code]
	if !ok {
		return Point{}, fmt.Errorf("not found '%s' on map", code)
	}

	return v, nil
}
