package trace

import (
	"fmt"
	"bufio"
	"github.com/toolkits/file"
	"errors"
	"os/exec"
)

func Checkpy() error {
	pyfile := "trace.py"
	if !file.IsExist(pyfile) {
		return errors.New("trace.py is not found")
	}
	return nil
}

func ExecCommand(commandName string, params []string) (*exec.Cmd, *bufio.Reader,error) {
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	//fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return nil,nil,err
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	return cmd,reader,nil


}
