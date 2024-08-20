package skill

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) UpdateSkillName(c *gin.Context) {
	var updateName struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET name = $1 WHERE key = $2`
	executeUpdate(c, h.Db, query, "not be able to update skill name", updateName.Name, c.Param("key"))
}

func (h *Handler) UpdateSkillDescription(c *gin.Context) {
	var updateDescription struct {
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateDescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET description = $1 WHERE key = $2`
	executeUpdate(c, h.Db, query, "not be able to update skill description", updateDescription.Description, c.Param("key"))
}

func (h *Handler) UpdateSkillLogo(c *gin.Context) {
	var updateLogo struct {
		Logo string `json:"logo" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateLogo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET logo = $1 WHERE key = $2`
	executeUpdate(c, h.Db, query, "not be able to update skill logo", updateLogo.Logo, c.Param("key"))
}

func (h *Handler) UpdateSkillTags(c *gin.Context) {
	var updateTags struct {
		Tags []string `json:"tags" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateTags); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET tags = $1 WHERE key = $2`
	executeUpdate(c, h.Db, query, "not be able to update skill tags", pq.Array(updateTags.Tags), c.Param("key"))
}
