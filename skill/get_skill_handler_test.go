package skill

import (
	"net/http"
	"net/http/httptest"
	"skillsapi/database"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllSkills(t *testing.T) {
	database.ResetDB()

	db := database.NewPostgres()
	defer db.Close()

	handler := Handler{Db: db}
	router := gin.Default()
	router.GET("/api/v1/skills", handler.GetSkills)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/skills", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	expected := `{
		"status": "success",
		"data": [
			{
				"key": "go",
				"name": "Go",
				"description": "Go is a statically typed, compiled programming language designed at Google.",
				"logo": "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
				"tags": ["programming language", "system"]
			},
			{
				"key": "nodejs",
				"name": "Node.js",
				"description": "Node.js is an open-source, cross-platform, JavaScript runtime environment that executes JavaScript code outside of a browser.",
				"logo": "https://upload.wikimedia.org/wikipedia/commons/d/d9/Node.js_logo.svg",
				"tags": ["runtime", "javascript"]
			}
		]
	}`
	assert.JSONEq(t, expected, recorder.Body.String())
}

func TestGetSkillById(t *testing.T) {
	database.ResetDB()

	db := database.NewPostgres()
	defer db.Close()

	handler := Handler{Db: db}
	router := gin.Default()
	router.GET("/api/v1/skills/:key", handler.GetSkill)

	t.Run("should return a skill successfully", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/skills/go", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		expected := `{
			"status": "success",
			"data": {
				"key": "go",
				"name": "Go",
				"description": "Go is a statically typed, compiled programming language designed at Google.",
				"logo": "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
				"tags": ["programming language", "system"]
			}
		}`
		assert.JSONEq(t, expected, recorder.Body.String())
	})

	t.Run("should return 404 when skill is not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/skills/unknown", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		expected := `{"status":"error","message":"Skill not found"}`
		assert.JSONEq(t, expected, recorder.Body.String())
	})
}
