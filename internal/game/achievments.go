package game

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type achievment struct {
	Name    string
	Total   int
	Current int
}

type achievmentService struct {
	account string
	client  *api.Client
	logger  *slog.Logger

	achievments map[string]achievment
	mu          sync.Mutex // protects achievments
}

func newAchievmentService(client *api.Client, account string) *achievmentService {
	s := &achievmentService{
		account:     account,
		client:      client,
		logger:      slog.Default().With(ycloggingslog.Stream, "game").With("service", "achievments"),
		achievments: make(map[string]achievment),
	}

	s.sync()
	go s.update()

	return s
}

func (s *achievmentService) sync() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, a := range s.actualAchievments() {
		s.achievments[a.Name] = a

		achievmentsStatus.Set(float64(a.Current)/float64(a.Total), a.Name)
	}
}

func (s *achievmentService) update() {
	for range time.Tick(30 * time.Second) {
		s.sync()
	}
}

func (s *achievmentService) get(name string) achievment {
	s.mu.Lock()
	defer s.mu.Unlock()

	if a, ok := s.achievments[name]; ok {
		return a
	}

	s.logger.Warn("requested achievment that not exists: " + name)
	return achievment{}
}

func (s *achievmentService) actualAchievments() []achievment {
	page := 1
	result := make([]achievment, 0)

	for {
		apiRequestCount.Inc("achievments")

		resp, err := s.client.GetAccountAchievementsAccountsAccountAchievementsGet(context.Background(), oas.GetAccountAchievementsAccountsAccountAchievementsGetParams{
			Account: s.account,
			Page:    oas.NewOptInt(page),
		})
		if err != nil {
			s.logger.With("error", err).Error("fail get achievments")
			continue
		}

		if data, ok := resp.(*oas.DataPageAccountAchievementSchema); ok {
			for _, a := range data.Data {
				result = append(result, achievment{
					Name:    a.Name,
					Total:   a.Total,
					Current: a.Current,
				})
			}

			if page >= data.Pages.Value.Int {
				break
			}
			page++
		} else {
			s.logger.With("response", resp).Error("fail get achievments")
			continue
		}
	}

	return result
}
