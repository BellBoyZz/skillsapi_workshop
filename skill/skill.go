package skill

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Skill struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	Tags        []string `json:"tags"`
}

type Handler struct {
	Db *sql.DB
}

func GetPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
