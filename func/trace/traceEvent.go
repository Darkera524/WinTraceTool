package trace

import (
	"io"
	"github.com/bitly/go-simplejson"
	"strconv"
	"bytes"
	"net/http"
	"io/ioutil"
	"unsafe"
	"github.com/Darkera524/WinTraceTool/g"
	"github.com/Darkera524/WinTraceTool/model"
	"time"
)

type TraceBuilder struct{
	/*rawData map[string]interface{}
	formmattedData map[string]interface{}*/
}

//var collectElements []string

func (t *TraceBuilder) getData(guid string) error {
	err := Checkpy()
	if err != nil {
		//差错处理
		g.Logger().Println(err.Error())
		return err
	}
	//t.rawData = make(map[string]interface{})

	command := "python"
	params := []string{"trace.py", g.ProviderMap.MetricMap[guid], guid}

	//
	g.Logger().Println(g.ProviderMap.MetricMap[guid])

	//collectElements = g.ProviderMap.Providers[guid]
	model.TracingList = append(model.TracingList, guid)
	cmd,reader,err := ExecCommand(command, params)

	if err != nil {
		g.Logger().Println(err.Error())
		model.TracingList = remove(model.TracingList, guid)
		return err
	}

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			model.TracingList = remove(model.TracingList, guid)
			break
		}

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
		//g.Logger().Println(line)
		go t.formatData(line)//, g.ProviderMap.MetricMap[guid])
	}

	cmd.Wait()
	return nil
}

func (t *TraceBuilder) formatData(line string) error {
	rs := []rune(string(line))
	ins_bytes := []byte(line)
	length := len(rs)

	var coma_index int
	for i,v := range ins_bytes{
		if v == ','{
			coma_index = i
			break
		}
	}

	eventCode,err := strconv.Atoi(string(rs[1:coma_index]))
	if err != nil {
		g.Logger().Println(err.Error())
		return err
	}

	//fmt.Println(string(rs[length-4:length]))

	json := string((rs[coma_index+1:length-3]))

	jsonobj,err := simplejson.NewJson([]byte(json))
	if err != nil {
		g.Logger().Println(err.Error())
		return err
	}

	timestamp_win,err := jsonobj.Get("EventHeader").Get("TimeStamp").Int64()
	if err != nil {
		g.Logger().Println(err.Error())
		return err
	}
	timestamp_unix := int64((float64(timestamp_win)/10000000.0 - 11644473600.0) * 1000.0)
	tm := time.Unix(timestamp_unix/1000, (timestamp_unix%1000)*1000000)
	time_string := tm.Format("2006-01-02 15:04:05.999999-07:00")

	switch eventCode {
	case 256:
		//query_received
		go t.handle_dns_256(time_string, jsonobj)
	case 257:
		//response_success
		go t.handle_dns_257(time_string, jsonobj)
	case 258:
		//response_failure
		go t.handle_dns_258(time_string, jsonobj)
	case 259:
		//ignored_query
		go t.handle_dns_259(time_string, jsonobj)
	case 260:
		//recurse_query_out
		go t.handle_dns_260(time_string, jsonobj)
	case 261:
		//recurse_response_in
		go t.handle_dns_261(time_string, jsonobj)
	case 262:
		//recurse_query_timeout
		go t.handle_dns_262(time_string, jsonobj)
	default:
		//other event that eventid > 262
		go t.handle_dns_other(time_string, jsonobj)
	}

	//fmt.Println(taskname)

	/*metric := "trace." + provider
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

	t.appendSql(send_json)

	data,err := send_json.MarshalJSON()
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	err = t.sendData(data)
	if err != nil {
		return err
	}
	*/

	return nil
}

