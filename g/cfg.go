package g

import (
	"github.com/toolkits/file"
	"fmt"
	"encoding/json"
	"sync"
	"os"
)

var (
	ConfigFile string
	config *Config
	lock = new(sync.RWMutex)
)



type Config struct {
	KairosDB_path	string		`json:"kairosDB_path"`
	Server_path		string		`json:"server_path"`
	Server_timeout 	int 		`json:"server_timeout"`
	Listen_port		string		`json:"listen_port"`
}

func GetConfig() *Config {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		fmt.Println(err.Error())
	}
	var c Config
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		fmt.Println(err.Error())
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

}

func Hostname() (string, error){
	hostname, err := os.Hostname()
	if err != nil {

	}
	return hostname, err
}
