package postgres

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go-pg/cache"
	"go-pg/modules"
	"net/http"
)

func CreateHeroHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// create an empty hero of type models.Hero
		var hero modules.Hero

		// decode the json request to hero
		err := json.NewDecoder(c.Request.Body).Decode(&hero)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid request body")
			return
		}

		// call insert user function and pass the user
		res := h.CreateHero(hero)

		if res.Status == http.StatusOK {
			// Delete the cache for listing all heros
			ch := cache.GetCache()
			ch.Remove(cache.HeroList)
		}

		c.String(res.Status, res.Message)
		return
	}
}

func AdjustHeroHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h.AdjustHero()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error in database %v", err)
			return
		}

		c.String(http.StatusOK, "Heros update attack and defense power successfully.")
	}
}

