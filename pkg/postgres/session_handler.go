package postgres

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go-pg/modules"
	"log"
	"net/http"
	"sort"
	"time"
)

type fightResponse struct {
	HeroBlood int
	BossBlood int
	Score int
	GameOver bool `json:"gameover"`
	NextLevel bool `json:"nextlevel"`
}

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

func GetSessionHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		s, err := h.LoadSessionFromDb(userId)
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to get data from database.%v", err)
			return
		}

		// 3. 若session为空，则新建一个, 并写入内存ssMap
		if s == nil {

			// 设置level 为第一关，并load第一关的boss到session
			b, err := h.LoadBossFromDB(1)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error in database %v", err)
				return
			}
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
}

func SelectHeroHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Params.ByName("id")

		heroName := c.Query("hero")
		hero, err := h.LoadHeroFromDB(heroName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error in database %v", err)
			return
		}

		ssMap.UpdateHero(sid, hero)

		sv, _ := ssMap.Get(sid)

		c.JSON(http.StatusOK, gin.H{
			"Session": sv.Session,
			"Hero":    sv.Hero,
			"Boss":    sv.Boss,
		})
	}
}

func ArchiveSessionHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Params.ByName("id")

		s := ssMap.GetSession(sid)

		if s.HeroName == "" {
			c.String(http.StatusBadRequest, "Please select a hero before archive")
			return
		}

		err := h.Archive(s)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error in database %v", err)
			return
		}

		c.String(http.StatusOK, "Session archived successfully")
	}
}

func NextLevelHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Params.ByName("id")

		// load boss from next level
		b, err := h.LoadBossFromDB(ssMap.GetCurrentLevel(sid)+1)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error in database %v", err)
			return
		}

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


}

func QuitHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Params.ByName("id")

		s := ssMap.GetSession(sid)

		err := h.Archive(s)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error in database %v", err)
			return
		}

		ssMap.RemoveSession(sid)

		c.JSON(http.StatusOK, gin.H{
			"Message": "Archived and quit successfully",
		})
	}
}

func OnlineRankingHandler(c *gin.Context) {
	var rk ranking
	rks := rk.Get()

	c.JSON(http.StatusOK, rks)
}

func FightHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			err := h.RemoveSessionFromDB(sid)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error in database %v", err)
				return
			}

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
}