package main

import (
	"github.com/gin-gonic/gin"
	"go-pg/pkg/postgres"
)

func main() {
	gin.DisableConsoleColor()

	r := gin.Default()
	// User rest api
	r.GET("/heros", postgres.GinGetAllHeros)
	r.GET("/session/:id/", postgres.GetSessionById) //获取session
	r.PUT("/session/:id", postgres.SelectHero) //选择英雄,更新session, PUT /session/:id?hero="张无忌"
	r.PUT("/session/:id/fight", postgres.Fight)  //打boss, 更新session里面的信息，得分
	r.POST("/session/:id/archive", postgres.Archive) // ssMap[:id]存档到db
	r.POST("/session/:id/level", postgres.NextLevel) // 过关，更新session[加血，currentLevel, boss信息]
	//r.GET("/session/rank", postgres.Rank) //查看排名

	// Admin rest api
	r.POST("/admin/hero", postgres.CreateHero) //添加新的hero
	r.PUT("/admin/hero", postgres.AdjustHero) //节假日对hero进行调整，比如每个战斗力提升20%


	r.Run(":8000")
}