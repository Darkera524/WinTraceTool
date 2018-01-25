package trace

import (
	"io"
	"fmt"
	"github.com/bitly/go-simplejson"
	"strconv"
	"bytes"
	"net/http"
	"io/ioutil"
	"unsafe"
	"strings"
	"github.com/Darkera524/WinTraceTool/g"
)

type TraceBuilder struct{
	rawData map[string]interface{}
	formmattedData map[string]interface{}
}

func (t *TraceBuilder) getData() error {
	err := Checkpy()
	if err != nil {
		//差错处理
		return err
	}
	t.rawData = make(map[string]interface{})

	command := "python"
	params := []string{"trace.py", "Microsoft-Windows-DNS-Server", "EB79061A-A566-4698-9119-3ED2807060E7"}

	cmd,reader,err := ExecCommand(command, params)

	if err != nil {
		return err
	}

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		//fmt.Print(line)

		/*bytesline := []byte(line)
		for i, ch := range bytesline {

			switch {
			case ch > '~':   bytesline[i] = ' '
			case ch == '\r':
			case ch == '\n':
			case ch == '\t':
			case ch == '\'':	bytesline[i] = []byte("'")[0]
			case ch < ' ':   bytesline[i] = ' '
			}
		}*/

		go t.formatData(line)
	}

	cmd.Wait()
	return nil
}

func (t *TraceBuilder) formatData(line string) error {
	rs := []rune(string(line))
	length := len(rs)
	eventCode,err := strconv.Atoi(string(rs[1:4]))
	if err != nil {
		return err
	}

	//fmt.Println(string(rs[length-4:length]))

	json := string((rs[6:length-3]))

	jsonobj,err := simplejson.NewJson([]byte(json))
	if err != nil {
		return err
	}

	timestamp_win,err := jsonobj.Get("EventHeader").Get("TimeStamp").Int64()
	if err != nil {
		return err
	}
	timestamp_unix := int64((float64(timestamp_win)/10000000.0 - 11644473600.0) * 1000.0)
	//fmt.Println(taskname)

	//get the collection element tuple form server depending on providers and eventcode
	collectElements := []string{"Task Name", "QNAME", "Source"}

	metric := "trace.dns"
	datapoints := [][]int64{[]int64{timestamp_unix, 1}}

	send_json := simplejson.New()
	send_json.Set("name", metric)
	send_json.Set("datapoints", datapoints)
	tags_json := simplejson.New()
	tags_json.Set("Eventcode", eventCode)
	for i:=0;i<len(collectElements);i++ {
		ins,err := jsonobj.Get(collectElements[i]).String()
		if err != nil {
			//应当有log记录
			continue
		}
		tags_json.Set(strings.Replace(collectElements[i]," ", "", -1), ins)
	}
	send_json.Set("tags", tags_json)

	data,err := send_json.MarshalJSON()
	if err != nil {
		return err
	}

	t.sendData(data)

	return nil
}

func (t *TraceBuilder) sendData(data []byte) error {
	fmt.Println(string(data))
	reader := bytes.NewBuffer(data)

	//url从配置文件中读取
	url := g.GetConfig().KairosDB_path

	request,err := http.Post(url,"application/json;charset=utf-8", reader)
	if err != nil {
		return err
	}

	respBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)

	return nil
}

