package main

import (
	"flag"
	"go-zy-log/common"
	"go-zy-log/db"
	"go-zy-log/docs"
	"go-zy-log/server"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var lg = logrus.New()

var (
	listenAddr string
	mongodb    string
)

func init() {
	flag.StringVar(&listenAddr, "addr", ":25506", "MS server address")
	flag.StringVar(&mongodb, "mongodb", "mongodb://127.0.0.1:30017/logs", "")
}

// @title 日志服务
// @version 1.0
// @description 智能投放
// @BasePath /
func main() {

	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	dbControl := db.NewDBCon(mongodb, lg)
	s, _ := server.New(lg, &dbControl, listenAddr)
	common.SendMessageToFeiShu("go-zy-log - ["+listenAddr+"]，server start", lg)
	r := gin.Default() 
	docs.SwaggerInfo.BasePath = "/"

	s.Start(r)
}
