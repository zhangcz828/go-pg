package postgres

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go-pg/modules"
	"go-pg/pkg/connection"
	"log"
	"net/http"
	"sort"
	"time"
)

type sessions struct {
	sessionMap map[string]modules.SessionView
}

func (ss *sessions) Update(id string, s modules.SessionView) {
	ss.sessionMap[id] = s
}

func (ss *sessions) UpdateHero(id string, hero modules.Hero) {
	sv := ss.sessionMap[id]
	sv.Hero = hero
	sv.Session.HeroName = hero.Name
	sv.Session.LiveHeroBlood = hero.Blood
	ss.sessionMap[id] = sv
}

func (ss *sessions) Get(id string) (modules.SessionView, bool) {
	v, ok := ss.sessionMap[id]
	return v, ok
}

func (ss *sessions) GetSession(id string) modules.Session {
	v, _ := ss.sessionMap[id]
	return v.Session
}

func (ss *sessions) GetCurrentLevel(id string) int {
	return ssMap.sessionMap[id].CurrentLevel
}

func (ss *sessions) RemoveSession(id string) {
	delete(ss.sessionMap, id)
}

// ssMap 初始化，注意map初始化的坑
var ssMap *sessions = &sessions{
	sessionMap: make(map[string]modules.SessionView),
}

type ranking struct {
	Key string `json:"userid"`
	Value int `json:"score"`
}

// 实现排序的map, sort.Slice, link: https://duchengqian.com/go-sort.html
func (r *ranking) Get() []ranking {
	var rs []ranking
	for k, v := range ssMap.sessionMap {
		rs = append(rs, ranking{
			Key: k,
			Value: v.Score,
	})
	}

	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Value> rs[j].Value
	})

	//for _, kv := range ss {
	//	fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	//}

	return rs

}

func GetSessionById(c *gin.Context) {
	userId := c.Params.ByName("id")

	// 1. 查找内存中是否存在此session，若有则返回
	if v, ok := ssMap.Get(userId); ok {
		log.Printf("Loaded session %s from memory", userId)
		c.JSON(http.StatusOK, gin.H{
			"Session": v.Session,
			"Hero": v.Hero,
			"Boss": v.Boss,
		})
		return
	}

	// 2. 查询database中是否archive了此session
	s, err := loadSessionFromDb(userId)
	if err != nil {
		log.Fatalf("Unable to get data from database. %v", err)
	}

	// 3. 若session为空，则新建一个, 并写入内存ssMap
	if s == nil {

		// 设置level 为第一关，并load第一关的boss到session
		b := loadBossFromDB(1)
		newSession := modules.SessionView{
			Hero: modules.Hero{},
			Boss: b,
			Session: modules.Session{
				UID: userId,
				LiveBossBlood: b.Blood,
				CurrentLevel: b.Level,
				ArchiveDate: time.Now(),
			},
		}

		ssMap.Update(userId, newSession)

		c.JSON(http.StatusOK, gin.H{
			"Session": newSession.Session,
			"Hero": newSession.Hero,
			"Boss": newSession.Boss,
		})
		return
	}

	// 4. 返回db中查找到的session
	ssMap.Update(userId, *s) // 更新ssMap
	log.Printf("Loaded session %s from database", userId)
	c.JSON(http.StatusOK, gin.H{
		"Session": s.Session,
		"Hero": s.Hero,
		"Boss": s.Boss,
	})
}

