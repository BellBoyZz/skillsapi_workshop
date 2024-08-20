package skill

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type UpdateSkill struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Logo        string   `json:"logo" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
}

func (h *Handler) UpdateSkill(c *gin.Context) {
	var updatedSkill UpdateSkill
	if err := c.ShouldBindJSON(&updatedSkill); err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET name = $1, description = $2, logo = $3, tags = $4 WHERE key = $5`
	executeUpdate(c, h.Db, query, "not be able to update skill", updatedSkill.Name, updatedSkill.Description, updatedSkill.Logo, pq.Array(updatedSkill.Tags), c.Param("key"))
}
