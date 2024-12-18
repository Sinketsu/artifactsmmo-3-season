package strategy

import (
	"context"
	"fmt"
	"slices"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
)

type tasksMonsters struct {
	character *generic.Character
	game      *game.Game

	cancel []string
	food   []string
	events []string

	current string
}

func TasksMonsters(character *generic.Character, game *game.Game) *tasksMonsters {
	return &tasksMonsters{
		character: character,
		game:      game,
	}
}

func (s *tasksMonsters) Name() string {
	return "do monster tasks"
}

func (s *tasksMonsters) Cancel(tasks ...string) *tasksMonsters {
	s.cancel = append(s.cancel, tasks...)
	return s
}

func (s *tasksMonsters) UseFood(food ...string) *tasksMonsters {
	s.food = append(s.food, food...)
	return s
}

func (s *tasksMonsters) AllowEvents(events ...string) *tasksMonsters {
	s.events = append(s.events, events...)
	return s
}

func (s *tasksMonsters) Do(ctx context.Context) error {
	task := s.character.Task()

	if task == generic.NoTask {
		macro.AcceptMonsterTask(ctx, s.character, s.game)

		task = s.character.Task()
		if task == generic.NoTask {
			// if some error occured
			return fmt.Errorf("no task...")
		}
	}

	if slices.Contains(s.cancel, task.Code) {
		macro.CancelMonsterTask(ctx, s.character, s.game)
		return nil
	}

	if task.Current == task.Total {
		macro.CompleteMonsterTask(ctx, s.character, s.game)

		macro.Deposit(ctx, s.character, s.game, "tasks_coin")
		macro.DepositGold(ctx, s.character, s.game)

		if s.character.Inventory()["tasks_coin"] > 5 {
			if err := s.character.Deposit(ctx, "tasks_coin", s.character.Inventory()["tasks_coin"]-5); err != nil {
				return err
			}
		}

		return nil
	}

	if err := macro.Heal(ctx, s.character, s.game, s.food...); err != nil {
		return fmt.Errorf("heal: %w", err)
	}

	if s.character.InventoryFull() {
		macro.Deposit(ctx, s.character, s.game, "tasks_coin")
		macro.DepositGold(ctx, s.character, s.game)

		if s.character.Inventory()["tasks_coin"] > 5 {
			if err := s.character.Deposit(ctx, "tasks_coin", s.character.Inventory()["tasks_coin"]-5); err != nil {
				return err
			}
		}
	}

	var monster game.Point
	for _, event := range s.events {
		if event, err := s.game.GetEvent(event); err == nil {
			monster = event
			break
		}
	}

	if monster.Name == "" {
		var err error
		monster, err = s.game.Find(task.Code, s.character.Location())
		if err != nil {
			return fmt.Errorf("get map: %w", err)
		}
	}

	if s.current != monster.Name {
		if err := macro.SwitchGear(ctx, s.character, s.game, monster); err != nil {
			return fmt.Errorf("switch gear: %w", err)
		}
		s.current = monster.Name
	}

	err := s.character.Move(ctx, monster)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = s.character.Fight(ctx)
	if err != nil {
		return fmt.Errorf("fight: %w", err)
	}

	return nil
}
