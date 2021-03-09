package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"go-pg/modules"
	"log"
	"net/http"
)

// create connection with postgres db
func createConnection() *sql.DB {

	config := getDbConfig()
	// connection string
	pgConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Addr,
		config.Port,
		config.Username,
		config.Password,
		config.DBName)

	// open database
	db, err := sql.Open("postgres", pgConn)
	CheckError(err)

	// check the connection
	err = db.Ping()
	CheckError(err)

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
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
	db := createConnection()
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
