package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

var Mgr map[string]*openwechat.Friend
var Self *openwechat.Self
var MsgChan chan MsgType

func messageHandler(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		fmt.Println("msg error", err)
		return
	}
	if msg.IsText() && sender.IsFriend() {
		if msg.Content == "ping" {
			msg.ReplyText("pong")
			return
		}
		fmt.Println(sender.ID(), sender.RemarkName, sender.NickName, sender == Self.User)
		b, s := getFollow(sender.ID()).HandleMessage(msg.RawContent)
		if b {
			msg.ReplyText(s)
		} else {
			//msg.ReplyText("aa")
		}

	} else {
		fmt.Println("recv ", msg)
	}
}
func SendMsg(id, msg string) {
	MsgChan <- MsgType{id, msg}
}
func RecvMsg() {
	tick := time.NewTicker(10 * time.Second)
	cache := map[string]string{}
	for {
		select {
		case <-tick.C:
			for id, msg := range cache {
				_sendMsg(id, msg)
			}
			cache = map[string]string{}
		case x := <-MsgChan:
			cache[x.Id] += x.Msg
		}
	}
}
func _sendMsg(id, msg string) {
	f, ok := Mgr[id]
	if ok {
		f.SendText(msg)
	} else {
		fmt.Println("f = nil", id)
	}
}
func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func RunWechat() {
	Mgr = map[string]*openwechat.Friend{}
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop)
	fmt.Println("ret")
	// 注册消息处理函数

	bot.MessageHandler = messageHandler
	// 可以设置通过该uuid获取到登录的二维码
	if runtime.GOOS == "linux" {
		bot.UUIDCallback = ConsoleQrCode
	} else if runtime.GOOS == "windows" {
		bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	}
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// 登录
	reloadStorage := openwechat.NewFileHotReloadStorage("../storage.json")
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
		//	fmt.Println(f.RemarkName, f.NickName)
		Mgr[f.ID()] = f
	}
	go RecvMsg()
	// 阻塞主程序,直到用户退出或发生异常
	StartWs <- struct{}{}
	bot.Block()
}
