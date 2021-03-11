package postgres

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"go-pg/cache"
	"go-pg/modules"
	"go-pg/pkg/connection"
	"log"
	"net/http"
)

func GetAllHeros(w http.ResponseWriter, r *http.Request) {
	var heros []modules.Hero

	// 缓存数据库的key
	k := cache.HeroList

	c := cache.GetCache()

	if value, ok := c.Get(k); ok {
		// v, ok := value.(modules.Hero)
		//if !ok {
		//	log.Fatal("It's not ok for type Hero")
		//}

		log.Printf("user get %s of value from Cache\n", k)

		json.NewEncoder(w).Encode(value)

		return
	}

	heros, err := getAllHero()
	if err != nil {
		log.Fatalf("Unable to get data from database. %v", err)

	}

	// 写缓存
	c.Add(k, heros)

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
