package main

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

var Mgr map[string]*openwechat.Friend

func messageHandler(msg *openwechat.Message) {
	_, err := msg.Sender()
	if err != nil {
		fmt.Println(err)
		return
	}
	if msg.IsText() {
		if msg.Content == "ping" {
			msg.ReplyText("pong")
			return
		}
		if HandleMessage(msg.RawContent) {
			msg.ReplyText("鸡哥是最厉害的")
		} else {
			msg.ReplyText("哎呀 你干嘛!")
		}
	} else {
		fmt.Println("recv ", msg)
	}
}
func SendMsg(name, msg string) {
	f, ok := Mgr[name]
	if ok {
		f.SendText(msg)
	} else {
		fmt.Println("f = nil", name)
	}
}
func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func Init() {
	Mgr = map[string]*openwechat.Friend{}
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop)
	fmt.Println("ret")
	// 注册消息处理函数

	bot.MessageHandler = messageHandler
	// 可以设置通过该uuid获取到登录的二维码
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// 登录
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	err := bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())
	if err != nil {
		fmt.Println(err)
		return
	}
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(self)
	friends, err := self.Friends()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range friends {
		//	fmt.Println(f.RemarkName, f.Sex, f.NickName, f.UserName)
		Mgr[f.RemarkName] = f
	}
	go Run()

	// 阻塞主程序,直到用户退出或发生异常
	bot.Block()
}
