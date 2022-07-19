package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v", err)
		return
	}
	defer cli.Close()

	// watch
	watchCh := cli.Watch(context.Background(), "xyy")

	for wresp := range watchCh {
		for _, evt := range wresp.Events {
			fmt.Printf("type:%s key:%s value:%s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
		}
	}
}
