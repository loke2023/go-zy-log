package common

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type QueryObj struct {
	mu sync.Mutex
	v  map[int]string
}

func (c *QueryObj)Init(i int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v = make(map[int]string)
	for ii := 1; ii <= i; ii++  {
		c.v[ii] = "free"
	}
}

func (c *QueryObj)PopFree() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.v {
		if v == "free" {
			c.v[k] = "use"
			return k
		}
	}
	return 0
}

func (c *QueryObj)RecvFree(i int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.v {
		if k == i && v == "use" {
			c.v[k] = "free"
			break
		}
	}
}

type StringTime struct {
	S 	string
	T 	time.Time
}

func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return nil == err
}

type SafeMapStringArray struct {
	mu sync.Mutex
	v  map[string][]StringTime
}

// Inc increments the counter for the given key.
func (c *SafeMapStringArray) Push(key string, s StringTime) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key] = append(c.v[key], s)
	defer c.mu.Unlock()
}


// Value returns the current value of the counter for the given key.
func (c *SafeMapStringArray) Pop(key string) []StringTime {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	if len(c.v[key]) > 10000 {
		sret := c.v[key][0:10000]
		c.v[key] = c.v[key][10000:len(c.v[key])]
		return sret
	}else{
		sret := c.v[key]
		c.v[key] = []StringTime{}
		return sret
	}
}

func (c *SafeMapStringArray) ArrayCount(key string) int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	return len(c.v[key])
}

func (c *SafeMapStringArray) Init() {
	c.mu.Lock()
	c.v = make(map[string][]StringTime)
	defer c.mu.Unlock()
}


func GetNowTimeMinute(t int64, timeLoc *time.Location) string {
	tt := time.Unix(t, 0)
	if tt.Year() < 2006 {
		tt = time.Now()
	}
	tm := tt.In(timeLoc).Format("200601021504")
	tmInt, _ := strconv.ParseInt(tm, 10, 64)
	tmInt = int64(tmInt/10) * 10
	return strconv.FormatInt(tmInt, 10)
}

func GetTimeDay(t int64, timeLoc *time.Location) string {
	tt := time.Unix(t, 0)
	if tt.Year() < 2006 {
		tt = time.Now()
	}
	return tt.In(timeLoc).Format("20060102")
}

func GetTimeMonth(t int64, timeLoc *time.Location) string {
	tt := time.Unix(t, 0)
	if tt.Year() < 2006 {
		tt = time.Now()
	}
	return tt.In(timeLoc).Format("200601")
}

func ParseInLocationTime(s string, timeLoc *time.Location) time.Time {
	timeFile, err := time.ParseInLocation("200601021504", s, timeLoc)
	if err != nil {
		timeFile, err = time.ParseInLocation("2006010215", s, timeLoc)
	}
	if err != nil {
		timeFile, err = time.ParseInLocation("20060102", s, timeLoc)
	}
	if err != nil {
		timeFile, err = time.ParseInLocation("2006-01-02 15:04:05", s, timeLoc)
	}
	if err != nil {
		timeFile, err = time.ParseInLocation("2006-01-02 15:04", s, timeLoc)
	}
	if err != nil {
		timeFile, err = time.ParseInLocation("2006-01-02 15", s, timeLoc)
	}
	if err != nil {
		timeFile, err = time.ParseInLocation("2006-01-02", s, timeLoc)
	}
	if err != nil {
		timeFile, err = time.ParseInLocation("2006-01-02T15:04:05.000+08:00", s, timeLoc)
	}
	if err != nil {
		timeFile, err = time.ParseInLocation("2006-01-02T15:04:05Z", s, timeLoc)
	}
	if err != nil {
		timeFile, err = time.ParseInLocation("20060102150405", s, timeLoc)
	}
	if timeFile.Year() > 2001 {
		return timeFile
	}
	return time.Date(1990,0,0,0,0,0,0, timeLoc)
}


func TimeFormat(t int64, loc *time.Location) string  {
	timeInfo := time.Unix(t, 0).In(loc)
	if timeInfo.Year() > 2100 {
		timeInfo = time.Unix(t/1000, 0).In(loc)
	}
	return timeInfo.Format("2006-01-02 15:04:05")
}


func TimeFormat_miute(t int64, loc *time.Location) string  {
	timeInfo := time.Unix(t, 0).In(loc)
	if timeInfo.Year() > 2100 {
		timeInfo = time.Unix(t/1000, 0).In(loc)
	}
	return timeInfo.Format("2006-01-02 15:04:00")
}


func TimeFormat_day(t int64, timeLoc *time.Location) string  {
	timeInfo := time.Unix(t, 0).In(timeLoc)
	if timeInfo.Year() > 2100 {
		timeInfo = time.Unix(t/1000, 0).In(timeLoc)
	}
	return timeInfo.Format("2006-01-02 00:00:00")
}

