package wechat

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
)

type MsgRes struct {
	Res string `json:"res"`
}

func handle(c *gin.Context) {
	id := c.Query("id")
	msg := c.Query("msg")
	log.Info(msg)
	b, res := GetFollow(id).HandleMessage(msg)
	if b {
		c.JSON(200, MsgRes{Res: res})
	} else {
		c.JSON(417, MsgRes{Res: "什么玩意"})
	}
}

func Run() {
	loadFollow()
	router := gin.Default()
	v1 := router.Group("/v1")
	v1.GET("handle", handle)
	golib.Go(func() {
		router.Run()
	})
}
