package main

import (
	"go-gin-jwt-demo/config"
	"go-gin-jwt-demo/model"
	"go-gin-jwt-demo/router"
	"log/slog"
	"os"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ //NewJSONHandler formats each log entry as json, Stdout send logs to terminal or be collected by a journal
		Level: slog.LevelInfo, //LevelInfo includes INFO, WARN and ERROR messages
	}))
	slog.SetDefault(logger)
	slog.Info("application starting")

	// 初始化数据库连接
	config.InitDB()
	// 自动迁移建表（可选，不用手动执行SQL）
	err := config.DB.AutoMigrate(&model.User{})
	if err != nil {
		panic("数据表迁移失败：" + err.Error())
	}

	// 初始化路由
	r := router.SetupRouter()
	// 启动服务，监听8080端口
	_ = r.Run(":8080")
}
