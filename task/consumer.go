package task

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/yufeifly/migrator/client"
	"github.com/yufeifly/migrator/model"
	"github.com/yufeifly/migrator/redis"
	"time"
)

var DefaultConsumer *Consumer

type Consumer struct {
	ServicePort string
}

func NewConsumer() *Consumer {
	return &Consumer{
		ServicePort: "6379",
	}
}

// Consume consume a log in task queue
func (c *Consumer) Consume() error {
	cli := client.NewClient()
	// infinity loop, consume logs
	for {
		//fmt.Println("queue: ", DefaultQueue)
		logrus.Info("tick")
		taskJson := DefaultQueue.PopFront()
		if taskJson == "" {
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		// unmarshall get serialized kv
		var task model.Log
		err := json.Unmarshal([]byte(taskJson), &task)
		if err != nil {
			return err
		}
		if len(task.LogQueue) > 0 {
			for _, kv := range task.GetLogQueue() {
				var sli []string
				json.Unmarshal([]byte(kv), &sli)
				redis.Set("service1", sli[0], sli[1])
			}
		}

		// stop this goroutine if it is the last task
		if task.GetLastFlag() {
			logrus.Warn("the last log consumed")
			return nil
		}
		// consumed a log, send this message to src
		logrus.Infof("consumed a log, msg send to src")
		//time.Sleep(200 * time.Millisecond)
		err = cli.ConsumeAdder()
		if err != nil {
			logrus.Errorf("cli.consumed failed, err: %v", err)
			return err
		}
	}
	return nil
}
