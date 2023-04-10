package server

import (
	"go-zy-log/db"
	"time"

	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Addr     string
	Lg       *logrus.Logger
	Db       *db.DBControl
	MapCache cmap.ConcurrentMap
}

func New(lg *logrus.Logger, db *db.DBControl, addr string) (Server, error) {
	c := Server{
		Addr:     addr,
		Lg:       lg,
		Db:       db,
		MapCache: cmap.New(),
	}
	return c, nil
}

func (this *Server) QeuePushStringData(data string) {
	topic := "logs"
	existData, exist := this.MapCache.Pop(topic)
	if exist {
		arrData := existData.([]string)
		arrData = append(arrData, data)
		this.MapCache.Set(topic, arrData)
	} else {
		this.MapCache.Set(topic, []string{data})
	}
}

func (this *Server) QueuPopStringData() []string {
	topic := "logs"
	existData, exist := this.MapCache.Pop(topic)
	if exist {
		return existData.([]string)
	} else {
		return []string{}
	}
}

// @Summary 日志接收
// @Description 无返回，收到日志数据直接写入mongodb
// @Tags 常用接口
// @Accept json
// @Success 200 {string} string
// @param data body string true "日志数据"
// @Router /pub/logs [post]
func (c *Server) APIRecvLogs(g *gin.Context) {
	data, _ := g.GetRawData()
	c.QeuePushStringData(string(data))
}

func (c *Server) Start(r *gin.Engine) {

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v := r.Group("/pub")
	{
		v.POST("/logs", c.APIRecvLogs)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			logs := c.QueuPopStringData()
			for _, d := range logs {
				c.Db.PushLog([]byte(d))
			}
		}
	}()

	c.Db.CleanHistoryLogs()

	cronTask := cron.New() //创建一个cron实例
	cronTask.AddFunc("1 1 1 * *", func() {
		c.Db.CleanHistoryLogs()
	})

	//启动/关闭
	cronTask.Start()
	defer cronTask.Stop()

	//启动
	r.Run(c.Addr)
}
