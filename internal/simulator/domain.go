package simulator

type Result struct {
	Win                  bool
	Turns                int
	Seconds              int
	RemainingMonsterHp   int
	RemainingCharacterHp int
	NeedHeal             int
}
