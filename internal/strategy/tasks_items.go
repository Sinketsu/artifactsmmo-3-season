package strategy

import (
	"context"
	"fmt"
	"slices"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
)

type tasksItems struct {
	character *generic.Character
	game      *game.Game

	cancel  []string
	events  []string
	current string
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

func (s *tasksItems) Cancel(tasks ...string) *tasksItems {
	s.cancel = append(s.cancel, tasks...)
	return s
}

func (s *tasksItems) AllowEvents(events ...string) *tasksItems {
	s.events = append(s.events, events...)
	return s
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

	if slices.Contains(s.cancel, task.Code) {
		macro.CancelItemsTask(ctx, s.character, s.game)
		return nil
	}

	if task.Current == task.Total {
		macro.CompleteItemTask(ctx, s.character, s.game)

		macro.Deposit(ctx, s.character, s.game, "tasks_coin")

		if s.character.Inventory()["tasks_coin"] > 5 {
			if err := s.character.Deposit(ctx, "tasks_coin", s.character.Inventory()["tasks_coin"]-5); err != nil {
				return err
			}
		}

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
		// TODO - may be support multi res items?
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

		macro.Deposit(ctx, s.character, s.game, "tasks_coin")
		return nil
	}

	var spot game.Point
	for _, event := range s.events {
		if event, err := s.game.GetEvent(event); err == nil {
			spot = event
			break
		}
	}

	if spot.Name == "" {
		spot, err = s.game.Find(goal.Code, s.character.Location())
		if err != nil {
			return fmt.Errorf("not found: %s", goal.Code)
		}
	}

	if s.current != spot.Name {
		if err := macro.SwitchTools(ctx, s.character, s.game, spot); err != nil {
			return fmt.Errorf("switch tools: %w", err)
		}
		s.current = spot.Name
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
