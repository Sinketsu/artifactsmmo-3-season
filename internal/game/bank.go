package game

import (
	"context"
	"log/slog"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type bankService struct {
	client *api.Client
	logger *slog.Logger
}

func newBankService(client *api.Client) *bankService {
	s := &bankService{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "bank"),
	}
	go s.update()

	return s
}

func (s *bankService) update() {
	for range time.Tick(30 * time.Second) {
		resp, err := s.client.GetBankDetailsMyBankGet(context.Background())
		if err != nil {
			s.logger.With("error", err).Error("fail to get bank details")
		}

		bankGoldCount.Set(int64(resp.Data.Gold))
	}
}

func (b *bankService) Items(ctx context.Context) map[string]int {
	result := map[string]int{}

	page := 1
	for {
		apiRequestCount.Inc("bank")

		resp, err := b.client.GetBankItemsMyBankItemsGet(ctx, oas.GetBankItemsMyBankItemsGetParams{
			Page: oas.NewOptInt(page),
		})
		if err != nil {
			b.logger.With("error", err).Error("fail get bank items")
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
