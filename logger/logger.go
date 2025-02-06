package logger

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var LOGGER = Logger{
	AlsoPrint: true,
}
var LOG_OUT = "./locksmith.log"

func getLogFile() *os.File {
	file, err := os.OpenFile(LOG_OUT, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		if os.Getenv("DO_NOT_PANIC_LOCKSMITH") != "true" {
			panic("error opening log file")
		}
	}

	return file
}

var (
	buf    bytes.Buffer
	logger = log.New(getLogFile(), "", log.Lmsgprefix)
)

type LogFormat struct {
	XMLName    xml.Name `xml:"decoder"`
	WazuhName  string   `xml:"name,attr"`
	Parent     string   `xml:"parent"`
	Regex      string   `xml:"regex"`
	RegexOrder []string `xml:"order>value"`
	FmtPattern string   `xml:"-"`
}

// FORMAT:
// timestamp hostname level src-ip: message
/*
<decoder name="locksmith">
  <program_name>^locksmith</program_name>
</decoder>
<decoder name="locksmith">
  <parent>locksmith</parent>
  <regex>User '(\w+)' logged from '(\d+.\d+.\d+.\d+)'</regex>
  <order>user, srcip</order>
</decoder>
<decoder name="locksmith">
  <regex>User '(\w+)' failed login from '(\d+.\d+.\d+.\d+)'</regex>
  <order>user, srcip</order>
</decoder>
*/

func appendToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

type Logger struct {
	AlsoPrint bool
}

func (l *Logger) Log(useFormat LogFormats, params ...interface{}) {
	selectedLog := Logs[useFormat]
	if len(params) != len(selectedLog.RegexOrder) {
		panic(fmt.Sprintf("too incorrect number of arguments passed to %+v - %+v\n", useFormat, params))
	}
	format := selectedLog.FmtPattern

	out := fmt.Sprintf(format, params...)

	// Jun 13 11:12:13 server locksmith[1111]:
	timestamp := time.Now().Format("Jan 02 15:04:05")
	processName := "locksmith"
	hostname, _ := os.Hostname()
	processID := os.Getpid()

	out = fmt.Sprintf("%s %s %s[%d]: %s", timestamp, hostname, processName, processID, out)
	logger.Println(out)

	if l.AlsoPrint {
		fmt.Println(out)
	}
}

func GetIPFromRequest(r http.Request) string {
	// Support for reverse proxies
	if forwardedIP := r.Header.Get("X-Forwarded-For"); len(forwardedIP) > 0 {
		return forwardedIP
	}

	return r.RemoteAddr
}
