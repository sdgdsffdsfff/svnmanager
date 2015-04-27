package client

import (
	"king/helper"
	"king/service/master"
	"king/service/task"
	"king/service/webSocket"
)

func init() {
	task.New("Heartbeat", func(this *task.Task) interface{} {
		Refresh()
		return nil
	})

	updating := false

	task.New("client.UpdateHostUnDeployList", func(this *task.Task) interface{} {

		if updating {
			return nil
		}

		updating = true

		defer func() {
			this.Enable = false
			updating = false
		}()

		webSocket.BroadCastAll(&webSocket.Message{
			Method: "syncDeployList",
		})

		if master.UnDeployList != nil {
			result := map[int64]error{}
			helper.AsyncMap(hostMap, func(key, value interface{}) bool {
				c := value.(*HostClient)
				err := c.UpdateUnDeployList(&master.UnDeployList)
				result[c.Id] = err
				return false
			})
			webSocket.BroadCastAll(&webSocket.Message{
				Method: "syncDeployList",
				Params: result,
			})
		}
		return nil
	})
}
