package simulator

import (
	"context"
	"os"
	"testing"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

type characterStub struct {
	level int
}

func (cs *characterStub) Level() int {
	return cs.level
}

func TestSimple(t *testing.T) {
	require.NoError(t, godotenv.Load("../../local.env"))

	client, err := api.New(os.Getenv("SERVER_URL"), os.Getenv("SERVER_TOKEN"))
	require.NoError(t, err)

	sim := New()

	items := []oas.ItemSchema{}

	for _, code := range []string{"death_knight_sword", "steel_shield", "piggy_helmet", "skeleton_armor",
		"skeleton_pants", "steel_boots", "forest_ring", "forest_ring", "skull_amulet"} {
		item, err := client.GetItemItemsCodeGet(context.Background(), oas.GetItemItemsCodeGetParams{
			Code: code,
		})
		require.NoError(t, err)

		items = append(items, item.(*oas.ItemResponseSchema).Data)
	}

	monster, err := client.GetMonsterMonstersCodeGet(context.Background(), oas.GetMonsterMonstersCodeGetParams{
		Code: "owlbear",
	})
	require.NoError(t, err)

	result := sim.Fight(&characterStub{level: 35}, items, monster.(*oas.MonsterResponseSchema).Data)
	t.Log("win:", result.Win)
	t.Log("turns:", result.Turns)
	t.Log("seconds:", result.Seconds)
	t.Log("remainig monster hp:", result.RemainingMonsterHp)
	t.Log("remainig character hp:", result.RemainingCharacterHp)
}

func TestAll(t *testing.T) {
	require.NoError(t, godotenv.Load("../../local.env"))

	client, err := api.New(os.Getenv("SERVER_URL"), os.Getenv("SERVER_TOKEN"))
	require.NoError(t, err)

	sim := New()

	items := []oas.ItemSchema{}
	for _, code := range []string{"elderwood_staff", "steel_shield", "piggy_helmet", "skeleton_armor",
		"skeleton_pants", "steel_boots", "forest_ring", "forest_ring", "skull_amulet"} {
		item, err := client.GetItemItemsCodeGet(context.Background(), oas.GetItemItemsCodeGetParams{
			Code: code,
		})
		require.NoError(t, err)

		items = append(items, item.(*oas.ItemResponseSchema).Data)
	}

	monsters, err := client.GetAllMonstersMonstersGet(context.Background(), oas.GetAllMonstersMonstersGetParams{
		Size: oas.NewOptInt(50),
	})
	require.NoError(t, err)

	for _, monster := range monsters.Data {
		result := sim.Fight(&characterStub{level: 33}, items, monster)
		t.Logf("%s (%d): %t\n", monster.Code, monster.Level, result.Win)
	}
}
