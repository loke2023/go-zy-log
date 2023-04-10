package db

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
	"sync"
	"time"
	"go-zy-log/common"
)

type TBLock struct {
	sync.RWMutex
}


type Logmessage struct {
	Type  	string `json:"type"`
	Time  	int64  `json:"time"`
	Level 	string `json:"level"`
	Key 	string `json:"key"`
	Info  	string `json:"info"`
}

type DBControl struct {
	lg 				*logrus.Logger
	timeLoc     	*time.Location
	map_db 			map[string]*mongo.Database
	map_collection 	map[string]*mongo.Collection
	mondb_con		string
}

func ConnectToDB(uri, name string, timeout time.Duration) (*mongo.Database, error)  {

	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// 通过传进来的uri连接相关的配置
	o := options.Client().ApplyURI(uri)

	// 发起链接
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}
	// 判断服务是不是可用
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}

	// 返回 client
	return client.Database(name), nil
}

func NewDBCon(mondb_con string, g *logrus.Logger) DBControl {
	dbCon := DBControl{}
	dbCon.mondb_con = mondb_con
	dbCon.timeLoc, _ = time.LoadLocation("Asia/Shanghai")
	dbCon.lg = g
	dbCon.map_db = make(map[string]*mongo.Database)
	dbCon.map_collection = make(map[string]*mongo.Collection)
	return dbCon
}

func (this *DBControl)PushLog(s []byte) {
	var data Logmessage
	err := json.Unmarshal(s, &data)
	if err != nil {
		this.lg.Error("json unmarshal faild, ", err.Error(), ", log:", string(s))
		return
	}
	this.AddOne(&data)								
}

func (this *DBControl)AddOne(data *Logmessage)  {
	this.CleanHistoryLogs()
	table := data.Level + "_" + time.Now().In(this.timeLoc).Format("2006-01-02")
	database_key := data.Type + "_" + table
	collection, exist := this.map_collection[database_key]
	if !exist {
		//创建一个新的连接
		db_ex, ex := this.map_db[data.Type]
		if !ex {
			maxTime := time.Duration(10) // 链接超时时间
			db, err := ConnectToDB(this.mondb_con, data.Type, maxTime)
			if err != nil {
				this.lg.Error("connect mongodb faild, ", err.Error(), ", log:", *data)
				return
			}
			this.map_db[data.Type] = db
			db_ex = db
		}
		col := db_ex.Collection(table)
		this.map_collection[database_key] = col
		collection = col
	}

	log_data := struct {
		Time 	string
		Info 	string
		Key 	string
	}{
		common.TimeFormat(data.Time, this.timeLoc),
		data.Info,
		data.Key,
	}
	_, err := collection.InsertOne(context.TODO(), log_data)
	if err != nil {
		this.lg.Error(err, ", logdata:", log_data)
		return
	}
}


func (this *DBControl) CleanHistoryLogs() {

	for _, vd := range this.map_db {
		tables, _ := vd.ListCollectionNames(context.Background(), bson.D{})
		for _, tb := range tables {
			del := false
			tb_array := strings.Split(tb, "_")
			if len(tb_array) < 2 {
				continue
			}
			level := tb_array[0]
			tb_time := common.ParseInLocationTime(tb_array[1],this.timeLoc)
			if strings.Contains(strings.ToLower(level), "error") ||
				strings.Contains(strings.ToLower(level), "warning") {
				//错误日志，删除3个月以前的
				if time.Now().Unix() - tb_time.Unix() > 60*60*24*30*3 {
					del = true
				}
			}else{
				//其他日志，删除一个月以前的
				if time.Now().Unix() - tb_time.Unix() > 60*60*24*30 {
					del = true
				}
			}

			if del {
				//删除
				col := vd.Collection(tb)
				col.Drop(context.Background())
				this.lg.Info("删除库", vd.Name(), ", table:", tb)
			}
		}
	}
}