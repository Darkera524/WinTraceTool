package http

import (
	"net/http"
	"encoding/json"
	"github.com/Darkera524/WinTraceTool/model"
	"github.com/Darkera524/WinTraceTool/g"
)

func configNormalRoute() {
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
	})
	
	http.HandleFunc("/query", func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var provider []*model.HttpRequest
		err := decoder.Decode(&provider)
		if err != nil {

		}

		traceornot := false

		for k,_ := range g.ProviderMap.MetricMap {
			if k == provider[0].Provider {
				writer.Write([]byte("ok"))
				traceornot = true
				break
			}
		}

		if traceornot == false {
			writer.Write([]byte("no"))
		}
	})

}
