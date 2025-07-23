package utils

import "github.com/gin-gonic/gin"

func SetAuthorizationToken(c *gin.Context, token string) {
	c.Header("Authorization", "Bearer "+token)
}

func SetError(err string) gin.H {
	res := gin.H{"error": err}
	return res
}

func SetMessage(msg string) gin.H {
	res := gin.H{"message": msg}
	return res
}
