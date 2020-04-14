package mock

import "github.com/gin-gonic/gin"

// Service interface
type Service interface {
	SendStatus(c *gin.Context)
	ReceiveMessage(c *gin.Context)
	DirService(c *gin.Context)
	DirInfoService(c *gin.Context)
}
