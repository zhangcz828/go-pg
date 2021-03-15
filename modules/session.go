package modules

import "time"

type Session struct {
	UID string
	HeroName string
	LiveHeroBlood int
	LiveBossBlood int
	CurrentLevel int
	Score int
	ArchiveDate time.Time
}

type SessionView struct {
	Session
	Boss
	Hero
}