func loadSessionFromDb(uid string) (*modules.SessionView, error){
	db := connection.CreateConnection()
	defer db.Close()

	var s modules.SessionView

	// create the select sql query
	sqlStatement := fmt.Sprintf("SELECT * FROM session_view WHERE sessionid = %s", uid)

	// execute the sql statement
	row, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer row.Close()

	// unmarshal the row object to hero
	if !row.Next() {
		// db 中没有数据，返回空
		return nil, nil
	}

	err = row.Scan(
			&s.Session.UID,
			&s.Session.HeroName,
			&s.Hero.Detail,
			&s.Hero.AttackPower,
			&s.Hero.DefensePower,
			&s.Hero.Blood,
			&s.Session.LiveHeroBlood,
			&s.Session.LiveBossBlood,
			&s.Session.CurrentLevel,
			&s.Session.Score,
			&s.Session.ArchiveDate,
			&s.Boss.Name,
			&s.Boss.Detail,
			&s.Boss.AttackPower,
			&s.Boss.DefensePower,
			&s.Boss.Blood)

	// Fill in other duplicated fields in SessionView
	//s.Hero.Name = s.Session.HeroName
	//s.Boss.Level = s.Session.CurrentLevel

	if err != nil {
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return &s, err
}

func loadBossFromDB(level int) modules.Boss {
	db := connection.CreateConnection()
	defer db.Close()

	var b modules.Boss

	// create the select sql query
	sqlStatement := fmt.Sprintf("SELECT * FROM boss WHERE level = %d", level)

	// execute the sql statement
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	//iterate over the rows
	rows.Next()

	err = rows.Scan(&b.Name, &b.Detail, &b.AttackPower, &b.DefensePower, &b.Blood, &b.Level)

	if err != nil {
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return b
}

func SelectHero(c *gin.Context) {
	sid := c.Params.ByName("id")

	heroName := c.Query("hero")
	hero := loadHeroFromDB(heroName)

	ssMap.UpdateHero(sid, hero)

	sv, _ := ssMap.Get(sid)

	c.JSON(http.StatusOK, gin.H{
		"Session": sv.Session,
		"Hero":    sv.Hero,
		"Boss":    sv.Boss,
	})
}

func loadHeroFromDB(heroName string) modules.Hero {
	db := connection.CreateConnection()
	defer db.Close()

	var h modules.Hero

	// create the select sql query
	sqlStatement := fmt.Sprintf("SELECT * FROM hero WHERE name = '%s'", heroName)
	//sqlStatement := fmt.Sprint("SELECT * FROM hero WHERE attackpower=10")


	// execute the sql statement
	row, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer row.Close()

	//iterate over the rows
	row.Next()

	err = row.Scan(&h.Name, &h.Detail, &h.AttackPower, &h.DefensePower, &h.Blood)

	if err != nil {
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return h
}

func Archive(c *gin.Context) {
	sid := c.Params.ByName("id")

	s := ssMap.GetSession(sid)

	archive(s)

	// format a response object
	res := struct{
		SessionID string
		Message string } {
		SessionID: sid,
		Message: "Session archived successfully",
	}

	// send the response
	c.JSON(http.StatusOK, res)
}

func archive(s modules.Session) {
	db := connection.CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO session(uid, heroname, heroblood, bossblood, currentlevel, score, archivedate) VALUES($1, $2, $3, $4, $5, $6, $7) ON conflict (uid) DO UPDATE SET heroblood = $8, bossblood = $9, currentlevel = $10, score = $11, archivedate = $12;`

	fmt.Println(sqlStatement)

	// execute the sql statement
	_, err := db.Exec(sqlStatement,
		s.UID,
		s.HeroName,
		s.LiveHeroBlood,
		s.LiveBossBlood,
		s.CurrentLevel,
		s.Score,
		time.Now(),
		s.LiveHeroBlood,
		s.LiveBossBlood,
		s.CurrentLevel,
		s.Score,
		time.Now())

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}

type fightResponse struct {
	HeroBlood int
	BossBlood int
	Score int
	GameOver bool `json:"gameover"`
	NextLevel bool `json:"nextlevel"`
}

func Fight(c *gin.Context) {
	sid := c.Params.ByName("id")

	sv, _ := ssMap.Get(sid)

	// 0. 战斗前的检查
	if sv.LiveHeroBlood <= 0 || sv.LiveBossBlood <= 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"Message": "GameOver or should be next level",
		})
	}

	// 1. 战斗开始, 模拟战斗过程，后期加入武器和大招以及判断得分
	sv.Session.LiveHeroBlood -= sv.Boss.AttackPower
	sv.Session.LiveBossBlood -= sv.Hero.AttackPower
	sv.Score += 10

	// 2. 判断hero 是否死亡？
	if sv.Session.LiveHeroBlood <= 0 {
		sv.Session.LiveHeroBlood = 0
		res := fightResponse{
			HeroBlood: sv.LiveHeroBlood,
			BossBlood: sv.LiveBossBlood,
			Score:     sv.Score,
			GameOver:  true,
			NextLevel: false,
		}
		ssMap.RemoveSession(sid)
		removeSessionFromDB(sid)
		c.JSON(http.StatusOK, res)
		return
	}

	// 3. 判断是否过关
	if sv.Session.LiveBossBlood <= 0 {
		sv.LiveBossBlood = 0
		res := fightResponse{
			HeroBlood: sv.LiveHeroBlood,
			BossBlood: sv.LiveBossBlood,
			Score: sv.Score,
			GameOver: false,
			NextLevel: true,
		}
		c.JSON(http.StatusOK, res)
		return
	}

	// 4. 正常返回
	res := fightResponse{
		HeroBlood: sv.LiveHeroBlood,
		BossBlood: sv.LiveBossBlood,
		Score: sv.Score,
		GameOver: false,
		NextLevel: false,
	}
	c.JSON(http.StatusOK, res)
}

func removeSessionFromDB(sid string) {
	db := connection.CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM session where uid = $1`

	// execute the sql statement
	_, err := db.Exec(sqlStatement, sid)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}

func NextLevel(c *gin.Context) {
	sid := c.Params.ByName("id")

	// load boss from next level
	b := loadBossFromDB(ssMap.GetCurrentLevel(sid)+1)

	sv, _ := ssMap.Get(sid)

	sv.Boss = b
	sv.Session.LiveHeroBlood = b.Blood
	sv.CurrentLevel += 1

	ssMap.Update(sid, sv)

	c.JSON(http.StatusOK, gin.H{
		"Session": sv.Session,
		"Message": "Go to the next Level!",
	})

}

func Quit(c *gin.Context) {
	sid := c.Params.ByName("id")

	s := ssMap.GetSession(sid)

	archive(s)

	ssMap.RemoveSession(sid)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Archived and quit successfully",
	})
}

func OnlineRanking(c *gin.Context) {
	var rk ranking
	rks := rk.Get()

	c.JSON(http.StatusOK, rks)
}