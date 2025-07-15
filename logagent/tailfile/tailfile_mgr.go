package tailfile

import (
	"github.com/sirupsen/logrus"
	"logagent/common"
)

//tailTask 的管理者

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask
	collectEntryList []common.CollectEntry
	confChan         chan []common.CollectEntry
}

var (
	ttMgr *tailTaskMgr
)

//main 函数调用
func Init(allConf []common.CollectEntry) (err error) {
	//allConf 里面存了若干个日志的收集项
	//针对每一个日志收集项创建了一个对应的tailObj

	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic) //创建一个日志收集任务
		err = tt.Init()                          //去打开日志文件准备读
		if err != nil {
			logrus.Errorf("create tailobj for:%s failed,err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt //把创建的这个tailTask任务记录在册，方便后续管理
		//起一个后台的goroutine去收集日志
		go tt.run()
	}
	go ttMgr.watch() //在后台等新的配置来
	return
}

func (t *tailTaskMgr) watch() {
	for {
		//派一个小弟等着新配置来
		newConf := <-t.confChan //取到值说明新的配置来啦
		//新配置来了之后应该管理一下我之前启动的哪些tailTask
		logrus.Infof("get new conf from etcd,conf :%v，start mange tailTask", newConf)

		for _, conf := range newConf {
			//1. 原来已经存在的任务就不用动
			if t.isExist(conf) {
				continue
			}
			//2. 原来没有的，我要新创建一个tailTask任务
			tt := newTailTask(conf.Path, conf.Topic) //创建一个日志收集任务
			err := tt.Init()                         //去打开日志文件准备读
			if err != nil {
				logrus.Errorf("create tailobj for:%s failed,err:%v", conf.Path, err)
				continue
			}
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			t.tailTaskMap[tt.path] = tt //把创建的这个tailTask任务记录在册，方便后期起一个后台的goroutine去收集日志
			go tt.run()
		}
		//3. 原来有的现在没有的要tailTask停掉
		//找出tailTaskMap中存在，但是newConf不存在的那些tailTask,把他们都停掉
		for key, task := range t.tailTaskMap {
			var found bool
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				//这个tailTask要停掉了
				logrus.Infof("the task collect path :%s need to stop.", task.path)
				task.cancel()
			}
		}
	}
}

//判断tailTaskMap是否存在该收集项
func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
	return
}
