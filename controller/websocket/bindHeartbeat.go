package websocket

import (
	"king/service/webSocket"
	"king/service/client"
	"king/utils/JSON"
	"king/helper"
)

func init(){
	webSocket.BindWebSocketMethod("heartbeat", func() JSON.Type {
		list := client.List()
		result := JSON.Type{}

		for _, c := range list {
			result[helper.Itoa64(c.Id)] = c.Status
		}

		return helper.Success(result)
	})
}
