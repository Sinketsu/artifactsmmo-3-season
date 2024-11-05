package strategy

import "github.com/Sinketsu/artifactsmmo-3-season/internal/generic"

type simpleFight struct {
	character *generic.Character

	monster string
	deposit []string
}

func SimpleFight(character *generic.Character) *simpleFight {
	return &simpleFight{
		character: character,
	}
}

func (s *simpleFight) With(monster string) *simpleFight {
	s.monster = monster
	return s
}

func (s *simpleFight) Deposit(items ...string) *simpleFight {
	s.deposit = items
	return s
}

func (s *simpleFight) Name() string {
	return "fight with " + s.monster
}

func (s *simpleFight) Do() error {
	return nil
}
