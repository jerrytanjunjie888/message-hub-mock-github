package mock

import "github.com/gin-gonic/gin"

//Respository interface
type Repository interface {
	SendStatus(c *gin.Context)
	ReceiveMessage(c *gin.Context)
	DirRepository(c *gin.Context)
	DirInfoRepository(c *gin.Context)
}
