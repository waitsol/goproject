package TcpNet

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"main/ipc"
	"main/proto/go/main/pb3"
	"net"
	"time"
	"unsafe"
)

func sim_send() {
	code := ""
	Num := ""
	Money := ""
	buy := false
	bs := &pb3.CSBuyOrSellReq{Code: &code, Num: &Num, Money: &Money, IsBuy: &buy}
	req := &pb3.PacketReq{Packet: &pb3.PacketReq_Bs{Bs: bs}}

	time.Sleep(time.Second * 2)
	*bs.Code = "603716"
	*bs.Money = "33.67"
	*bs.Num = "100"
	*bs.IsBuy = true
	log.Info("......chan.....")
	ipc.ReqChan <- *req
	log.Info("......chan..... ok ")

	time.Sleep(time.Second * 2)
	*bs.Code = "301078"
	*bs.Money = "13"
	*bs.Num = "100"
	*bs.IsBuy = true
	ipc.ReqChan <- *req

	time.Sleep(time.Second * 2)
	*bs.Code = "600203"
	*bs.Money = "15.2"
	*bs.Num = "100"
	*bs.IsBuy = false
	ipc.ReqChan <- *req

	time.Sleep(time.Second * 2)
	*bs.Code = "003023"
	*bs.Money = "23"
	*bs.Num = "100"
	*bs.IsBuy = false
	ipc.ReqChan <- *req

}
func Run() {
	log.Infof("启动 TCP 服务器...")

	// 监听本地 8080 端口
	listener, err := net.Listen("tcp", ":9654")
	if err != nil {
		log.Infof("监听失败:", err)
		return
	}
	defer listener.Close()

	log.Infof("服务器正在等待连接...")

	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Infof("接受连接失败:", err)
			continue
		}
		//go sim_send()
		// 为每个连接创建 goroutine 处理
		handleConnection(conn)
	}
}
func read200(conn net.Conn) []byte {
	buf := make([]byte, 1024)
	c := 0
	for c < 200 {
		n, err := conn.Read(buf[c:])
		log.Info("n = ", n, "  c = ", c, err)
		c += n
		if err != nil {
			log.Error("read error ", err, " c = ", c)
			return nil
		}
	}
	return buf
}
func readMsg(conn net.Conn) []byte {
	buf := read200(conn)
	if buf == nil {
		return nil
	}
	ptr := (*int32)(unsafe.Pointer(&buf[0]))
	// 取值
	value := *ptr

	return buf[4 : value+4]
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Infof("客户端 %s 已连接\n", conn.RemoteAddr().String())

	data := read200(conn)

	if data == nil {
		log.Error("data nil")
		return
	}
	int64Ptr := (*int64)(unsafe.Pointer(&data[0]))
	// 取值
	value := *int64Ptr
	if value != 0x9865743211 {
		log.Error("验证失败 close")
		return
	}
	log.Info("验证通过")
	for {

		req := <-ipc.ReqChan
		log.Info("chan ... recv ", req.String())
		arr, err := proto.Marshal(&req)
		buf := make([]byte, 200)
		*(*int32)(unsafe.Pointer(&buf[0])) = int32(len(arr))
		copy(buf[4:], arr)
		// 发送给ths
		log.Info("send ths ..........", req.String())
		_, err = conn.Write(buf[0:200])
		if err != nil {
			log.Error("发送回复失败:", err)
			return
		}
		// 读取客户端发送的数据
		data = readMsg(conn)
		if data == nil {
			log.Errorf("客户端 %s 断开连接\n", conn.RemoteAddr().String())
			return
		}
		var rsp pb3.PacketRsp
		err = proto.Unmarshal(data, &rsp)
		if err != nil {
			log.Error("Unmarshal msg  error ", err)
			return
		}
		log.Info("recv ...........   ", rsp.String())
	}
}
