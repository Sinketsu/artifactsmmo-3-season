package game

import (
	"context"
	"log/slog"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
	"github.com/jellydator/ttlcache/v3"
)

type Point struct {
	X, Y int
	Name string
}

type maps struct {
	client *api.Client
	logger *slog.Logger

	cache *ttlcache.Cache[string, Point]
}

func newMaps(client *api.Client) *maps {
	m := &maps{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "map"),

		cache: ttlcache.New[string, Point](),
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

			s.cache.Set(
				m.Content.MapContentSchema.Code,
				Point{
					X:    m.X,
					Y:    m.Y,
					Name: m.Content.MapContentSchema.Code,
				},
				ttlcache.NoTTL,
			)
		}

		if page >= resp.Pages.Value.Int {
			break
		}
		page++
	}
	// TODO - maybe ignore event monsters?
}

func (m *maps) Get(code string) Point {
	// TODO
	return Point{}
}
