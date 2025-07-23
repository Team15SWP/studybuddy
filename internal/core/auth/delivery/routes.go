package delivery

import "github.com/gin-gonic/gin"

func MapAuthRoutes(router *gin.Engine, h Handlers) {
	router.POST("/login", h.LogIn)
	router.POST("/signup", h.SignUp)
	router.GET("/confirm", h.Confirm)
}
