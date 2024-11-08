package generic

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"unsafe"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
)

func (c *Character) Move(ctx context.Context, to game.Point) error {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionMoveMyNameActionMovePost(ctx, &oas.DestinationSchema{
		X: to.X,
		Y: to.Y,
	}, oas.ActionMoveMyNameActionMovePostParams{
		Name: c.name,
	})
	if err != nil {
		return err
	}

	switch v := resp.(type) {
	case *oas.CharacterMovementResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		if to.Name != "" {
			c.logger.Debug(fmt.Sprintf("go to: %s", to.Name))
		} else {
			c.logger.Debug(fmt.Sprintf("go to: (%d, %d)", to.X, to.Y))
		}

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return nil
	case *oas.ActionMoveMyNameActionMovePostNotFound:
		return fmt.Errorf("target position not found")
	case *oas.ActionMoveMyNameActionMovePostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionMoveMyNameActionMovePostCode490:
		// character already at destination
		return nil
	case *oas.ActionMoveMyNameActionMovePostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionMoveMyNameActionMovePostCode499:
		return fmt.Errorf("cooldown...")
	}

	return fmt.Errorf("unknown error: %v", resp)
}

func (c *Character) Rest(ctx context.Context) error {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionRestMyNameActionRestPost(ctx, oas.ActionRestMyNameActionRestPostParams{
		Name: c.name,
	})
	if err != nil {
		return err
	}

	switch v := resp.(type) {
	case *oas.CharacterRestResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		c.logger.Debug(fmt.Sprintf("rest: %d hp restored", v.Data.HpRestored))

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return nil
	case *oas.ActionRestMyNameActionRestPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionRestMyNameActionRestPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionRestMyNameActionRestPostCode499:
		return fmt.Errorf("cooldown...")
	}

	return fmt.Errorf("unknown error: %v", resp)
}

func (c *Character) Fight(ctx context.Context) (*oas.CharacterFightDataSchemaFight, error) {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionFightMyNameActionFightPost(ctx, oas.ActionFightMyNameActionFightPostParams{
		Name: c.name,
	})
	if err != nil {
		return nil, err
	}

	switch v := resp.(type) {
	case *oas.CharacterFightResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		c.logger.Debug("fight", slog.Int("xp", v.Data.Fight.Xp), slog.Int("gold", v.Data.Fight.Gold), slog.Any("items", v.Data.Fight.Drops))

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return &v.Data.Fight, nil
	case *oas.ActionFightMyNameActionFightPostCode486:
		return nil, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionFightMyNameActionFightPostCode497:
		return nil, fmt.Errorf("inventory is full")
	case *oas.ActionFightMyNameActionFightPostCode498:
		return nil, fmt.Errorf("character not found")
	case *oas.ActionFightMyNameActionFightPostCode499:
		return nil, fmt.Errorf("cooldown...")
	case *oas.ActionFightMyNameActionFightPostCode598:
		return nil, fmt.Errorf("monster not found on this map")
	}

	return nil, fmt.Errorf("unknown error: %v", resp)
}

func (c *Character) Deposit(ctx context.Context, item string, quantity int) error {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionDepositBankMyNameActionBankDepositPost(ctx, &oas.SimpleItemSchema{
		Code:     item,
		Quantity: quantity,
	}, oas.ActionDepositBankMyNameActionBankDepositPostParams{
		Name: c.name,
	})
	if err != nil {
		return err
	}

	switch v := resp.(type) {
	case *oas.BankItemTransactionResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		c.logger.Debug(fmt.Sprintf("deposit: %d %s", quantity, item))

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return nil
	case *oas.ActionDepositBankMyNameActionBankDepositPostNotFound:
		return fmt.Errorf("item not found")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode461:
		return fmt.Errorf("transaction is already in progress with this item/your gold in your bank")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode462:
		return fmt.Errorf("bank is full")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode499:
		return fmt.Errorf("cooldown...")
	case *oas.ActionDepositBankMyNameActionBankDepositPostCode598:
		return fmt.Errorf("bank not found on this map")
	}

	return fmt.Errorf("unknown error: %v", resp)
}

func (c *Character) Withdraw(ctx context.Context, item string, quantity int) error {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionWithdrawBankMyNameActionBankWithdrawPost(ctx, &oas.SimpleItemSchema{
		Code:     item,
		Quantity: quantity,
	}, oas.ActionWithdrawBankMyNameActionBankWithdrawPostParams{
		Name: c.name,
	})
	if err != nil {
		return err
	}

	switch v := resp.(type) {
	case *oas.BankItemTransactionResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		c.logger.Debug(fmt.Sprintf("withdraw: %d %s", quantity, item))

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return nil
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostNotFound:
		return fmt.Errorf("item not found")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode461:
		return fmt.Errorf("transaction is already in progress with this item/your gold in your bank")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode478:
		return fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode497:
		return fmt.Errorf("inventory is full")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode499:
		return fmt.Errorf("cooldown...")
	case *oas.ActionWithdrawBankMyNameActionBankWithdrawPostCode598:
		return fmt.Errorf("bank not found on this map")
	}

	return fmt.Errorf("unknown error: %v", resp)
}

