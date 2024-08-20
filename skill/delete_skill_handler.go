package skill

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteSkill(c *gin.Context) {
	key := c.Param("key")

	var exists bool
	err := h.Db.QueryRow(`SELECT EXISTS(SELECT 1 FROM skills WHERE key = $1)`, key).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Skill not found",
		})
		return
	}

	_, err = h.Db.Exec(`DELETE FROM skills WHERE key = $1`, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "not be able to delete skill",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Skill deleted",
	})
}
