package game

import (
	"context"
	"log/slog"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type bank struct {
	client *api.Client
	logger *slog.Logger
}

func newBank(client *api.Client) *bank {
	return &bank{
		client: client,
		logger: slog.Default().With(ycloggingslog.Stream, "game").With("service", "bank"),
	}
}

func (b *bank) Items() map[string]int {
	result := map[string]int{}

	page := 1
	for {
		apiRequestCount.Inc("bank")

		resp, err := b.client.GetBankItemsMyBankItemsGet(context.Background(), oas.GetBankItemsMyBankItemsGetParams{
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
