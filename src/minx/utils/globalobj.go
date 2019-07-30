package utils

import (
	"../ziface"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GlobalObj struct {
	TcpServer     ziface.IServer //全局Server对象
	Host          string         //IP地址
	TcpPort       int            //IP端口
	Name          string         //服务器名字
	Version       string         //服务器版本
	MaxPacketSize uint32         //单个包的最大字节数
	MaxConn       int            // 当前服务器的最大连接数
}

func (g *GlobalObj) Reload() {
	// 读取配置文件
	file := "conf/zinx.json"
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(dir)
	return
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	// 配置文件转对象
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:          "Zinx Main Server",
		Version:       "V0.4",
		TcpPort:       7777,
		Host:          "127.0.0.1",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}

	//GlobalObject.Reload()
}
