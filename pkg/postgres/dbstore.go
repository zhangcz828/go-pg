package postgres

import (
	"fmt"
	"go-pg/cache"
	"go-pg/modules"
	"go-pg/pkg/connection"
	"log"
	"net/http"
	"time"
)

type DbStoreInterface interface {
	GetAllHeros() modules.Heros
	CreateHero(h modules.Hero) response
	LoadSessionFromDb(sid string) (*modules.SessionView, error)
	LoadBossFromDB(level int) (modules.Boss, error)
	LoadHeroFromDB(heroName string) (modules.Hero, error)
	Archive(session modules.Session) error
	RemoveSessionFromDB(sid string) error
	AdjustHero() error
}

// response format
type response struct {
	Status      int  `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// App use to talk with PG and implement DbStoreInterface
type DbStore struct {}

func (s *DbStore) GetAllHeros() modules.Heros {
	// 缓存数据库的key
	k := cache.HeroList

	ch := cache.GetCache()

	if value, ok := ch.Get(k); ok {
		log.Printf("user get %s of value from Cache\n", k)

		return value.(modules.Heros)
	}

	heros := getAllHero()

	// 写缓存
	ch.Add(k, heros)

	return heros
}

func (s *DbStore) CreateHero(hero modules.Hero) response {

	// create the postgres db connection
	db := connection.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO hero VALUES ($1, $2, $3, $4, $5);`

	// execute the sql statement
	// Scan function will save the insert id in the id
	_, err := db.Exec(sqlStatement,
		hero.Name,
		hero.Detail,
		hero.AttackPower,
		hero.DefensePower,
		hero.Blood)

	if err != nil {
		return response{
			Status: http.StatusInternalServerError,
			Message: "Unable to execute the query",
		}
	}

	// Insert successfully
	log.Printf("Inserted a hero: %v", hero.Name)

	return response{
		Status: http.StatusOK,
		Message: "Hero created successfully",
	}
}

func getAllHero() modules.Heros {
	db := connection.CreateConnection()
	defer db.Close()

	heros := modules.Heros{}

	// create the select sql query
	sqlStatement := `SELECT * FROM Hero`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	//iterate over the rows
	for rows.Next() {
		hero := modules.Hero{}

		// unmarshal the row object to hero
		err = rows.Scan(&hero.Name, &hero.Detail, &hero.AttackPower, &hero.DefensePower, &hero.Blood)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the hero to the heros slice
		heros = append(heros, &hero)
	}

	return heros
}

func (s *DbStore) LoadSessionFromDb(sid string) (*modules.SessionView, error) {
	db := connection.CreateConnection()
	defer db.Close()

	var sv modules.SessionView

	// create the select sql query
	sqlStatement := fmt.Sprintf("SELECT * FROM session_view WHERE sessionid = %s", sid)

	// execute the sql statement
	row, err := db.Query(sqlStatement)
	defer row.Close()
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err
	}

	// unmarshal the row object to hero
	if !row.Next() {
		// db 中没有数据，返回空
		return nil, nil
	}

	err = row.Scan(
		&sv.Session.UID,
		&sv.Session.HeroName,
		&sv.Hero.Detail,
		&sv.Hero.AttackPower,
		&sv.Hero.DefensePower,
		&sv.Hero.Blood,
		&sv.Session.LiveHeroBlood,
		&sv.Session.LiveBossBlood,
		&sv.Session.CurrentLevel,
		&sv.Session.Score,
		&sv.Session.ArchiveDate,
		&sv.Boss.Name,
		&sv.Boss.Detail,
		&sv.Boss.AttackPower,
		&sv.Boss.DefensePower,
		&sv.Boss.Blood)

	// Fill in other duplicated fields in SessionView
	//s.Hero.Name = s.Session.HeroName
	//s.Boss.Level = s.Session.CurrentLevel

	if err != nil {
		log.Printf("Unable to scan the row. %v", err)
		return nil, err
	}

	return &sv, err
}

func (s *DbStore) LoadBossFromDB(level int) (modules.Boss, error) {
	db := connection.CreateConnection()
	defer db.Close()

	b := modules.Boss{}

	// create the select sql query
	sqlStatement := fmt.Sprintf("SELECT * FROM boss WHERE level = %d", level)

	// execute the sql statement
	rows, err := db.Query(sqlStatement)
	defer rows.Close()
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return b, err
	}

	//iterate over the rows
	rows.Next()

	err = rows.Scan(&b.Name, &b.Detail, &b.AttackPower, &b.DefensePower, &b.Blood, &b.Level)

	if err != nil {
		log.Printf("Unable to scan the row. %v", err)
		return b, err
	}

	return b, nil
}

func (s *DbStore) LoadHeroFromDB(heroName string) (modules.Hero, error) {
	db := connection.CreateConnection()
	defer db.Close()

	h := modules.Hero{}

	// create the select sql query
	sqlStatement := fmt.Sprintf("SELECT * FROM hero WHERE name = '%s'", heroName)
	//sqlStatement := fmt.Sprint("SELECT * FROM hero WHERE attackpower=10")


	// execute the sql statement
	row, err := db.Query(sqlStatement)
	defer row.Close()
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return h, err
	}

	//iterate over the rows
	row.Next()

	err = row.Scan(&h.Name, &h.Detail, &h.AttackPower, &h.DefensePower, &h.Blood)

	if err != nil {
		log.Printf("Unable to scan the row. %v", err)
		return h, err
	}

	return h, nil
}

func (s *DbStore) Archive(session modules.Session) error {
	db := connection.CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO session(uid, heroname, heroblood, bossblood, currentlevel, score, archivedate) VALUES($1, $2, $3, $4, $5, $6, $7) ON conflict (uid) DO UPDATE SET heroblood = $8, bossblood = $9, currentlevel = $10, score = $11, archivedate = $12;`

	// execute the sql statement
	_, err := db.Exec(sqlStatement,
		session.UID,
		session.HeroName,
		session.LiveHeroBlood,
		session.LiveBossBlood,
		session.CurrentLevel,
		session.Score,
		time.Now(),
		session.LiveHeroBlood,
		session.LiveBossBlood,
		session.CurrentLevel,
		session.Score,
		time.Now())

	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return err
	}

	return nil
}

func (s *DbStore) RemoveSessionFromDB(sid string) error {
	db := connection.CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM session where uid = $1`

	// execute the sql statement
	_, err := db.Exec(sqlStatement, sid)

	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return err
	}

	return nil
}

func (s *DbStore) AdjustHero() error {
	// create the postgres db connection
	db := connection.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `UPDATE hero SET attackpower = attackpower*1.2, defensepower = defensepower*1.2`

	// execute the sql statement
	// Scan function will save the insert id in the id
	_, err := db.Exec(sqlStatement)

	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}

	return err
}
