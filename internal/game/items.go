package game

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type itemService struct {
	client *api.Client
	logger *slog.Logger

	cache map[string]oas.ItemSchema
}

func newItemService(client *api.Client) *itemService {
	s := &itemService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "items"),

		cache: make(map[string]oas.ItemSchema),
	}

	start := time.Now()
	s.sync()
	s.logger.Info("items sync done: " + time.Since(start).String())

	return s
}

func (s *itemService) sync() {
	page := 1
	for {
		apiRequestCount.Inc("items")

		resp, err := s.client.GetAllItemsItemsGet(context.Background(), oas.GetAllItemsItemsGetParams{
			Page: oas.NewOptInt(page),
		})
		if err != nil {
			s.logger.With("error", err).Error("fail get all items")
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

func (s *itemService) get(code string) (oas.ItemSchema, error) {
	v, ok := s.cache[code]
	if !ok {
		return oas.ItemSchema{}, fmt.Errorf("not found '%s' item", code)
	}

	return v, nil
}
