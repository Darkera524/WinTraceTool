package g

import (
	"time"
	"os"
	"log"
	"github.com/axgle/mahonia"
)

var (
	HbsClient *SingleConnRpcClient
	logger *log.Logger
	Enc mahonia.Decoder
)

func InitRpcClients() {
		HbsClient = &SingleConnRpcClient{
			RpcServer: GetConfig().Server_path,
			Timeout:   time.Duration(GetConfig().Server_timeout) * time.Millisecond,
		}

}

func InitLog() {
	fileName := "wintracetool.log"
	logFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("open file error !")
	}
	logger = log.New(logFile, "[Debug]", log.LstdFlags)
	log.Println("logging on", fileName)
}

func Logger() *log.Logger {
	lock.RLock()
	defer lock.RUnlock()
	return logger
}

//convert GBK to utf-8
func InitEncoding(){
	Enc = mahonia.NewDecoder("gbk")
}
