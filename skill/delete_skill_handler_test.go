package skill

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"skillsapi/database"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSkill(t *testing.T) {
	db := database.NewPostgres()
	defer db.Close()

	h := Handler{Db: db}
	r := gin.Default()
	r.POST("/api/v1/skills", h.CreateSkill)
	r.DELETE("/api/v1/skills/:key", h.DeleteSkill)

	newSkill := Skill{
		Key:         "testDeleteSkill",
		Name:        "Test",
		Description: "test",
		Logo:        "test",
		Tags:        []string{"test"},
	}
	jsonValue, _ := json.Marshal(newSkill)
	req, _ := http.NewRequest("POST", "/api/v1/skills", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	t.Run("delete existing skill", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/skills/%v", newSkill.Key), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, "Skill deleted", response["message"])
	})

	t.Run("delete non-existent skill", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/skills/nonexistentkey", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]string
		json.NewDecoder(w.Body).Decode(&response)
		assert.Equal(t, "error", response["status"])
		assert.Equal(t, "Skill not found", response["message"])
	})
}
