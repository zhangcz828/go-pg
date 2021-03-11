package modules

import "time"

type Session struct {
	UID int
	HeroName string
	HeroBlood int
	BossBlood int
	CurrentLevel int
	Score int
	ArchiveDate time.Time
}
