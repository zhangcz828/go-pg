package postgres

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-pg/cache"
	"go-pg/modules"
	"go-pg/pkg/connection"
	"log"
	"net/http"
)

type DbStoreInterface interface {
	GetAllHeros() modules.Heros
	CreateHero(c *gin.Context) response
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

func (s *DbStore) CreateHero(c *gin.Context) response {
	// create an empty hero of type models.Hero
	var hero modules.Hero

	// decode the json request to hero
	err := json.NewDecoder(c.Request.Body).Decode(&hero)

	if err != nil {
		return response{
			Status: http.StatusBadRequest,
			Message: "Invalid request body",
		}
	}

	// call insert user function and pass the user
	res := createHero(hero)

	if res.Status == http.StatusOK {
		// Delete the cache for listing all heros
		ch := cache.GetCache()
		ch.Remove(cache.HeroList)
	}

	return res
}

func createHero(hero modules.Hero) response {

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

