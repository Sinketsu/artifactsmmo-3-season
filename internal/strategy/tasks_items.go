package strategy

import (
	"context"
	"fmt"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
	"golang.org/x/exp/maps"
)

type tasksItems struct {
	character *generic.Character
	game      *game.Game
}

func TasksItems(character *generic.Character, game *game.Game) *tasksItems {
	return &tasksItems{
		character: character,
		game:      game,
	}
}

func (s *tasksItems) Name() string {
	return "do items tasks"
}

func (s *tasksItems) Do(ctx context.Context) error {
	task := s.character.Task()

	if task == generic.NoTask {
		macro.AcceptItemTask(ctx, s.character, s.game)

		task = s.character.Task()
		if task == generic.NoTask {
			// if some error occured
			return fmt.Errorf("no task...")
		}
	}

	if task.Current == task.Total {
		macro.CompleteTask(ctx, s.character, s.game)

		macro.Deposit(ctx, s.character, s.game, maps.Keys(s.character.Inventory())...)

		return nil
	}

	if s.character.Inventory()[task.Code] >= task.Total-task.Current {
		macro.TradeItemTask(ctx, s.character, s.game, task.Code, task.Total-task.Current)

		return nil
	}

	craft := ""
	goal, err := s.game.GetItem(task.Code)
	if err != nil {
		return fmt.Errorf("fail to get item %s: %w", task.Code, err)
	}

	// TODO good check craftable or not - now it is fail(
	if goal.Craft.IsSet() && len(goal.Craft.Value.CraftSchema.Items) > 0 {
		// TODO support multi res items
		craft = goal.Code
		code := goal.Craft.Value.CraftSchema.Items[0].Code

		goal, err = s.game.GetItem(code)
		if err != nil {
			return fmt.Errorf("fail to get item %s: %w", code, err)
		}
	}

	if s.character.InventoryFull() {
		if craft != "" {
			macro.CraftFromInventory(ctx, s.character, s.game, craft)
		}

		macro.TradeItemTask(ctx, s.character, s.game, task.Code, min(s.character.Inventory()[task.Code], task.Total-task.Current))

		return nil
	}

	spot, err := s.game.Find(goal.Code, s.character.Location())
	if err != nil {
		return fmt.Errorf("not found: %s", goal.Code)
	}

	err = s.character.Move(ctx, spot)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = s.character.Gather(ctx)
	if err != nil {
		return fmt.Errorf("gather: %w", err)
	}

	return nil
}
