package wmi_collector

import (
	"github.com/Darkera524/WinTraceTool/model"
	"github.com/StackExchange/wmi"
	"reflect"

	"bytes"
	"strconv"
	"github.com/Darkera524/WinTraceTool/g"
	"strings"
	"time"
)

type WmiInfos struct {

}

func (w *WmiInfos) getData(wmiClass string) error {

	switch wmiClass {
	case "Win32_PerfRawData_perfProc_Process":
		model.WmiInfoList = append(model.WmiInfoList, wmiClass)
		class := []model.Win32_PerfRawData_perfProc_Process{}
		q := wmi.CreateQuery(&class, "")
		err := wmi.Query(q, &class)
		if err != nil {
			g.Logger().Println(err.Error())
		}

		for _,v := range class {
			go w.formatData(v)
		}
	}



	return nil
}

func (w *WmiInfos) formatData(info model.Wmi_Base) error {
	detect_type := reflect.TypeOf(info)
	field_num := detect_type.NumField()
	field_list := []string{}
	for i:=0;i<field_num;i++ {
		field_list = append(field_list, detect_type.Field(i).Name)
	}

	type_total := strings.Split(detect_type.String(),".")
	type_name := type_total[len(type_total)-1]

	detect_value := reflect.ValueOf(info)

	sql := generateSql(type_name, field_list, detect_value)

	w.sendData(sql)

	return nil
}

func (w *WmiInfos) sendData(sql string) error {
	err := g.ExecSql(sql)
	if err != nil {
		g.Logger().Println("Error occurred when exec the sql" + sql)
	}

	return nil
}

func generateSql(name string, field_list []string, detect_value reflect.Value) string {
	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(name)
	buffer.WriteString(" (")
	for i:=0;i<len(field_list);i++ {
		buffer.WriteString(field_list[i])
		if i != len(field_list)-1 {
			buffer.WriteString(",")
		}
 	}
 	buffer.WriteString(",eventtime")
 	buffer.WriteString(") VALUES ('")
 	for i:=0;i<len(field_list);i++ {
 		field := detect_value.FieldByName(field_list[i])
 		if field.Type().String() == "string" {
			buffer.WriteString(field.String())
		} else if field.Type().String() == "int64" {
			buffer.WriteString(strconv.FormatInt(field.Int(),10))
		} else if field.Type().String() == "float64" {
			buffer.WriteString(strconv.FormatFloat(field.Float(), 'g', 15, 64))
		} else {
			g.Logger().Println("uninitialed type:" + field.Type().String())
		}

 		if i != len(field_list)-1 {
 			buffer.WriteString("','")
		} else {
			buffer.WriteString("','")
			ins := strings.Split(time.Now().String(),".")[0]
			buffer.WriteString(ins)
			buffer.WriteString("')")
		}
	}

	return buffer.String()
}