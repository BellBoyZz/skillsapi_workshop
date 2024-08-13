package main

import (
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

var skills = []Skill{
	{
		Key: "go",
		Name: "Go",
		Description: "Go is a statically typed, compiled programming language designed at Google.",
		Logo: "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
		Tags: []string{"programming language", "system"},
	},
	{
		Key: "python",
		Name: "Python",
		Description: "Python is an interpreted, high-level, general-purpose programming language.",
		Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
		Tags: []string{"programming language", "scripting"},
	},
}

func GetSkills(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skills,
	})
}

func GetSkill(c *gin.Context) {
	key := c.Param("key")

	for _, skill := range skills {
		if skill.Key == key {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   skill,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "error",
		"message": "Skill not found",
	})
}

func CreateSkill(c *gin.Context) {
	var newSkill Skill

	if err := c.ShouldBindJSON(&newSkill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	for _, skill := range skills {
		if skill.Key == newSkill.Key {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Skill already exists",
			})
			return
		}
	}

	skills = append(skills, newSkill)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newSkill,
	})
}

func main() {
	r := gin.Default()

	r.GET("/api/v1/skills", GetSkills)
	r.GET("/api/v1/skills/:key", GetSkill)
	r.POST("/api/v1/skills", CreateSkill)

	r.Run(":8080")
}