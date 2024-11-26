package skill

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) CreateSkill(c *gin.Context) {
	var newSkill Skill

	if err := c.ShouldBindJSON(&newSkill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	var exists bool
	err := h.Db.QueryRow(`SELECT EXISTS(SELECT 1 FROM skills WHERE key=$1)`, newSkill.Key).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
		})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Skill already exists",
		})
		return
	}

	_, err = h.Db.Exec(`INSERT INTO skills (key, name, description, logo, tags) VALUES ($1, $2, $3, $4, $5)`,
		newSkill.Key, newSkill.Name, newSkill.Description, newSkill.Logo, pq.Array(newSkill.Tags))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newSkill,
	})
}
