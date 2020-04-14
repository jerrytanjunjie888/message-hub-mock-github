package repository

import (
	"encoding/json"
	"message-hub-mock/mock"
	"message-hub-mock/mock/model"
	"net/http"
	"os"
	"sync"

	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"go.uber.org/zap"
)

type mockRepository struct {
	Logger                *zap.Logger
	CallbackServerAddress string
}

// NewRepository Constructs a new mock repo instance
func NewRepository(logger *zap.Logger) mock.Repository {
	callbackServerAddress := os.Getenv("MESSAGE_HUB_CALLBACK_SERVER_ADDRESS")

	return &mockRepository{
		Logger:                logger,
		CallbackServerAddress: callbackServerAddress,
	}
}

// SendStatus method calls the callback endpoint to update a message received
func (repo *mockRepository) SendStatus(c *gin.Context) {
	input := model.Request{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		repo.Logger.Error(err.Error())
	}

	callbackPath := "/message/status"
	callbackReqURL := repo.CallbackServerAddress + callbackPath

	t, _ := json.Marshal(input)
	repo.Logger.Debug("Core | Mock |" + string(t))

	req.Debug = true
	r, err := req.Post(callbackReqURL, req.BodyJSON(&input))

	resp := r.Response()
	repo.Logger.Info(resp.Status)

	result := model.Response{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		repo.Logger.Error(err.Error())
	}

	switch resp.StatusCode {
	case http.StatusOK:
		c.SecureJSON(http.StatusOK, "Success")

	default:
		c.SecureJSON(http.StatusBadRequest, "fail")
	}
}

func (repo *mockRepository) ReceiveMessage(c *gin.Context) {
	input := model.Request{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		repo.Logger.Error(err.Error())
	}

	// obj := map[string]interface{}{}
	// obj["sequenceID"] = input.SequenceID
	// obj["status"] = input.Status

	// json, err := json.Marshal(&obj)

	// client := config.LoadRedis()

	// redisErr := client.Set("NewMessage", json, 0).Err()
	// if redisErr != nil {
	// 	fmt.Println(redisErr)
	// }

	// redisVal, err := client.Get("NewMessage").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(redisVal)

	result := model.Response{}
	result.Code = 200
	result.Message = "成功"
	c.SecureJSON(http.StatusOK, result)
}

func (repo *mockRepository) DirRepository(c *gin.Context) {
	//golang获取文件夹下面的文件列表
	path := c.Param("path")

	baseRestURL := "./root/"
	apiSpecificPath := path
	reqURL := baseRestURL + apiSpecificPath
	files, _ := ioutil.ReadDir(reqURL)

	var DirParam model.DirParam
	var DirParamList []model.DirParam

	for _, f := range files {
		DirParam.Name = f.Name()
		DirParam.IsDir = f.IsDir()
		DirParam.Size = f.Size()
		DirParamList = append(DirParamList, DirParam)
	}

	result := model.DirResponse{}

	result.Path = reqURL
	result.Dirs = DirParamList

	c.SecureJSON(http.StatusOK, result)

}

func (repo *mockRepository) getDirRepository(path string) (res *model.DirResponse) {
	//golang获取文件夹下面的文件列表

	baseRestURL := "./root/"
	reqURL := baseRestURL + path
	files, _ := ioutil.ReadDir(reqURL)

	var DirParam model.DirParam
	var DirParamList []model.DirParam

	for _, f := range files {
		DirParam.Name = f.Name()
		DirParam.IsDir = f.IsDir()
		DirParam.Size = f.Size()
		DirParamList = append(DirParamList, DirParam)
	}

	result := &model.DirResponse{}

	result.Path = reqURL
	result.Dirs = DirParamList

	return result

}

func (repo *mockRepository) counterWaitGroup(path string) (res *model.DirResponse, err error) {
	var mut sync.Mutex
	var wg sync.WaitGroup
	result := &model.DirResponse{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock()
			result = repo.getDirRepository(path)
			wg.Done()
		}()
	}
	wg.Wait()
	return result, err
}

func (repo *mockRepository) DirInfoRepository(c *gin.Context) {
	//golang获取文件夹下面的文件列表
	path := c.Param("path")

	dirData, err := repo.counterWaitGroup(path)
	if err != nil {
		repo.Logger.Error(err.Error())
	}

	files, _ := ioutil.ReadDir(dirData.Path)

	var DirInfoResponse model.DirInfoResponse

	var fileCount, dirCount, totalSize int64
	for _, f := range files {
		if f.IsDir() {
			dirCount++
		}
		totalSize = totalSize + f.Size()
		fileCount++
	}

	DirInfoResponse.Path = dirData.Path
	DirInfoResponse.DirCount = dirCount
	DirInfoResponse.FileCount = fileCount
	DirInfoResponse.TotalSize = totalSize

	c.SecureJSON(http.StatusOK, DirInfoResponse)

}
