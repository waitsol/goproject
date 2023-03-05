package main

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

var Mgr map[string]*openwechat.Friend
var Self *openwechat.Self

func messageHandler(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		fmt.Println(err)
		return
	}
	if msg.IsText() && sender.IsFriend() {
		if msg.Content == "ping" {
			msg.ReplyText("pong")
			return
		}
		SendMsg("by2", "ss")
		fmt.Println(sender.RemarkName, sender.NickName, sender == Self.User)
		b, s := getFoller(sender.RemarkName).HandleMessage(msg.RawContent)
		if b {
			msg.ReplyText(s)
		} else {
			msg.ReplyText("aa")
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
	Self = self
	fmt.Println(self)
	friends, err := self.Friends()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range friends {
		//fmt.Println(f.RemarkName, f.Sex, f.NickName, f.UserName)
		fmt.Println(f.RemarkName, f.NickName)
		Mgr[f.RemarkName] = f
	}
	//go Run()

	// 阻塞主程序,直到用户退出或发生异常
	bot.Block()
}
