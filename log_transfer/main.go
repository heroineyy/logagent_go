package main

import (
	"code.oldboy.com/studygolang/day16/log_transfer/es"
	"code.oldboy.com/studygolang/day16/log_transfer/kafka"
	"code.oldboy.com/studygolang/day16/log_transfer/model"
	"fmt"
	"github.com/go-ini/ini"
)

// log transfer
// 从kafka消费日志数据,写入ES

func main(){
	// 1. 加载配置文件
	var cfg = new(model.Config)
	err := ini.MapTo(cfg, "./config/logtransfer.ini")
	if err != nil {
		fmt.Printf("load config failed,err:%v\n", err)
		panic(err)
	}
	fmt.Printf("%#v\n", *cfg)
	fmt.Println("load config success")
	// 2. 连接ES
	err = es.Init(cfg.ESConf.Address, cfg.ESConf.Index, cfg.ESConf.GoNum, cfg.ESConf.MaxSize)
	if err != nil {
		fmt.Printf("Init es failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("Init ES success")
	// 3. 连接kafka
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Topic)
	if err != nil {
		fmt.Printf("connect to kafka failed,err:%v\n", err)
		panic(err)
	}
	fmt.Println("Init kafka success")
	// 在这儿停顿!
	select {}
}
