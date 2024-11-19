package game

import (
	"context"
	"log/slog"
	"maps"
	"sync"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type bankService struct {
	client *api.Client
	logger *slog.Logger

	items map[string]int
	mu    sync.Mutex // protects items
}

func newBankService(client *api.Client) *bankService {
	s := &bankService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "bank"),
	}

	start := time.Now()
	s.sync()
	s.logger.Info("bank sync done: " + time.Since(start).String())

	go s.update()

	return s
}

func (s *bankService) sync() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = s.actualItems()

	bankItemCount.ResetAll()
	for item, quantity := range s.items {
		bankItemCount.Set(int64(quantity), item)
	}

	resp, err := s.client.GetBankDetailsMyBankGet(context.Background())
	if err != nil {
		s.logger.With("error", err).Error("fail to get bank details")
		return
	}

	bankGoldCount.Set(int64(resp.Data.Gold))
}

func (s *bankService) update() {
	for range time.Tick(30 * time.Second) {
		s.sync()
	}
}

func (s *bankService) actualItems() map[string]int {
	result := map[string]int{}

	page := 1
	for {
		apiRequestCount.Inc("bank")

		resp, err := s.client.GetBankItemsMyBankItemsGet(context.Background(), oas.GetBankItemsMyBankItemsGetParams{
			Page: oas.NewOptInt(page),
		})
		if err != nil {
			s.logger.With("error", err).Error("fail get bank items")
			continue
		}

		for _, data := range resp.Data {
			result[data.Code] = data.Quantity
		}

		if page >= resp.Pages.Value.Int {
			break
		}
		page++
	}

	return result
}

func (s *bankService) Items() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return maps.Clone(s.items)
}
