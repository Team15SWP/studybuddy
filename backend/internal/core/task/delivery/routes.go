package delivery

import (
	"study_buddy/internal/middlewares"
	"study_buddy/pkg/constants"

	"github.com/gin-gonic/gin"
)

func MapTaskRoutes(router *gin.Engine, h Handlers, authMiddleware gin.HandlerFunc) {
	router.POST("/generate_task", authMiddleware, middlewares.PermissionHandler(constants.TaskEntity, constants.Get), h.GenerateTask)
	router.POST("/evaluate_code", authMiddleware, middlewares.PermissionHandler(constants.TaskEntity, constants.Create), h.EvaluateCodeForTask)
	router.POST("/submit_code", authMiddleware, middlewares.PermissionHandler(constants.TaskEntity, constants.Create), h.EvaluateCodeForTask)
	router.GET("/get_stats", authMiddleware, middlewares.PermissionHandler(constants.TaskEntity, constants.Get), h.GetStatistics)
	router.GET("/notification-settings", authMiddleware, middlewares.PermissionHandler(constants.TaskEntity, constants.Get), h.GetNotificationSettings)
	router.POST("/notification-settings", authMiddleware, middlewares.PermissionHandler(constants.TaskEntity, constants.Get), h.SetNotificationSettings)
}
