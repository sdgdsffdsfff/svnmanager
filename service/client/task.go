package client

import (
	"king/service/task"
	"king/helper"
	"king/service/master"
	"fmt"
	"king/service/webSocket"
)

func init(){
	task.New("Heartbeat", func(this *task.Task) interface{} {
		Refresh()
		return nil
	})

	updating := false

	task.New("UpdateHostUnDeployList", func(this *task.Task) interface{} {

		if updating {
			return
		}

		this.Enable = false
		updating = true

		defer func(){
			this.Enable = true
			updating = false
		}()

		if master.UnDeployList != nil {
			helper.AsyncMap(hostMap, func(key, value interface{}) bool {
				c := value.(*HostClient)
				if err := c.UpdateUnDeployList(&master.UnDeployList); err != nil {
					fmt.Println(err)
				} else {
					webSocket.Notify("Update")
				}
				return false
			})
		}
		return nil
	})
}
