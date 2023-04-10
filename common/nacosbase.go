package common

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)


func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}


func getRemoteUrlConfigInfo(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte("get error"), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("read all error"), err
	}
	return body, nil
}


// config should be a pointer to structure, if not, panic
func  loadCacheConfig(cf string) (string, error) {
	f, err := ioutil.ReadFile(cf)
	if err != nil {
		return "", fmt.Errorf("Can not load config at %s. Error: %v", cf, err)
	}
	return string(f), nil
}


func  writeConfigFile(info string, cfile string) (written int64, err error) {

	//打开dstFileName
	dstFile, err := os.OpenFile(cfile, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		return 0, fmt.Errorf("open file err = %v\n", err)
	}

	writer := bufio.NewWriter(dstFile)
	defer func() {
		writer.Flush() //把缓冲区的内容写入到文件
		dstFile.Close()

	}()

	return io.Copy(writer, strings.NewReader(info))
}


func LoadRemoteNacosConfigInfo(addr, id, groupname, dataid string, loadcache bool) (string, error) {

	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	workPath = workPath + "/transform_cache/nacos/"
 	os.MkdirAll(workPath, os.ModePerm);
	cfile := workPath + getMd5String(id+groupname+dataid) + ".json"

 	//读取key
	Key_public_url := addr + "&group=" + groupname + "&tenant=" + id + "&dataId=" + dataid
	info, err := getRemoteUrlConfigInfo(Key_public_url)
	if err != nil {
		if loadcache {
			cacheInfo, errred := loadCacheConfig(cfile)
			if errred != nil {
				return cacheInfo, fmt.Errorf("remote read err = ", err.Error(),
					", errred:", errred.Error(), ", load cache faild")
			}
			return cacheInfo, fmt.Errorf("remote read err = ", err.Error(), ", load cache info")
		}else{
			return "", fmt.Errorf("remote read err = ", err.Error(), ", no load cache")
		}
	}else{
		//保存info到缓存
		wl, err := writeConfigFile(string(info), cfile)
		if err != nil {
			return string(info), fmt.Errorf("writeConfigFile err = ", err.Error(), ", len:", wl)
		}
	}

	return string(info), nil
}

type NacosConfigChanged func(id, groupname, dataid, newinfo string)

func LoadRemoteNacosForSDK(addr, id, groupname, dataid string, f NacosConfigChanged) (string, error) {
	cc := constant.ClientConfig{
		NamespaceId:         id, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "error",
	}

	sc := []constant.ServerConfig{{IpAddr: addr, Port: 8848,},}

	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	configInfo, err := client.GetConfig(vo.ConfigParam{DataId: dataid, Group: groupname,})
	if err != nil {
		return "", fmt.Errorf("load config err [add:" + addr + ", id:" + id + "] , err:" + err.Error())
	}

	if f != nil {
		//Listen config change,key=dataId+group+namespaceId.
		client.ListenConfig(vo.ConfigParam{
			DataId: dataid,
			Group:  groupname,
			OnChange: f,
		})
	}

	return configInfo, nil
}