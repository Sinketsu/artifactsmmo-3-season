package simulator

import (
	"math"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
)

type Character interface {
	Level() int
}

type service struct{}

func New() *service {
	return &service{}
}

func (s *service) Fight(c Character, items []oas.ItemSchema, monster oas.MonsterSchema) Result {
	baseHp := 120 + (c.Level()-1)*5

	effects := map[string]float64{}
	for i := range items {
		for _, e := range items[i].Effects {
			effects[e.Name] += float64(e.Value)
		}
	}

	characterDamage := round(round(effects["attack_fire"]*(1.+effects["dmg_fire"]/100+effects["boost_dmg_fire"]/100))*(1.-float64(monster.ResFire)/100)) +
		round(round(effects["attack_air"]*(1.+effects["dmg_air"]/100+effects["boost_dmg_air"]/100))*(1.-float64(monster.ResAir)/100)) +
		round(round(effects["attack_earth"]*(1.+effects["dmg_earth"]/100+effects["boost_dmg_earth"]/100))*(1.-float64(monster.ResEarth)/100)) +
		round(round(effects["attack_water"]*(1.+effects["dmg_water"]/100+effects["boost_dmg_water"]/100))*(1.-float64(monster.ResWater)/100))

	monsterDamage := round(float64(monster.AttackFire)*(1.-effects["res_fire"]/100-effects["boost_res_fire"]/100)) +
		round(float64(monster.AttackAir)*(1.-effects["res_air"]/100-effects["boost_res_air"]/100)) +
		round(float64(monster.AttackEarth)*(1.-effects["res_earth"]/100-effects["boost_res_earth"]/100)) +
		round(float64(monster.AttackWater)*(1.-effects["res_water"]/100-effects["boost_res_water"]/100))

	totalHp := baseHp + int(effects["hp"]) + int(effects["boost_hp"])
	speed := effects["speed"]

	turnsToWin := math.Ceil(float64(monster.Hp) / characterDamage)
	turnsToLoose := math.Ceil(float64(totalHp) / monsterDamage)

	if turnsToWin > turnsToLoose {
		return Result{
			Win:                false,
			Turns:              int(turnsToLoose),
			Seconds:            int(round(2 * 2 * turnsToLoose * (1 - speed/100))),
			RemainingMonsterHp: monster.Hp - int(turnsToLoose*characterDamage),
		}
	}

	return Result{
		Win:                  true,
		Turns:                int(turnsToWin),
		Seconds:              int(round((2*2*turnsToWin - 1) * (1 - speed/100))),
		RemainingCharacterHp: totalHp - int((turnsToWin-1)*monsterDamage),
	}
}

func round(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) > 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}
