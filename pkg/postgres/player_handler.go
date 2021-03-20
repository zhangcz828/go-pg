package postgres

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

func GetHerosHandler(h DbStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		heros := h.GetAllHeros()
		c.JSON(http.StatusOK, heros)
	}
}