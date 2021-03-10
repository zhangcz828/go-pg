package postgres

import (
	"encoding/json"
	"github.com/golang/groupcache"
	_ "github.com/lib/pq"
	"go-pg/cache"
	"go-pg/modules"
	"go-pg/pkg/connection"
	"log"
	"net/http"
)

func GetAllHeros(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var heros []modules.Hero

	// 缓存数据库的key
	k := r.URL.Path + r.Method

	log.Printf("user get %s of value from groupcache\n", k)
	cache.CreateHerosCacheGroup().Get(nil, k, groupcache.AllocatingByteSliceSink(&data))

	json.Unmarshal(data, &heros)
	json.NewEncoder(w).Encode(heros)
}

// GetAllHero will return all the Hero
func GetAllHero(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	allHero, err := getAllHero()
	if err != nil {
		log.Fatalf("Unable to get all heros. %v", err)
	}

	// send all the hero as response
	json.NewEncoder(w).Encode(allHero)
}

func getAllHero() ([]modules.Hero, error){
	db := connection.CreateConnection()
	defer db.Close()

	var heros []modules.Hero

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
		var hero modules.Hero

		// unmarshal the row object to hero
		err = rows.Scan(&hero.Name, &hero.Detail, &hero.AttackPower, &hero.DefensePower, &hero.Blood)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the hero to the heros slice
		heros = append(heros, hero)
	}

	return heros, err
}
