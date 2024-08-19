package main

import (
	"database/sql"
	"net/http"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Skill struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	Tags        []string `json:"tags"`
}

var db *sql.DB

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	r.GET("/api/v1/skills", getSkills)
	r.GET("/api/v1/skills/:key", getSkill)
	r.POST("/api/v1/skills", createSkill)
	r.PUT("/api/v1/skills/:key", updateSkill)
	r.PATCH("/api/v1/skills/:key/actions/name", updateSkillName)
	r.PATCH("/api/v1/skills/:key/actions/description", updateSkillDescription)
	r.PATCH("/api/v1/skills/:key/actions/logo", updateSkillLogo)
	r.PATCH("/api/v1/skills/:key/actions/tags", updateSkillTags)
	r.DELETE("/api/v1/skills/:key", deleteSkill)

	r.Run(":8080")
}

func initDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}
}

func getSkills(c *gin.Context) {
	rows, err := db.Query(`SELECT key, name, description, logo, tags FROM skills`)
	if err != nil {
		// log.Printf("Error executing query: %v", err)
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
			// log.Printf("Error scanning row: %v", err)
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
		// log.Printf("Error with rows: %v", rows.Err())
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

func getSkill(c *gin.Context) {
	key := c.Param("key")

	row := db.QueryRow(`SELECT key, name, description, logo, tags FROM skills WHERE key = $1`, key)

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
		// log.Printf("Error querying skill: %v", err)
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

func createSkill(c *gin.Context) {
	var newSkill Skill

	if err := c.ShouldBindJSON(&newSkill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM skills WHERE key=$1)`, newSkill.Key).Scan(&exists)
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

	_, err = db.Exec(`INSERT INTO skills (key, name, description, logo, tags) VALUES ($1, $2, $3, $4, $5)`,
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

// Helper function to execute update queries and return the updated skill
func executeUpdate(c *gin.Context, query, errorMessage string, args ...interface{}) {
	key := c.Param("key")

	_, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("Error executing update: %v", err)
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
		log.Printf("Error querying updated skill: %v", err)
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

func updateSkill(c *gin.Context) {
	var updatedSkill Skill
	if err := c.ShouldBindJSON(&updatedSkill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET name = $1, description = $2, logo = $3, tags = $4 WHERE key = $5`
	executeUpdate(c, query, "not be able to update skill", updatedSkill.Name, updatedSkill.Description, updatedSkill.Logo, pq.Array(updatedSkill.Tags), c.Param("key"))
}

func updateSkillName(c *gin.Context) {
	var updateName struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&updateName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET name = $1 WHERE key = $2`
	executeUpdate(c, query, "not be able to update skill name", updateName.Name, c.Param("key"))
}

func updateSkillDescription(c *gin.Context) {
	var updateDescription struct {
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&updateDescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET description = $1 WHERE key = $2`
	executeUpdate(c, query, "not be able to update skill description", updateDescription.Description, c.Param("key"))
}

func updateSkillLogo(c *gin.Context) {
	var updateLogo struct {
		Logo string `json:"logo"`
	}
	if err := c.ShouldBindJSON(&updateLogo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET logo = $1 WHERE key = $2`
	executeUpdate(c, query, "not be able to update skill logo", updateLogo.Logo, c.Param("key"))
}

func updateSkillTags(c *gin.Context) {
	var updateTags struct {
		Tags []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&updateTags); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	query := `UPDATE skills SET tags = $1 WHERE key = $2`
	executeUpdate(c, query, "not be able to update skill tags", pq.Array(updateTags.Tags), c.Param("key"))
}

func deleteSkill(c *gin.Context) {
	key := c.Param("key")

	_, err := db.Exec(`DELETE FROM skills WHERE key = $1`, key)
	if err != nil {
		log.Printf("Error deleting skill: %v", err)
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
