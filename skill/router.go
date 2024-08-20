package skill

import (
	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine, h *Handler) {
	r.GET("/ping", GetPing)
	r.GET("/api/v1/skills", h.GetSkills)
	r.GET("/api/v1/skills/:key", h.GetSkill)
	r.POST("/api/v1/skills", h.CreateSkill)
	r.PUT("/api/v1/skills/:key", h.UpdateSkill)
	r.PATCH("/api/v1/skills/:key/actions/name", h.UpdateSkillName)
	r.PATCH("/api/v1/skills/:key/actions/description", h.UpdateSkillDescription)
	r.PATCH("/api/v1/skills/:key/actions/logo", h.UpdateSkillLogo)
	r.PATCH("/api/v1/skills/:key/actions/tags", h.UpdateSkillTags)
	r.DELETE("/api/v1/skills/:key", h.DeleteSkill)
}
