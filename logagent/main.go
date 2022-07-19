package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
)

// Config 收集指定目录下的日志文件，发送到kafka中
//现在的技能包
//往kafka发数据
//使用tail读日志文件

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

func run() {
	select {}
}

func main() {
	//-1：获取本机IP,为后续去etcd取配置文件打下基础
	ip, err := common.GetOutboundIP()
	if err != nil {
		logrus.Errorf("get ip failed,err:%v", err)
		return
	}
	var configObj = new(Config)
	//0:配置文件`go-ini`
	err = ini.MapTo(configObj, "./config/config.ini")
	if err != nil {
		logrus.Errorf("load config failed,err:%v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)

	//1.初始化连接kafka（做好准备工作）
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed,err:%v", err)
		return
	}
	logrus.Info("init kafka success!")

	//2.初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd faild,err:%v", err)
		return
	}

	//从etcd中拉取要收集日志的配置项
	collectkey := fmt.Sprintf(configObj.EtcdConfig.CollectKey, ip)
	//fmt.Println("123hhhhhh:", collectkey)
	allConf, err := etcd.GetConf(collectkey)
	if err != nil {
		logrus.Errorf("get conf from etcd faild,err:%v", err)
		return
	}
	fmt.Println(allConf)

	//派一个小弟去监控etcd中 configObj.EtsdConfig.Collectkey 对应值得变化
	go etcd.WatchConf(collectkey)

	//2.根据配置中的日志路径使用tail去收集日志
	err = tailfile.Init(allConf) //把从etcd中获取的日志项读取到tailfile中
	if err != nil {
		logrus.Errorf("init tailfile failed,err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")

	run()

}
