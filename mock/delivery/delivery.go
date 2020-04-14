package delivery

import (
	"message-hub-mock/mock/repository"
	"message-hub-mock/mock/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TestAlive to test if the application is running successfully
func TestAlive(c *gin.Context, logger *zap.Logger) {

	c.SecureJSON(http.StatusOK, "Mock Service is working")
}

// SendStatus function which will be delivered to main.go
func SendStatus(c *gin.Context, logger *zap.Logger) {
	service.NewService(repository.NewRepository(logger)).SendStatus(c)
}

// ReceiveMessage function will emulate response for pushing a message to lien
func ReceiveMessage(c *gin.Context, logger *zap.Logger) {
	service.NewService(repository.NewRepository(logger)).ReceiveMessage(c)
}

func DirService(c *gin.Context, logger *zap.Logger) {
	service.NewService(repository.NewRepository(logger)).DirService(c)
}

func DirInfoService(c *gin.Context, logger *zap.Logger) {
	service.NewService(repository.NewRepository(logger)).DirInfoService(c)
}
