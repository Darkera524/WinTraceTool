package g

import "time"

var (
	HbsClient *SingleConnRpcClient
)

func InitRpcClients() {
		HbsClient = &SingleConnRpcClient{
			RpcServer: GetConfig().Server_path,
			Timeout:   time.Duration(GetConfig().Server_timeout) * time.Millisecond,
		}

}