func TimeFormat_hour(t int64, timeLoc *time.Location) string  {
	timeInfo := time.Unix(t, 0).In(timeLoc)
	if timeInfo.Year() > 2100 {
		timeInfo = time.Unix(t/1000, 0).In(timeLoc)
	}
	return timeInfo.In(timeLoc).Format("2006-01-02 15:00:00")
}


func IP_InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func NarrowInt(x int64) (int) {
	y := int(x)
	if math.Abs(float64(y)) >= math.MaxInt32 {
		y = math.MaxInt32-1
	}
	return y
}

func GetNowTimeCumstMinute(t int64, m int, timeLoc *time.Location) string {
	tt := time.Unix(t, 0)
	if tt.Year() < 2006 {
		tt = time.Now()
	}
	ttUnix := tt.Unix() - tt.Unix()%(int64(m)*60)
	tt = time.Unix(ttUnix, 0)
	tm := tt.In(timeLoc).Format("200601021504")
	return tm
}

func GetNowTimeOneMinute(t int64, timeLoc *time.Location) string {
	tt := time.Unix(t, 0)
	if tt.Year() < 2006 {
		tt = time.Now()
	}
	tm := tt.In(timeLoc).Format("200601021504")
	tmInt, _ := strconv.ParseInt(tm, 10, 64)
	return strconv.FormatInt(tmInt, 10)
}

func FileMd5sum(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "Open err, " + filepath, err
	}
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return "ioutil.ReadAll, " + filepath, err
	}
	md5 := fmt.Sprintf("%x", md5.Sum(body))
	return md5, nil
}


func  CopyFile(dstFileName string, srcFileName string) (written int64, err error) {

	srcFile, err := os.Open(srcFileName)

	if err != nil {
		return 0, err
	}

	defer srcFile.Close()

	//通过srcFile，获取到Reader
	reader := bufio.NewReader(srcFile)

	//打开dstFileName
	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}

	writer := bufio.NewWriter(dstFile)
	defer func() {
		writer.Flush() //把缓冲区的内容写入到文件
		dstFile.Close()
	}()

	return io.Copy(writer, reader)
}


func GetFileLastModifyTime(f string) time.Time {
	finfo, err := os.Stat(f)
	if err != nil {
		return time.Now()
	}
	linuxFileAttr := finfo.Sys().(*syscall.Stat_t)
	valueOfFile := reflect.ValueOf(*linuxFileAttr)
	fileTime := time.Now()
	ftValue := valueOfFile.FieldByName("Atim")
	if !ftValue.IsValid() {
		ftValue = valueOfFile.FieldByName("Atimespec")
	}
	if ftValue.IsValid() && !ftValue.IsZero() {
		ftValueSec := ftValue.FieldByName("Sec")
		ftValueNSec := ftValue.FieldByName("Nsec")
		ftIntSec := int64(0)
		ftIntNSec := int64(0)
		if ftValueSec.IsValid() && !ftValueSec.IsZero() {
			ftIntSec = ftValueSec.Int()
		}
		if ftValueNSec.IsValid() && !ftValueNSec.IsZero() {
			ftIntNSec = ftValueNSec.Int()
		}
		fileTime = time.Unix(ftIntSec, ftIntNSec)
	}
	ftValueModify := valueOfFile.FieldByName("Mtim")
	if !ftValueModify.IsValid() {
		ftValueModify = valueOfFile.FieldByName("Mtimespec")
	}
	if ftValueModify.IsValid() && !ftValueModify.IsZero() {
		ftValueSec := ftValueModify.FieldByName("Sec")
		ftValueNSec := ftValueModify.FieldByName("Nsec")
		ftIntSec := int64(0)
		ftIntNSec := int64(0)
		if ftValueSec.IsValid() && !ftValueSec.IsZero() {
			ftIntSec = ftValueSec.Int()
		}
		if ftValueNSec.IsValid() && !ftValueNSec.IsZero() {
			ftIntNSec = ftValueNSec.Int()
		}
		if fileTime.Unix() < ftIntSec {
			fileTime = time.Unix(ftIntSec, ftIntNSec)
		}
	}
	return fileTime
}


func SubstrByByte(str string, length int) string {
	if len(str) <= length {
		return str
	}

	bs := []byte(str)[:length]
	bl := 0
	for i:=len(bs)-1; i>=0; i-- {
		switch {
		case bs[i] >= 0 && bs[i] <= 127:
			return string(bs[:i+1])
		case bs[i] >= 128 && bs[i] <= 191:
			bl++;
		case bs[i] >= 192 && bs[i] <= 253:
			cl := 0
			switch {
			case bs[i] & 252 == 252:
				cl = 6
			case bs[i] & 248 == 248:
				cl = 5
			case bs[i] & 240 == 240:
				cl = 4
			case bs[i] & 224 == 224:
				cl = 3
			default:
				cl = 2
			}
			if bl+1 == cl {
				return string(bs[:i+cl])
			}
			return string(bs[:i])
		}
	}
	return ""
}