func (c *Character) DepositGold(ctx context.Context, quantity int) error {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionDepositBankGoldMyNameActionBankDepositGoldPost(ctx, &oas.DepositWithdrawGoldSchema{
		Quantity: quantity,
	}, oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostParams{
		Name: c.name,
	})
	if err != nil {
		return err
	}

	switch v := resp.(type) {
	case *oas.BankGoldTransactionResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		c.logger.Debug(fmt.Sprintf("deposit gold: %d", quantity))

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return nil
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode461:
		return fmt.Errorf("transaction is already in progress with this item/your gold in your bank")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode486:
		return fmt.Errorf("action is already in progress by your character")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode492:
		return fmt.Errorf("insufficient gold")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode498:
		return fmt.Errorf("character not found")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode499:
		return fmt.Errorf("cooldown...")
	case *oas.ActionDepositBankGoldMyNameActionBankDepositGoldPostCode598:
		return fmt.Errorf("bank not found on this map")
	}

	return fmt.Errorf("unknown error: %v", resp)
}

func (c *Character) Gather(ctx context.Context) (*oas.SkillDataSchemaDetails, error) {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionGatheringMyNameActionGatheringPost(ctx, oas.ActionGatheringMyNameActionGatheringPostParams{
		Name: c.name,
	})
	if err != nil {
		return nil, err
	}

	switch v := resp.(type) {
	case *oas.SkillResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		c.logger.Debug("gather", slog.Int("xp", v.Data.Details.Xp), slog.Any("items", v.Data.Details.Items))

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return &v.Data.Details, nil
	case *oas.ActionGatheringMyNameActionGatheringPostCode486:
		return nil, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionGatheringMyNameActionGatheringPostCode493:
		return nil, fmt.Errorf("skill level is too low")
	case *oas.ActionGatheringMyNameActionGatheringPostCode497:
		return nil, fmt.Errorf("inventory is full")
	case *oas.ActionGatheringMyNameActionGatheringPostCode498:
		return nil, fmt.Errorf("character not found")
	case *oas.ActionGatheringMyNameActionGatheringPostCode499:
		return nil, fmt.Errorf("cooldown...")
	case *oas.ActionGatheringMyNameActionGatheringPostCode598:
		return nil, fmt.Errorf("resource not found on this map")
	}

	return nil, fmt.Errorf("unknown error: %v", resp)
}

func (c *Character) Craft(ctx context.Context, code string, quantity int) (*oas.SkillDataSchemaDetails, error) {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionCraftingMyNameActionCraftingPost(ctx, &oas.CraftingSchema{
		Code:     code,
		Quantity: oas.NewOptInt(quantity),
	}, oas.ActionCraftingMyNameActionCraftingPostParams{
		Name: c.name,
	})
	if err != nil {
		return nil, err
	}

	switch v := resp.(type) {
	case *oas.SkillResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		c.logger.Debug(fmt.Sprintf("craft %d %s", quantity, code), slog.Int("xp", v.Data.Details.Xp), slog.Any("items", v.Data.Details.Items))

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return &v.Data.Details, nil
	case *oas.ActionCraftingMyNameActionCraftingPostNotFound:
		return nil, fmt.Errorf("craft not found")
	case *oas.ActionCraftingMyNameActionCraftingPostCode478:
		return nil, fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionCraftingMyNameActionCraftingPostCode486:
		return nil, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionCraftingMyNameActionCraftingPostCode493:
		return nil, fmt.Errorf("skill level is too low")
	case *oas.ActionCraftingMyNameActionCraftingPostCode497:
		return nil, fmt.Errorf("inventory is full")
	case *oas.ActionCraftingMyNameActionCraftingPostCode498:
		return nil, fmt.Errorf("character not found")
	case *oas.ActionCraftingMyNameActionCraftingPostCode499:
		return nil, fmt.Errorf("cooldown...")
	case *oas.ActionCraftingMyNameActionCraftingPostCode598:
		return nil, fmt.Errorf("workshop not found on this map")
	}

	return nil, fmt.Errorf("unknown error: %v", resp)
}

func (c *Character) Recycle(ctx context.Context, code string, quantity int) (*oas.RecyclingDataSchemaDetails, error) {
	apiRequestCount.Inc()

	resp, err := c.cli.ActionRecyclingMyNameActionRecyclingPost(ctx, &oas.RecyclingSchema{
		Code:     code,
		Quantity: oas.NewOptInt(quantity),
	}, oas.ActionRecyclingMyNameActionRecyclingPostParams{
		Name: c.name,
	})
	if err != nil {
		return nil, err
	}

	switch v := resp.(type) {
	case *oas.RecyclingResponseSchema:
		c.syncState(unsafe.Pointer(&v.Data.Character))
		c.logger.Debug(fmt.Sprintf("recycle %d %s", quantity, code), slog.Any("items", v.Data.Details.Items))

		time.Sleep(time.Duration(v.Data.Cooldown.RemainingSeconds) * time.Second)
		return &v.Data.Details, nil
	case *oas.ActionRecyclingMyNameActionRecyclingPostNotFound:
		return nil, fmt.Errorf("item not found")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode473:
		return nil, fmt.Errorf("item cannot be recycled")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode478:
		return nil, fmt.Errorf("missing item or insufficient quantity")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode486:
		return nil, fmt.Errorf("action is already in progress by your character")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode493:
		return nil, fmt.Errorf("skill level is too low")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode497:
		return nil, fmt.Errorf("inventory is full")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode498:
		return nil, fmt.Errorf("character not found")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode499:
		return nil, fmt.Errorf("cooldown...")
	case *oas.ActionRecyclingMyNameActionRecyclingPostCode598:
		return nil, fmt.Errorf("workshop not found on this map")
	}

	return nil, fmt.Errorf("unknown error: %v", resp)
}
