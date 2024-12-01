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

type monsterService struct {
	client *api.Client
	logger *slog.Logger

	cache map[string]oas.MonsterSchema
}

func newMonsterService(client *api.Client) *monsterService {
	s := &monsterService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "monsters"),

		cache: make(map[string]oas.MonsterSchema),
	}

	start := time.Now()
	s.sync()
	s.logger.Info("resources sync done: " + time.Since(start).String())

	return s
}

func (s *monsterService) sync() {
	page := 1
	for {
		apiRequestCount.Inc("mosnters")

		resp, err := s.client.GetAllMonstersMonstersGet(context.Background(), oas.GetAllMonstersMonstersGetParams{
			Page: oas.NewOptInt(page),
		})
		if err != nil {
			s.logger.With("error", err).Error("fail get all monsters")
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

func (s *monsterService) get(code string) (oas.MonsterSchema, error) {
	v, ok := s.cache[code]
	if !ok {
		return oas.MonsterSchema{}, fmt.Errorf("not found '%s' monster", code)
	}

	return v, nil
}
