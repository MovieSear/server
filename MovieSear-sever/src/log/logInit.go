package logger

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("controllers")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {

	// 程序入口main.go函数，log.txt文件相对main.go函数的路径为./log/
	// 文件不存在则创建一个；文件属性为0666，表示可读可写
	logFile, err := os.OpenFile("./log/log.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open log file failed: ", err)
	}

	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(logFile, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)

}

