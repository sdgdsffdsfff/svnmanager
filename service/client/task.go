package client

import (
	"king/service/task"
)

func init(){
	task.New("Heartbeat", func(this *task.Task) interface{} {
		Refresh()
		return nil
	})
}