func PareTimeRangeForString(s string, tLocal *time.Location) (start, end string) {
	tn := time.Now()
	nowTimeString := tn.In(tLocal).Format("2006-01-02 15:04:05")

	pre_day_t := tn.In(tLocal).AddDate(0,0,-1)
	pre_day_unix := time.Date(pre_day_t.Year(), pre_day_t.Month(), pre_day_t.Day(), 23,59,59,0, tLocal)
	preDayEnd := pre_day_unix.In(tLocal).Format("2006-01-02 15:04:05")

	switch s {
	case "当天":
		return tn.In(tLocal).Format("2006-01-02 00:00:00"), nowTimeString;
	case "昨天":
		return tn.In(tLocal).AddDate(0, 0, -1).Format("2006-01-02 00:00:00"), preDayEnd;
	}
	pre_val := 999
	if strings.Contains(s, "1") || strings.Contains(s, "一") {
		pre_val = -1
	}
	if strings.Contains(s, "2") || strings.Contains(s, "两") || strings.Contains(s, "二天"){
		pre_val = -2
	}
	if strings.Contains(s, "3") || strings.Contains(s, "三") {
		pre_val = -3
	}
	if strings.Contains(s, "4") || strings.Contains(s, "四") {
		pre_val = -4
	}
	if strings.Contains(s, "5") || strings.Contains(s, "五") {
		pre_val = -5
	}
	if strings.Contains(s, "6") || strings.Contains(s, "六") {
		pre_val = -6
	}
	if strings.Contains(s, "7") || strings.Contains(s, "七") {
		pre_val = -7
	}
	if strings.Contains(s, "10") || strings.Contains(s, "十") {
		pre_val = -10
	}
	if strings.Contains(s, "12") || strings.Contains(s, "十二") {
		pre_val = -12
	}
	if strings.Contains(s, "14") || strings.Contains(s, "十四") {
		pre_val = -14
	}
	if strings.Contains(s, "24") || strings.Contains(s, "二十四") {
		pre_val = -24
	}
	if pre_val != 999 {
		if strings.Contains(s, "天") {
			if strings.Contains(s, "过去") {
				//过去指的是第一天的00点00分，截止到前一天的23点59分
				return tn.In(tLocal).AddDate(0,0, pre_val).Format("2006-01-02 00:00:00"), preDayEnd;
			}else if strings.Contains(s, "最近") {
				//最近指的是从当前时间往前推指定的时间值
				return tn.In(tLocal).AddDate(0,0, pre_val).Format("2006-01-02 15:04:05"), nowTimeString;
			}
		}
		if strings.Contains(s, "小时") {
			td_now := time.Unix(time.Now().Unix()+int64(pre_val*60*60), 0)
			if strings.Contains(s, "过去") {
				//过去指的是第一天的00点00分，截止到前一天的23点59分
				pre_hour := ParseInLocationTime(tn.In(tLocal).Format("2006-01-02 15:00:00"), tLocal).Unix()-1
				return td_now.In(tLocal).Format("2006-01-02 15:00:00"), time.Unix(pre_hour, 0).In(tLocal).Format("2006-01-02 15:04:05");
			}else if strings.Contains(s, "最近") {
				//最近指的是从当前时间往前推指定的时间值
				return td_now.In(tLocal).Format("2006-01-02 15:04:05"), nowTimeString;
			}
		}
	}
	//默认当天
	return tn.In(tLocal).Format("2006-01-02 00:00:00"), nowTimeString;
}

func StringInArray(s string, arr []string) bool {
	for _, k := range arr {
		if k == s {
			return true
		}
	}
	return false
}

//取并集
func GetMapUnionData(src, des map[int64]int64) map[int64]int64 {
	ret_map := des
	for _, s := range src {
		ret_map[s] = s
	}
	return ret_map
}

//取交集
func GetMapInnerData(src, des map[int64]int64) map[int64]int64 {
	ret_map := map[int64]int64{}
	for _, s := range src {
		isExist := false
		for _, d := range des {
			if s == d {
				isExist = true
				break
			}
		}
		if isExist {
			ret_map[s] = s
		}
	}
	return ret_map
}

func ArrayItemDeduplication(src, app []int64) []int64 {
	map_int := map[int64]int64{}
	for _, s := range src {
		map_int[s] = s
	}
	for _, a := range app {
		map_int[a] = a
	}

	ret_arr := []int64{}
	for _, m := range map_int {
		ret_arr = append(ret_arr, m)
	}
	return ret_arr
}
