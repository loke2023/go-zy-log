package common

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var globle int64
type FlagTime struct {
	Flag 		string
	TimeMS      int64
}

type AutoFuncTime struct {
	Lg 			*logrus.Logger
	Function    string
	Flags 		[]FlagTime
}

func (a *AutoFuncTime) Begin(fun string, lg *logrus.Logger){
	a.Function = fun
	a.Lg = lg
	a.Flags = append(a.Flags, FlagTime{
		"Begin",
		time.Now().UnixNano()/1e6,
	})
}

func (a *AutoFuncTime) BeginUser(fun string, lg *logrus.Logger, userID int64){
	a.Function = fun
	a.Lg = lg
	a.Flags = append(a.Flags, FlagTime{
		"Begin",
		time.Now().UnixNano()/1e6,
	})
}

func (a *AutoFuncTime) BeginUserS(fun string, lg *logrus.Logger, userID string){
	a.Function = fun
	a.Lg = lg
	a.Flags = append(a.Flags, FlagTime{
		"Begin",
		time.Now().UnixNano()/1e6,
	})
}

func (a *AutoFuncTime) Mark(f string){
	a.Flags = append(a.Flags, FlagTime{
		f,
		time.Now().UnixNano()/1e6,
	})
}


func (a *AutoFuncTime) End(s string){
	a.Flags = append(a.Flags, FlagTime{
		"End",
		time.Now().UnixNano()/1e6,
	})
	beginTime := int64(0)
	PreFlag := FlagTime{}
	sExtern := "Detail["
	for i, v := range a.Flags {
		if i == 0 {
			beginTime = v.TimeMS
		}else{
			if i > 1 {
				sExtern += ","
			}
			sExtern += v.Flag + ":" + strconv.FormatInt(v.TimeMS-PreFlag.TimeMS,10) + "ms"
		}
		PreFlag = v
	}
	sExtern += "]"
	ms := time.Now().UnixNano()/1e6 - beginTime
	if a.Lg != nil {
		tExTime := int64(0)
		if ms > 5000 + tExTime {
			a.Lg.Error( "*****Error*****" + s + "[globle:" + strconv.FormatInt(globle, 10) +
				"]函数[" + a.Function + "]处理耗时[" + strconv.FormatInt(ms,10) + "ms], " + sExtern)
		} else if ms > 2000 + tExTime{
			a.Lg.Warn( "*****Danger*****" + s + "[globle:" + strconv.FormatInt(globle, 10) +
				"]函数[" + a.Function + "]处理耗时[" + strconv.FormatInt(ms,10) + "ms], " + sExtern)
		} else if ms >= 900 + tExTime {
			a.Lg.Warn("*****Warn*****" + s + "[globle:" + strconv.FormatInt(globle, 10) +
				"]函数[" + a.Function + "]处理耗时[" + strconv.FormatInt(ms,10) + "ms], " + sExtern)
		}
		globle++
	}
}


func (a *AutoFuncTime) Print(s string, toFs bool){
	a.Flags = append(a.Flags, FlagTime{
		"End",
		time.Now().UnixNano()/1e6,
	})
	beginTime := int64(0)
	PreFlag := FlagTime{}
	sExtern := "Detail["
	for i, v := range a.Flags {
		if i == 0 {
			beginTime = v.TimeMS
		}else{
			if i > 1 {
				sExtern += ","
			}
			sExtern += v.Flag + ":" + strconv.FormatInt(v.TimeMS-PreFlag.TimeMS,10) + "ms"
		}
		PreFlag = v
	}
	sExtern += "]"
	ms := time.Now().UnixNano()/1e6 - beginTime
	sLogInfo := ""
	if a.Lg != nil {
		sLogInfo = s + "-[" + a.Function + "]处理耗时[" + strconv.FormatInt(ms,10) + "ms]"
		a.Lg.Info(sLogInfo + ", " + sExtern)
		globle++
	}
	if toFs {
		SendMessageToFeiShu(sLogInfo + ", " + sExtern, a.Lg)
	}
}

