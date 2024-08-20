package skill

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) GetSkills(c *gin.Context) {
	rows, err := h.Db.Query(`SELECT key, name, description, logo, tags FROM skills`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
		})
		return
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var skill Skill
		var tags pq.StringArray
		err := rows.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, &tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Internal server error",
			})
			return
		}
		skill.Tags = []string(tags)
		skills = append(skills, skill)
	}

	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skills,
	})
}

func (h *Handler) GetSkill(c *gin.Context) {
	key := c.Param("key")

	row := h.Db.QueryRow(`SELECT key, name, description, logo, tags FROM skills WHERE key = $1`, key)

	var skill Skill
	var tags pq.StringArray
	err := row.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, &tags)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Skill not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
		})
		return
	}
	skill.Tags = []string(tags)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}
