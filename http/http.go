package http

import (
	"github.com/Darkera524/WinTraceTool/g"
	"net/http"
	"fmt"
	_ "net/http/pprof"
	"encoding/json"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func init() {
	configNormalRoute()
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}

	RenderDataJson(w, data)
}


func Start() {
	port := g.GetConfig().Listen_port

	fmt.Println(port)

	s := &http.Server{
		Addr: port,
		MaxHeaderBytes: 1 << 30,
	}

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
