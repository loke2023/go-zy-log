package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

type DefaultFieldsHook struct {
	LogstashAddr 		string
	Name         		string
	LogstashEnv  		string
	LogErrToFeishu 		bool
}

type logmessage struct {
	Type  string `json:"type"`
	Time  int64  `json:"time"`
	Level string `json:"level"`
	Info  string `json:"info"`
}

var Log_LocalIP string = ""

func (df *DefaultFieldsHook) Fire(entry *logrus.Entry) error {

	//删除无效的utf-8字符
	if !utf8.ValidString(entry.Message) {
		v := make([]rune, 0, len(entry.Message))
		for i, r := range entry.Message {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(entry.Message[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		entry.Message = string(v)
	}

	var sigFile string
	_, file, line, ok := runtime.Caller(7)
	if ok {
		arrfile := strings.SplitAfterN(file, "/", 99)
		for i := 4; i >= 1; i-- {
			if i <= len(arrfile) {
				sigFile += arrfile[len(arrfile)-i]
			}
		}
		sigFile += ":" + strconv.Itoa(line)
	}
	entry.Data["file"] = sigFile
	logInfo := entry.Message

	if len(Log_LocalIP) <= 0 {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			Log_LocalIP = "0.0.0.0"
		} else {
			for _, address := range addrs {
				// 检查ip地址判断是否回环地址
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						Log_LocalIP = ipnet.IP.String()
					}
				}
			}
		}
	}

	//往自建日志服务器投递
	if len(df.LogstashAddr) > 0 {
		go func() {
			logKeyName := df.Name + "_" + df.LogstashEnv
			vvvv := logmessage{
				logKeyName,
				entry.Time.Unix(),
				entry.Level.String(),
				logInfo + ", file[" + sigFile + "]",
			}
			vInfo, err := json.Marshal(&vvvv)
			req, err := http.NewRequest("POST", df.LogstashAddr, bytes.NewBuffer(vInfo))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				fmt.Errorf(err.Error())
			}
			c := &http.Client{}
			resp, err := c.Do(req)
			if err != nil {
				fmt.Errorf(err.Error())
			}
			if resp != nil {
				defer resp.Body.Close()
			}
		}()
	}

	if df.LogErrToFeishu && entry.Level.String() == "error" {
		go func() {
			SendMessageToFeiShu(logInfo + ", file[" + sigFile + "]", nil)
		}()
	}

	return nil
}


func (df *DefaultFieldsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}


func LgPrintObject(logPre string, obj interface{}, lg *logrus.Logger) {
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			lgPrintObject(logPre + "- " + strconv.Itoa(i) + " ", v.Index(i).Interface(), lg)
		}
	default:
		lgPrintObject(logPre, obj, lg)
	}
}

func lgPrintObject(logPre string, obj interface{}, lg *logrus.Logger) {
	values := reflect.ValueOf(obj)
	types := reflect.TypeOf(obj)
	for i := 0; i < types.NumField(); i++  {
		lg.Info(logPre, "[", types.Name() ,"], [",
			types.Field(i).Name," = ", values.Field(i),"]")
	}
}