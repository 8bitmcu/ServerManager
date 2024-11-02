package sm

import (
	"bufio"
	"os/exec"
	"path/filepath"
)

var lines string
var cmd *exec.Cmd

func Start() {
	cmd = exec.Command(filepath.Join(Dba.Basepath(), "server", "acServer"))
	cmd.Dir = filepath.Join(Dba.Basepath(), "server")
	var stdOut, _ = cmd.StdoutPipe()
	cmd.Start()

	out := bufio.NewReader(stdOut)

	lines = ""
	go func() {
		for cmd != nil {
			b, _, _ := out.ReadLine()
			lines += string(b) + "<br />"
		}
	}()

}
func Get_Content() string {
	return lines
}

func Is_Running() bool {
	if cmd != nil && cmd.ProcessState != nil {
		return cmd.ProcessState.Exited()
	} else if cmd != nil && cmd.Process != nil {
		return true
	}
	return false
}

func Stop() {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
		cmd = nil
	}
}
