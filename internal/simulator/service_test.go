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

	for _, code := range []string{
		"bloodblade",
		"celest_ring",
		"celest_ring",
		"magic_stone_amulet",
		"cultist_boots",
		"cursed_hat",
		"fire_shield",
		"mithril_platebody",
		"mithril_platelegs",
		// artifacts
		"malefic_crystal",
		"christmas_star",
		"life_crystal",
		// elixirs
		"fire_res_potion",
		"enchanted_boost_potion",
	} {
		item, err := client.GetItemItemsCodeGet(context.Background(), oas.GetItemItemsCodeGetParams{
			Code: code,
		})
		require.NoError(t, err)

		items = append(items, item.(*oas.ItemResponseSchema).Data)
	}

	monster, err := client.GetMonsterMonstersCodeGet(context.Background(), oas.GetMonsterMonstersCodeGetParams{
		Code: "rosenblood",
	})
	require.NoError(t, err)

	result := sim.Fight(&characterStub{level: 40}, items, monster.(*oas.MonsterResponseSchema).Data)
	t.Log("monster:", monster.(*oas.MonsterResponseSchema).Data.Code)
	t.Log("win:", result.Win)
	t.Log("turns:", result.Turns)
	t.Log("seconds:", result.Seconds)
	t.Log("remainig monster hp:", result.RemainingMonsterHp)
	t.Log("remainig character hp:", result.RemainingCharacterHp)
	t.Log("")

	items = []oas.ItemSchema{}
	for _, code := range []string{
		"lightning_sword",
		"eternity_ring",
		"eternity_ring",
		"greater_emerald_amulet",
		"mithril_boots",
		"white_knight_helmet",
		"mithril_shield",
		"white_knight_armor",
		"white_knight_pants",
		// artifacts
		"malefic_crystal",
		"christmas_star",
		"life_crystal",
		// elixirs
		"fire_res_potion",
		"air_res_potion",
	} {
		item, err := client.GetItemItemsCodeGet(context.Background(), oas.GetItemItemsCodeGetParams{
			Code: code,
		})
		require.NoError(t, err)

		items = append(items, item.(*oas.ItemResponseSchema).Data)
	}

	monster, err = client.GetMonsterMonstersCodeGet(context.Background(), oas.GetMonsterMonstersCodeGetParams{
		Code: "efreet_sultan",
	})
	require.NoError(t, err)

	result = sim.Fight(&characterStub{level: 40}, items, monster.(*oas.MonsterResponseSchema).Data)
	t.Log("monster:", monster.(*oas.MonsterResponseSchema).Data.Code)
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
