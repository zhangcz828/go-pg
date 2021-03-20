package postgres

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go-pg/pkg/connection"
	"log"
	"net/http"
)

func CreateHeroHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := h.CreateHero(c)
		c.JSON(r.Status, r.Message)
	}
}

func AdjustHero(c *gin.Context) {
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
		log.Fatalf("Unable to execute the query. %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Heros update attack and defense power successfully.",
	})
}

