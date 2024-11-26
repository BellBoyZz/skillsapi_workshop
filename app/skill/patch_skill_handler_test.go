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

func createTestSkill(r *gin.Engine, skill Skill) {
	jsonValue, _ := json.Marshal(skill)
	req, _ := http.NewRequest("POST", "/api/v1/skills", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		panic(fmt.Sprintf("Failed to create test skill: %v", w.Body.String()))
	}
}

func TestPatchSkillName(t *testing.T) {
	db := database.NewPostgres()
	defer db.Close()

	h := Handler{Db: db}
	r := gin.Default()
	r.POST("/api/v1/skills", h.CreateSkill)
	r.PUT("/api/v1/skills/:key/action/name", h.UpdateSkillName)

	skill := Skill{
		Key:         "testUpdateName",
		Name:        "Test",
		Description: "test",
		Logo:        "test",
		Tags:        []string{"test"},
	}

	createTestSkill(r, skill)

	t.Run("update skill name", func(t *testing.T) {
		name := struct {
			Name string `json:"name"`
		}{
			Name: "updateName",
		}
		jsonValue, _ := json.Marshal(name)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/name", skill.Key), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("update skill name to unexisted skill", func(t *testing.T) {
		name := struct {
			Name string `json:"name"`
		}{
			Name: "updateName",
		}
		jsonValue, _ := json.Marshal(name)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/name", "unexistedkey"), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("update skill name with empty skill name", func(t *testing.T) {
		name := struct {
			Name string `json:"name"`
		}{
			Name: "",
		}
		jsonValue, _ := json.Marshal(name)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/name", skill.Key), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPatchSkillDescription(t *testing.T) {
	db := database.NewPostgres()
	defer db.Close()

	h := Handler{Db: db}
	r := gin.Default()
	r.POST("/api/v1/skills", h.CreateSkill)
	r.PUT("/api/v1/skills/:key/action/description", h.UpdateSkillDescription)

	skill := Skill{
		Key:         "testUpdateDescription",
		Name:        "Test",
		Description: "test",
		Logo:        "test",
		Tags:        []string{"test"},
	}

	createTestSkill(r, skill)

	t.Run("update skill description", func(t *testing.T) {
		desc := struct {
			Desc string `json:"description"`
		}{
			Desc: "updateDescription",
		}
		jsonValue, _ := json.Marshal(desc)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/description", skill.Key), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("update skill description to unexisted skill", func(t *testing.T) {
		desc := struct {
			Desc string `json:"description"`
		}{
			Desc: "updateDescription",
		}
		jsonValue, _ := json.Marshal(desc)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/description", "unexistedkey"), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("update skill description with empty skill description", func(t *testing.T) {
		desc := struct {
			Desc string `json:"description"`
		}{
			Desc: "",
		}
		jsonValue, _ := json.Marshal(desc)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/description", skill.Key), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPatchSkillLogo(t *testing.T) {
	db := database.NewPostgres()
	defer db.Close()

	h := Handler{Db: db}
	r := gin.Default()
	r.POST("/api/v1/skills", h.CreateSkill)
	r.PUT("/api/v1/skills/:key/action/logo", h.UpdateSkillLogo)

	skill := Skill{
		Key:         "testUpdateLogo",
		Name:        "Test",
		Description: "test",
		Logo:        "test",
		Tags:        []string{"test"},
	}

	createTestSkill(r, skill)

	t.Run("update skill logo", func(t *testing.T) {
		logo := struct {
			Logo string `json:"logo"`
		}{
			Logo: "updateLogo",
		}
		jsonValue, _ := json.Marshal(logo)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/logo", skill.Key), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("update skill logo to unexisted skill", func(t *testing.T) {
		logo := struct {
			Logo string `json:"logo"`
		}{
			Logo: "updateLogo",
		}
		jsonValue, _ := json.Marshal(logo)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/logo", "unexistedkey"), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("update skill logo with empty skill logo", func(t *testing.T) {
		logo := struct {
			Logo string `json:"logo"`
		}{
			Logo: "",
		}
		jsonValue, _ := json.Marshal(logo)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/logo", skill.Key), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPatchSkillTags(t *testing.T) {
	db := database.NewPostgres()
	defer db.Close()

	h := Handler{Db: db}
	r := gin.Default()
	r.POST("/api/v1/skills", h.CreateSkill)
	r.PUT("/api/v1/skills/:key/action/tags", h.UpdateSkillTags)

	skill := Skill{
		Key:         "testUpdateTags",
		Name:        "Test",
		Description: "test",
		Logo:        "test",
		Tags:        []string{"test"},
	}

	createTestSkill(r, skill)

	t.Run("update skill tags", func(t *testing.T) {
		tags := struct {
			Tags []string `json:"tags"`
		}{
			Tags: []string{"update"},
		}
		jsonValue, _ := json.Marshal(tags)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/tags", skill.Key), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("update skill tags to unexisted skill", func(t *testing.T) {
		tags := struct {
			Tags []string `json:"tags"`
		}{
			Tags: []string{"update"},
		}
		jsonValue, _ := json.Marshal(tags)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/tags", "unexistedkey"), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("update skill tags with empty skill tags", func(t *testing.T) {
		tags := struct {
			Tags []string `json:"tags"`
		}{
			Tags: nil,
		}
		jsonValue, _ := json.Marshal(tags)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/skills/%v/action/tags", skill.Key), bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
