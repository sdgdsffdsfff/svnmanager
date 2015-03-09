package websocket

import (
	"king/service/webSocket"
	"king/service/client"
	"king/utils/JSON"
	"king/helper"
)

func init(){
	webSocket.BindWebSocketMethod("procstat", func() JSON.Type {
		list := client.GetAliveList()
		result := JSON.Type{}

		for _, c := range list {
			result[helper.Itoa64(c.Id)] = JSON.Parse(c.Proc)
		}

		return helper.Success(result)
	})
}