//send to kairosDB
func (t *TraceBuilder) sendData(data []byte) error {
	reader := bytes.NewBuffer(data)
	//url从配置文件中读取
	url := g.GetConfig().KairosDB_path

	request,err := http.Post(url,"application/json;charset=utf-8", reader)
	if err != nil {
		g.Logger().Println(err.Error())
	}

	respBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		g.Logger().Println(err.Error())
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	g.Logger().Println(*str)

	return nil
}

func remove(slice []string, elem string) []string {
	for i,v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (t *TraceBuilder)handle_dns_256(time_string string,jsonobj *simplejson.Json){
	//InterfaceIP、Source、QNAME

	interface_ip,err := jsonobj.Get("InterfaceIP").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	source,err := jsonobj.Get("Source").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	qname,err := jsonobj.Get("QNAME").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Insert into dns_256(eventtime, eventcode, interfaceip, source, qname, dnsserver) VALUES")
	buffer.WriteString("('")
	buffer.WriteString(time_string)
	buffer.WriteString("','")
	buffer.WriteString("256")
	buffer.WriteString("','")
	buffer.WriteString(interface_ip)
	buffer.WriteString("','")
	buffer.WriteString(source)
	buffer.WriteString("','")
	buffer.WriteString(qname)
	buffer.WriteString("','")
	hostname,err := g.Hostname()
	buffer.WriteString(hostname)
	buffer.WriteString("')")

	sql := buffer.String()

	//exec sql
	err = g.ExecSql(sql)
	if err != nil {
		g.Logger().Println("Error occurred when exec the sql" + sql)
	}
}

func (t *TraceBuilder)handle_dns_257(time_string string,jsonobj *simplejson.Json){
	//InterfaceIP Destination QNAME
	interface_ip,err := jsonobj.Get("InterfaceIP").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	destination,err := jsonobj.Get("Destination").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	qname,err := jsonobj.Get("QNAME").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Insert into dns_257(eventtime, eventcode, interfaceip, destination, qname, dnsserver) VALUES")
	buffer.WriteString("('")
	buffer.WriteString(time_string)
	buffer.WriteString("','")
	buffer.WriteString("257")
	buffer.WriteString("','")
	buffer.WriteString(interface_ip)
	buffer.WriteString("','")
	buffer.WriteString(destination)
	buffer.WriteString("','")
	buffer.WriteString(qname)
	buffer.WriteString("','")
	hostname,err := g.Hostname()
	buffer.WriteString(hostname)
	buffer.WriteString("')")

	sql := buffer.String()

	//exec sql
	err = g.ExecSql(sql)
}

func (t *TraceBuilder)handle_dns_258(time_string string,jsonobj *simplejson.Json){
	//InterfaceIP Reason Destination QNAME
	interface_ip,err := jsonobj.Get("InterfaceIP").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	reason,err := jsonobj.Get("Reason").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	destination,err := jsonobj.Get("Destination").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	qname,err := jsonobj.Get("QNAME").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Insert into dns_258(eventtime, eventcode, interfaceip, reason, destination, qname, dnsserver) VALUES")
	buffer.WriteString("('")
	buffer.WriteString(time_string)
	buffer.WriteString("','")
	buffer.WriteString("258")
	buffer.WriteString("','")
	buffer.WriteString(interface_ip)
	buffer.WriteString("','")
	buffer.WriteString(reason)
	buffer.WriteString("','")
	buffer.WriteString(destination)
	buffer.WriteString("','")
	buffer.WriteString(qname)
	buffer.WriteString("','")
	hostname,err := g.Hostname()
	buffer.WriteString(hostname)
	buffer.WriteString("')")

	sql := g.Enc.ConvertString(buffer.String())

	//exec sql
	err = g.ExecSql(sql)
}

func (t *TraceBuilder)handle_dns_259(time_string string,jsonobj *simplejson.Json){
	//InterfaceIP Reason QNAME
	interface_ip,err := jsonobj.Get("InterfaceIP").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	reason,err := jsonobj.Get("Reason").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	qname,err := jsonobj.Get("QNAME").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Insert into dns_259(eventtime, eventcode, interfaceip, reason, qname, dnsserver) VALUES")
	buffer.WriteString("('")
	buffer.WriteString(time_string)
	buffer.WriteString("','")
	buffer.WriteString("259")
	buffer.WriteString("','")
	buffer.WriteString(interface_ip)
	buffer.WriteString("','")
	buffer.WriteString(reason)
	buffer.WriteString("','")
	buffer.WriteString(qname)
	buffer.WriteString("','")
	hostname,err := g.Hostname()
	buffer.WriteString(hostname)
	buffer.WriteString("')")

	sql := g.Enc.ConvertString(buffer.String())

	//exec sql
	err = g.ExecSql(sql)
}

func (t *TraceBuilder)handle_dns_260(time_string string,jsonobj *simplejson.Json){
	//InterfaceIP Destination QNAME
	interface_ip,err := jsonobj.Get("InterfaceIP").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	destination,err := jsonobj.Get("Destination").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	qname,err := jsonobj.Get("QNAME").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Insert into dns_260(eventtime, eventcode, interfaceip, destination, qname, dnsserver) VALUES")
	buffer.WriteString("('")
	buffer.WriteString(time_string)
	buffer.WriteString("','")
	buffer.WriteString("260")
	buffer.WriteString("','")
	buffer.WriteString(interface_ip)
	buffer.WriteString("','")
	buffer.WriteString(destination)
	buffer.WriteString("','")
	buffer.WriteString(qname)
	buffer.WriteString("','")
	hostname,err := g.Hostname()
	buffer.WriteString(hostname)
	buffer.WriteString("')")

	sql := buffer.String()

	//exec sql
	err = g.ExecSql(sql)
}

func (t *TraceBuilder)handle_dns_261(time_string string,jsonobj *simplejson.Json){
	//InterfaceIP Source QNAME
	interface_ip,err := jsonobj.Get("InterfaceIP").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	source,err := jsonobj.Get("Source").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	qname,err := jsonobj.Get("QNAME").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Insert into dns_261(eventtime, eventcode, interfaceip, source, qname, dnsserver) VALUES")
	buffer.WriteString("('")
	buffer.WriteString(time_string)
	buffer.WriteString("','")
	buffer.WriteString("261")
	buffer.WriteString("','")
	buffer.WriteString(interface_ip)
	buffer.WriteString("','")
	buffer.WriteString(source)
	buffer.WriteString("','")
	buffer.WriteString(qname)
	buffer.WriteString("','")
	hostname,err := g.Hostname()
	buffer.WriteString(hostname)
	buffer.WriteString("')")

	sql := buffer.String()

	//exec sql
	err = g.ExecSql(sql)
}

func (t *TraceBuilder)handle_dns_262(time_string string,jsonobj *simplejson.Json){
	//InterfaceIP Destination QNAME
	interface_ip,err := jsonobj.Get("InterfaceIP").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	destination,err := jsonobj.Get("Destination").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}
	qname,err := jsonobj.Get("QNAME").String()
	if err != nil {
		//log
		g.Logger().Println(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Insert into dns_262(eventtime, eventcode, interfaceip, destination, qname, dnsserver) VALUES")
	buffer.WriteString("('")
	buffer.WriteString(time_string)
	buffer.WriteString("','")
	buffer.WriteString("262")
	buffer.WriteString("','")
	buffer.WriteString(interface_ip)
	buffer.WriteString("','")
	buffer.WriteString(destination)
	buffer.WriteString("','")
	buffer.WriteString(qname)
	buffer.WriteString("','")
	hostname,err := g.Hostname()
	buffer.WriteString(hostname)
	buffer.WriteString("')")

	sql := buffer.String()

	//exec sql
	err = g.ExecSql(sql)
}

func (t *TraceBuilder)handle_dns_other(time_string string,jsonobj *simplejson.Json){

}



