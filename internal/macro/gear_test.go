package macro

import (
	"context"
	"os"
	"testing"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_compact(t *testing.T) {
	assert.Equal(t, []oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
	}, compact([]oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
	}, 2))

	assert.Equal(t, []oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
	}, compact([]oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
	}, 2))

	assert.Equal(t, []oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
	}, compact([]oas.ItemSchema{
		{Code: "a"},
		{Code: "a"},
		{Code: "a"},
		{Code: "b"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
		{Code: "c"},
	}, 4))
}

func TestAllWithBestGear(t *testing.T) {
	require.NoError(t, godotenv.Load("../../local.env"))

	client, err := api.New(os.Getenv("SERVER_URL"), os.Getenv("SERVER_TOKEN"))
	require.NoError(t, err)

	c := generic.NewCharacter("Emilia", client)
	c.Init()
	g := game.New(client)

	monsters, err := client.GetAllMonstersMonstersGet(context.Background(), oas.GetAllMonstersMonstersGetParams{
		Size: oas.NewOptInt(50),
	})
	require.NoError(t, err)

	for _, monster := range monsters.Data {
		if len(GetBestGearForMonster(c, g, monster.Code)) > 0 {
			t.Logf("%s (%d): true\n", monster.Code, monster.Level)
		} else {
			t.Logf("%s (%d): false\n", monster.Code, monster.Level)
		}
	}
}
