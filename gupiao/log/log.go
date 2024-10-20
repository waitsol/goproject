package log2

import (
	"bytes"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

type Formatter struct {
}

func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	buffer := entry.Buffer
	if buffer == nil {
		tmp := []byte{}
		buffer = bytes.NewBuffer(tmp)
	}
	//[时间][函数][文件][行]:data
	line := fmt.Sprintf("[%s][%s][%s][%d]:%s\n", entry.Time.Format(time.TimeOnly), path.Base(entry.Caller.Func.Name()),
		path.Base(entry.Caller.File), entry.Caller.Line, entry.Message,
	)
	buffer.WriteString(line)
	return buffer.Bytes(), nil
}
func init() {
	logfile := "log/server_"
	linkName := "log/latest_log.log"
	// 日志文件保留的时间

	// 创建日志轮转器
	rotator, err := rotatelogs.New(
		logfile+"%Y%m%d_%H.log",                          // 日志文件名加时间
		rotatelogs.WithLinkName(linkName),                // 始终指向最新的日志文件
		rotatelogs.WithRotationTime(time.Second*60*60*2), //2小时
		rotatelogs.WithMaxAge(time.Second*60*60*24*7),    //7天
	)
	if err != nil {
		log.Fatal("Failed to create rotator: ", err)
	}
	log.SetReportCaller(true)
	log.SetFormatter(&Formatter{})

	multiWriter := io.MultiWriter(os.Stdout, rotator)
	// 设置日志输出到文件
	log.SetOutput(multiWriter)
	// 打印日志
	log.Info("This is a log message.")
}
