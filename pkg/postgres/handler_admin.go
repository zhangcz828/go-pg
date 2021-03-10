package postgres

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"go-pg/modules"
	"go-pg/pkg/connection"
	"log"
	"net/http"
)

// response format
type response struct {
	Name      string  `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateHero(w http.ResponseWriter, r *http.Request) {
	// create an empty hero of type models.Hero
	var hero modules.Hero

	// decode the json request to hero
	err := json.NewDecoder(r.Body).Decode(&hero)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertName := insertHero(hero)

	// format a response object
	res := response{
		Name:      insertName,
		Message: "Hero created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func insertHero(hero modules.Hero) string {

	// create the postgres db connection
	db := connection.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO hero VALUES ($1, $2, $3, $4, $5) RETURNING name`

	// the inserted id will store in this id
	var  name string

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement,
		hero.Name,
		hero.Detail,
		hero.AttackPower,
		hero.DefensePower,
		hero.Blood).Scan(&name)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	log.Printf("Inserted a hero: %v", name)

	// return the inserted id
	return name
}