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

type resourceService struct {
	client *api.Client
	logger *slog.Logger

	cache map[string]oas.ResourceSchema
}

func newResourceService(client *api.Client) *resourceService {
	s := &resourceService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "resources"),

		cache: make(map[string]oas.ResourceSchema),
	}

	start := time.Now()
	s.sync()
	s.logger.Info("resources sync done: " + time.Since(start).String())

	return s
}

func (s *resourceService) sync() {
	page := 1
	for {
		apiRequestCount.Inc("resources")

		resp, err := s.client.GetAllResourcesResourcesGet(context.Background(), oas.GetAllResourcesResourcesGetParams{
			Page: oas.NewOptInt(page),
		})
		if err != nil {
			s.logger.With("error", err).Error("fail get all resources")
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

func (s *resourceService) get(code string) (oas.ResourceSchema, error) {
	v, ok := s.cache[code]
	if !ok {
		return oas.ResourceSchema{}, fmt.Errorf("not found '%s' resource", code)
	}

	return v, nil
}
