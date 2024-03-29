package utils

import (
	"../ziface"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// GlobalObj 配置结构
type GlobalObj struct {
	// server
	TCPServer ziface.IServer //全局Server对象
	Host      string         //IP地址
	TCPPort   int            //IP端口
	// zinx
	Name           string //服务器名字
	Version        string //服务器版本
	MaxPacketSize  uint32 //单个包的最大字节数
	MaxConn        int    //当前服务器的最大连接数
	WorkPoolSize   uint32 //工作线程池的数量
	MaxWorkTaskLen uint32 //最大的任务存量
	MaxMsgChanLen  uint32 //最多的发送消息的缓冲长度
	// other
	// TODO 配置
}

// Reload 载入配置文件
func (g *GlobalObj) Reload() {
	root := os.Getenv("ZINX_ROOT_PATH")
	if len(strings.TrimSpace(root)) == 0 {
		panic("error root path")
	}
	file := "conf/zinx.json"
	data, err := ioutil.ReadFile(root + "/" + file)
	if err != nil {
		panic(err)
	}

	// 配置文件转对象
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// GlobalObject 全局配置文件
var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:           "Zinx Main Server",
		Version:        "V0.4",
		TCPPort:        7777,
		Host:           "127.0.0.1",
		MaxConn:        12000,
		MaxPacketSize:  4096,
		WorkPoolSize:   10,
		MaxWorkTaskLen: 1024,
		MaxMsgChanLen:  8192,
	}

	//GlobalObject.Reload()
}
