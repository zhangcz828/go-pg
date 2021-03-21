package main

import (
	"github.com/gin-gonic/gin"
	"go-pg/pkg/postgres"
)

func main() {
	gin.DisableConsoleColor()

	store := postgres.DbStore{}

	r := gin.Default()

	// User rest api
	r.GET("/heros", postgres.GetHerosHandler(&store))
	r.GET("/session/:id/", postgres.GetSessionHandler(&store)) //获取session
	r.PUT("/session/:id", postgres.SelectHeroHandler(&store)) //选择英雄,更新session, PUT /session/:id?hero="张无忌"
	r.PUT("/session/:id/fight", postgres.FightHandler(&store))  //打boss, 更新session里面的信息，得分
	r.POST("/session/:id/archive", postgres.ArchiveSessionHandler(&store)) // ssMap[:id]存档到db
	r.POST("/session/:id/level", postgres.NextLevelHandler(&store)) // 过关，更新session[加血，currentLevel, boss信息]
	r.POST("/session/:id/quit", postgres.QuitHandler(&store)) // 玩家下线，自动存档
	r.GET("/session/:id/rank", postgres.OnlineRankingHandler) //在线玩家积分排名

	// Admin rest api
	r.POST("/admin/hero", postgres.CreateHeroHandler(&store)) //添加新的hero
	r.PUT("/admin/hero", postgres.AdjustHeroHandler(&store)) //节假日对hero进行调整，比如每个战斗力提升20%


	r.Run(":8000")
}