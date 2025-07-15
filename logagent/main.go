package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
)

//现在的技能包
//往kafka发数据
//使用tail读日志文件

type Config struct {
	KafkaConfig `ini:"kafka"`
	EtcdConfig  `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

func run() {
	select {}
}

func main() {
	var err error
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

	//3. 从etcd中拉取要收集日志的配置项
	collectKey := configObj.EtcdConfig.CollectKey
	allConf, err := etcd.GetConf(collectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd faild,err:%v", err)
		return
	}
	fmt.Println(allConf)

	//4、派一个小弟去监控etcd中 configObj.EtsdConfig.collectKey 对应值得变化
	go etcd.WatchConf(collectKey)

	//5.根据配置中的日志路径使用tail去收集日志
	err = tailfile.Init(allConf) //把从etcd中获取的日志项读取到tailfile中
	if err != nil {
		logrus.Errorf("init tailfile failed,err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")

	run()

}
