package tailfile

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/kafka"
	"strings"
	"time"
)

//tail相关
type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	return tt
}

func (t *tailTask) Init() (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tObj, err = tail.TailFile(t.path, cfg)
	return
}

func (t *tailTask) run() {
	logrus.Info("path:%s is running...", t.path)
	//读取数据，发往kafka
	//logfile-->TailObj --> log --> Client -->kafka
	//循环读数据
	for {
		select {
		case <-t.ctx.Done():
			logrus.Infof("path:%s is stopping...", t.path)
			return
		case line, ok := <-t.tObj.Lines: //chan tail.Line
			if !ok {
				logrus.Warn("tail file close reopen,path:%s\n", t.path)
				time.Sleep(time.Second) //读取出错等一秒
				continue
			}
			//如果是空行就略过
			if len(strings.Trim(line.Text, "\r")) == 0 {
				logrus.Info("出现空行啦，直接跳过...")
				continue
			}
			//利用通道将同步的日志代码改为异步的
			//把读出来的一行日志包装给kafka里面的msg类型
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			//丢到通道中
			//kafka.MsgChan <- msg
			kafka.ToMsgChan(msg)
		}

	}
}
