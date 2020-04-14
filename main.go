package main

import (
	"message-hub-mock/mock/delivery"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	// Load redis
	// config.LoadRedis()

	// LumberJack Implementation
	var logger *zap.Logger

	logMaxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
	logMaxBackups, _ := strconv.Atoi(os.Getenv("LOG_MAX_BACKUPS"))
	logMaxAge, _ := strconv.Atoi(os.Getenv("LOG_MAX_AGE"))
	logCompress, _ := strconv.ParseBool(os.Getenv("LOG_COMPRESS"))

	hook := &lumberjack.Logger{
		Filename:   os.Getenv("LOG_DIR") + os.Getenv("LOG_FILENAME"),
		MaxSize:    logMaxSize,
		MaxBackups: logMaxBackups,
		MaxAge:     logMaxAge,
		Compress:   logCompress,
	}

	// Configuration of the log format based on the DKT format for ELK.
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "@level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "log_time",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(hook), infoLevel),
	)

	switch gin.Mode() {
	case gin.ReleaseMode:
		logger = zap.New(core, zap.AddCaller())
	default:
		logger = zap.New(core, zap.AddCaller())
		err := godotenv.Load()
		if err != nil {
			logger.Warn("Error loading .env file")
		}
	}
	defer logger.Sync()

	r := gin.Default()
	v1 := r.Group("/trigger-msg/v1/api")
	{
		v1.GET("/test", func(c *gin.Context) {
			delivery.TestAlive(c, logger)
		})
		v1.POST("/push", func(c *gin.Context) {
			delivery.ReceiveMessage(c, logger)
		})
		v1.POST("/sendstatus", func(c *gin.Context) {
			delivery.SendStatus(c, logger)
		})
		v1.GET("/getDir/:path", func(c *gin.Context) {
			delivery.DirService(c, logger)
		})
		v1.GET("/getInfoDir/:path", func(c *gin.Context) {
			delivery.DirInfoService(c, logger)
		})
	}
	r.Run(":3000")
}
