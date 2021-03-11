package postgres

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go-pg/modules"
	"go-pg/pkg/connection"
	"log"
	"net/http"
	"time"
)

type sessions struct {
	sessionMap map[string]modules.Session
}

func (ss *sessions) Insert(id string, s modules.Session) {
	ss.sessionMap[id] = s
}

func (ss *sessions) Get(id string) (modules.Session, bool) {
	v, ok := ss.sessionMap[id]
	return v, ok
}

// ssMap 初始化，注意map初始化的坑
var ssMap *sessions = &sessions{
	sessionMap: make(map[string]modules.Session),
}

func GetSessionById(c *gin.Context) {
	userId := c.Params.ByName("id")

	// 1. 查找内存中是否存在此session，若有则返回
	if v, ok := ssMap.Get(userId); ok {
		log.Printf("Loaded session %s from memory", userId)
		c.JSON(http.StatusOK, v)
		return
	}

	// 2. 查询database中是否archive了此session
	session, err := loadSessionFromDb(userId)
	if err != nil {
		log.Fatalf("Unable to get data from database. %v", err)
	}

	// 3. 若session为空，则新建一个, 并写入内存ssMap
	if session == nil {
		newSession := modules.Session{
			ArchiveDate: time.Now(),
		}

		ssMap.Insert(userId, newSession)

		c.JSON(http.StatusOK, newSession)
		return
	}

	// 4. 返回db中查找到的session
	log.Printf("Loaded session %s from database", userId)
	c.JSON(http.StatusOK, session)
}

func loadSessionFromDb(uid string) (*modules.Session, error){
	db := connection.CreateConnection()
	defer db.Close()

	var s modules.Session

	// create the select sql query
	sqlStatement := fmt.Sprintf("SELECT * FROM Session WHERE uid = %s", uid)

	// execute the sql statement
	row, err := db.Query(sqlStatement)
	row.Scan()
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
			&s.UID,
			&s.HeroName,
			&s.HeroBlood,
			&s.BossBlood,
			&s.CurrentLevel,
			&s.Score,
			&s.ArchiveDate)

	if err != nil {
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return &s, err
}
