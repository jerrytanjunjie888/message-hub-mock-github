package service

import (
	"message-hub-mock/mock"

	"github.com/gin-gonic/gin"
)

type mockService struct {
	repo mock.Repository
}

// NewService constructs new service
func NewService(repo mock.Repository) mock.Service {
	return &mockService{
		repo: repo,
	}
}

func (service *mockService) SendStatus(c *gin.Context) {
	service.repo.SendStatus(c)
}

func (service *mockService) ReceiveMessage(c *gin.Context) {
	service.repo.ReceiveMessage(c)
}

func (service *mockService) DirService(c *gin.Context) {
	service.repo.DirRepository(c)
}

func (service *mockService) DirInfoService(c *gin.Context) {
	service.repo.DirInfoRepository(c)
}
