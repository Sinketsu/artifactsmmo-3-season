package game

import (
	"context"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type Order struct {
	Id       string
	Quantity int
	Price    int
}

type geService struct {
	client *api.Client
	logger *slog.Logger

	orders map[string][]Order
	mu     sync.RWMutex // protects items
}

func newGEService(client *api.Client) *geService {
	s := &geService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "ge"),
	}

	start := time.Now()
	s.sync()
	s.logger.Info("ge sync done: " + time.Since(start).String())

	go s.update()

	return s
}

func (s *geService) sync() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders = s.actualOrders()
}

func (s *geService) update() {
	for range time.Tick(30 * time.Second) {
		s.sync()
	}
}

func (s *geService) actualOrders() map[string][]Order {
	result := map[string][]Order{}

	page := 1
	for {
		apiRequestCount.Inc("ge")

		resp, err := s.client.GetGeSellOrdersGrandexchangeOrdersGet(context.Background(), oas.GetGeSellOrdersGrandexchangeOrdersGetParams{
			Page: oas.NewOptInt(page),
		})
		if err != nil {
			s.logger.With("error", err).Error("fail get ge orders")
			continue
		}

		for _, data := range resp.Data {
			result[data.Code] = append(result[data.Code], Order{
				Id:       data.ID,
				Quantity: data.Quantity,
				Price:    data.Price,
			})
		}

		if page >= resp.Pages.Value.Int {
			break
		}
		page++
	}

	for item := range result {
		slices.SortFunc(result[item], func(a, b Order) int {
			return a.Price - b.Price
		})
	}

	return result
}

func (s *geService) Get(item string) []Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.orders[item]
}
