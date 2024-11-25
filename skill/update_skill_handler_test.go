package skill

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"skillsapi/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSkill(t *testing.T) {
	db := database.NewPostgres()
	defer db.Close()
	h := Handler{Db: db}
	r := gin.Default()

	r.POST("/api/v1/skills", h.CreateSkill)
	r.PUT("/api/v1/skills/:key", h.UpdateSkill)

	t.Run("successfully update an existing skill", func(t *testing.T) {
		newSkill := Skill{
			Key:         "skillTestUpdate",
			Name:        "Initial Skill",
			Description: "Initial Description",
			Logo:        "initial_logo.png",
			Tags:        []string{"initial", "skill"},
		}
		createSkill(t, r, newSkill)

		updatedSkill := UpdateSkill{
			Name:        "Updated Skill",
			Description: "Updated Description",
			Logo:        "updated_logo.png",
			Tags:        []string{"updated", "skill"},
		}
		updateSkill(t, r, newSkill.Key, updatedSkill)

		assert.Equal(t, http.StatusOK, getResponseStatus(t, r, "PUT", "/api/v1/skills/skillTestUpdate", updatedSkill))
	})

	t.Run("fail to update due to invalid JSON payload", func(t *testing.T) {
		newSkill := Skill{
			Key:         "skillTestInvalidJSON",
			Name:        "Initial Skill",
			Description: "Initial Description",
			Logo:        "initial_logo.png",
			Tags:        []string{"initial", "skill"},
		}
		createSkill(t, r, newSkill)

		invalidPayload := UpdateSkill{
			Description: "Updated Description",
			Logo:        "updated_logo.png",
			Tags:        []string{"updated", "skill"},
		}
		assert.Equal(t, http.StatusBadRequest, getResponseStatus(t, r, "PUT", "/api/v1/skills/skillTestInvalidJSON", invalidPayload))
	})

	t.Run("fail to update due to mismatched data types", func(t *testing.T) {
		newSkill := Skill{
			Key:         "skillTestMismatchedType",
			Name:        "Initial Skill",
			Description: "Initial Description",
			Logo:        "initial_logo.png",
			Tags:        []string{"initial", "skill"},
		}
		createSkill(t, r, newSkill)

		type InvalidUpdateType struct {
			Name        int      `json:"name"`
			Description string   `json:"description"`
			Logo        string   `json:"logo"`
			Tags        []string `json:"tags"`
		}
		invalidUpdate := InvalidUpdateType{
			Name:        123,
			Description: "Updated Description",
			Logo:        "updated_logo.png",
			Tags:        []string{"updated", "skill"},
		}
		assert.Equal(t, http.StatusBadRequest, getResponseStatus(t, r, "PUT", "/api/v1/skills/skillTestMismatchedType", invalidUpdate))
	})

	t.Run("fail to update a nonexistent skill", func(t *testing.T) {
		nonexistentSkill := UpdateSkill{
			Name:        "Nonexistent Skill",
			Description: "This skill does not exist.",
			Logo:        "nonexistent_logo.png",
			Tags:        []string{"nonexistent"},
		}
		assert.Equal(t, http.StatusInternalServerError, getResponseStatus(t, r, "PUT", "/api/v1/skills/nonexistentKey", nonexistentSkill))
	})

	t.Run("fail to update due to missing required fields", func(t *testing.T) {
		missingFields := `{"name": "Updated Skill", "description": "Updated Description"}`
		req, _ := http.NewRequest("PUT", "/api/v1/skills/skillTestUpdate", bytes.NewBuffer([]byte(missingFields)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func createSkill(t *testing.T, r *gin.Engine, skill Skill) {
	t.Helper()
	jsonValue, _ := json.Marshal(skill)
	req, _ := http.NewRequest("POST", "/api/v1/skills", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func updateSkill(t *testing.T, r *gin.Engine, key string, update UpdateSkill) {
	t.Helper()
	jsonValue, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%s", key), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func getResponseStatus(t *testing.T, r *gin.Engine, method, url string, body interface{}) int {
	t.Helper()
	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}
