package postgres

import (
	"go-pg/modules"
)

type DbStoreMock struct{}

func (s DbStoreMock) GetAllHeros() modules.Heros {
	heros := modules.Heros{
		{
			Name:"杨过",
			Detail:"武功：黯然销魂掌, 蛤蟆功 xxxxxxxxxxxxx",
			AttackPower:30,
			DefensePower:15,
			Blood:100,
		},
		{
			Name: "xx",
			Detail: "xx",
			AttackPower: 2,
			DefensePower: 4,
			Blood: 100,
		},
	}

	return heros
}

func (s DbStoreMock) CreateHero(h modules.Hero) response {
	return response{
	}
}

func (s DbStoreMock) LoadSessionFromDb(sid string) (*modules.SessionView, error) {
	return &modules.SessionView{
		Session: modules.Session{},
		Boss:    modules.Boss{},
		Hero:    modules.Hero{},
	}, nil
}

func (s DbStoreMock) LoadBossFromDB(level int) (modules.Boss, error) {
	return modules.Boss{}, nil
}

func (s DbStoreMock) LoadHeroFromDB(heroName string) (modules.Hero, error) {
	return modules.Hero{}, nil
}

func (s DbStoreMock) Archive(session modules.Session) error {
	return nil
}

func (s DbStoreMock) RemoveSessionFromDB(sid string) error {
	return nil
}

func (s DbStoreMock) AdjustHero() error {
	return nil
}
