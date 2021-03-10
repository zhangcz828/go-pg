package connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// create connection with postgres db
func CreateConnection() *sql.DB {

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

	log.Println("Successfully connected!")
	// return the connection
	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
