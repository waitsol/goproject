package notify

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/waitsol/golib"
)

type MsgRes struct {
	Res string `json:"res"`
}

var zhanghao map[string]string
var ip2id map[string]string

func init() {
	ip2id = map[string]string{}
	zhanghao = map[string]string{}
	zhanghao["15358698379"] = "1"
	zhanghao["13910692031"] = "1"
	zhanghao["18343007254"] = "1"

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
func innerHandler(c *gin.Context, id string, args ...string) {
	msg := ""
	for i, x := range args {
		if i != 0 {
			msg += "_"
		}
		msg += x
	}
	b, res := GetFollow(id).HandleMessage(msg)
	fmt.Println("ret", res)
	if b {
		h := strings.ReplaceAll(res, "\n", "    ")
		fmt.Println(h)
		c.JSON(200, MsgRes{Res: h})

	} else {
		c.JSON(417, MsgRes{Res: "什么玩意"})
	}
}
func login(c *gin.Context) {
	user := c.PostForm("username")
	passwd := c.PostForm("password")
	log.Info(user, passwd)
	ip := c.Request.RemoteAddr
	fmt.Println("ip = ", ip)
	fmt.Println(c.Request.Header)
	if x, ok := zhanghao[user]; ok && x == passwd || x == "1" {
		if x == "1" {
			user = "15358698379"
		}
		ip2id[ip] = user
		fmt.Println(c.Keys)
		c.HTML(http.StatusOK, "app.html", gin.H{})
	} else {
		c.JSON(417, MsgRes{Res: "密码xx"})

	}
}
func list(c *gin.Context) {

	if x, ok := ip2id[c.Request.RemoteAddr]; ok {
		innerHandler(c, x, "list")
	} else {
		c.JSON(200, "登录过期")
	}
}

func clear(c *gin.Context) {
	if x, ok := ip2id[c.Request.RemoteAddr]; ok {
		innerHandler(c, x, "clear")
	} else {
		c.JSON(200, "账号过期")
	}
}

func set(c *gin.Context) {

	if x, ok := ip2id[c.Request.RemoteAddr]; ok {
		innerHandler(c, x, "set", c.PostForm("id"), c.PostForm("count"))
	} else {
		c.JSON(200, "账号过期")
	}

}
func add(c *gin.Context) {
	gpid := c.PostForm("id")

	if x, ok := ip2id[c.Request.RemoteAddr]; ok {
		innerHandler(c, x, gpid)
	} else {
		c.JSON(200, "账号过期")
	}

}
func Run() {
	loadFollow()
	router := gin.Default()
	router.LoadHTMLGlob("v1/*.html")

	v1 := router.Group("/v1")
	v1.GET("handle", handle)
	v1.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	v1.POST("/login", login)
	v1.POST("/add", add)
	v1.POST("/list", list)
	v1.POST("/clear", clear)
	v1.POST("/set", set)
	golib.Go(func() {
		router.Run("0.0.0.0:9876")
	})
}
