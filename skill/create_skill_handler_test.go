package skill

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"skillsapi/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateSkill(t *testing.T) {
	database.ResetDB()

	db := database.NewPostgres()
	handler := Handler{Db: db}

	router := gin.Default()
	router.POST("/api/v1/skills", handler.CreateSkill)

	t.Run("should create a skill successfully", func(t *testing.T) {
		newSkill := Skill{
			Key:         "rust",
			Name:        "Rust",
			Description: "Rust is a multi-paradigm system programming language.",
			Logo:        "https://upload.wikimedia.org/wikipedia/commons/d/d5/Rust_programming_language_black_logo.svg",
			Tags:        []string{"rust", "systems programming"},
		}

		jsonData, err := json.Marshal(newSkill)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/skills", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"key":         "rust",
				"name":        "Rust",
				"description": "Rust is a multi-paradigm system programming language.",
				"logo":        "https://upload.wikimedia.org/wikipedia/commons/d/d5/Rust_programming_language_black_logo.svg",
				"tags":        []interface{}{"rust", "systems programming"},
			},
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return an error when the skill already exists", func(t *testing.T) {
		newSkill := Skill{
			Key:         "rust",
			Name:        "Rust",
			Description: "Rust is a multi-paradigm system programming language.",
			Logo:        "https://upload.wikimedia.org/wikipedia/commons/d/d5/Rust_programming_language_black_logo.svg",
			Tags:        []string{"rust", "systems programming"},
		}

		jsonData, err := json.Marshal(newSkill)
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/skills", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"status":  "error",
			"message": "Skill already exists",
		}
		assert.Equal(t, expected, response)
	})

	t.Run("should return an error for invalid request payload", func(t *testing.T) {
		invalidJSON := `{"key": "rust", "name": "Rust", "description": "Rust is a multi-paradigm system programming language."`

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/skills", bytes.NewBuffer([]byte(invalidJSON)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"status":  "error",
			"message": "Invalid request payload",
		}
		assert.Equal(t, expected, response)
	})
}
