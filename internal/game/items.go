package game

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type items struct {
	client *api.Client
	logger *slog.Logger

	cache map[string]oas.ItemSchema
}

func newItems(client *api.Client) *items {
	s := &items{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "items"),

		cache: make(map[string]oas.ItemSchema),
	}

	s.init()
	return s
}

func (s *items) init() {
	page := 1
	for {
		apiRequestCount.Inc("items")

		resp, err := s.client.GetAllItemsItemsGet(context.Background(), oas.GetAllItemsItemsGetParams{
			Page: oas.NewOptInt(page),
		})
		if err != nil {
			s.logger.With("error", err).Error("fail get all maps")
			continue
		}

		for _, m := range resp.Data {
			s.cache[m.Code] = m
		}

		if page >= resp.Pages.Value.Int {
			break
		}
		page++
	}
}

func (s *items) Get(ctx context.Context, code string) (oas.ItemSchema, error) {
	v, ok := s.cache[code]
	if !ok {
		return oas.ItemSchema{}, fmt.Errorf("not found '%s' item", code)
	}

	return v, nil
}
