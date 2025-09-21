package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/ztrue/tracerr"
)

// open opens the specified URL in the default browser of the user.
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func PrintInterface(t interface{}) {
	s, _ := json.MarshalIndent(t, "", "\t")
	fmt.Print(string(s))
}

func FormatErrorHTML(err error) template.HTML {
	newErr := FormatError(err)

	r := regexp.MustCompile("\n")
	newErr = r.ReplaceAllString(newErr, "<br />")

	return template.HTML(newErr)
}

func FormatError(err error) string {
	newErr := tracerr.Sprint(err)

	r2 := regexp.MustCompile(".*ServerManager/")
	newErr = r2.ReplaceAllString(newErr, "")

	r3 := regexp.MustCompile(".*pkg/mod/")
	newErr = r3.ReplaceAllString(newErr, "")

	return newErr
}
