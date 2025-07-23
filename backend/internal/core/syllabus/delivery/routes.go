package delivery

import (
	"study_buddy/internal/middlewares"
	"study_buddy/pkg/constants"

	"github.com/gin-gonic/gin"
)

func MapSyllabusRoutes(router *gin.Engine, h Handlers, authMiddleware gin.HandlerFunc) {
	router.GET("/get_syllabus", authMiddleware, middlewares.PermissionHandler(constants.SyllabusEntity, constants.Get), h.GetSyllabus)
	router.POST("/save_syllabus", authMiddleware, middlewares.PermissionHandler(constants.SyllabusEntity, constants.Create), h.SaveSyllabus)
	router.DELETE("/delete_syllabus", authMiddleware, middlewares.PermissionHandler(constants.SyllabusEntity, constants.Delete), h.DeleteSyllabus)
}
