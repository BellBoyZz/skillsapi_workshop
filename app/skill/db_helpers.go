package skill

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func executeUpdate(c *gin.Context, db *sql.DB, query string, errorMessage string, args ...interface{}) {
	key := c.Param("key")

	_, err := db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errorMessage,
		})
		return
	}

	var skill Skill
	var tags pq.StringArray
	row := db.QueryRow(`SELECT key, name, description, logo, tags FROM skills WHERE key = $1`, key)
	err = row.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, &tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": errorMessage,
		})
		return
	}
	skill.Tags = []string(tags)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}
