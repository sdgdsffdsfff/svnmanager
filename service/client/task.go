package client

import (
	"fmt"
	"king/helper"
	"king/service/master"
	"king/service/task"
)

func init() {
	task.New("Heartbeat", func(this *task.Task) interface{} {
		Refresh()
		return nil
	})

	updating := false

	task.New("UpdateHostUnDeployList", func(this *task.Task) interface{} {

		if updating {
			return nil
		}

		updating = true

		defer func() {
			this.Enable = false
			updating = false
		}()

		if master.UnDeployList != nil {
			helper.AsyncMap(hostMap, func(key, value interface{}) bool {
				c := value.(*HostClient)
				if err := c.UpdateUnDeployList(&master.UnDeployList); err != nil {
					fmt.Println(err)
				}
				return false
			})
		}
		return nil
	})
}
