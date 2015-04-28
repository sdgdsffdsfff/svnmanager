package webSocket

import (
	"github.com/go-martini/martini"
	"king/helper"
	"king/service/client"
	"king/service/master"
	"king/service/webSocket"
	"king/utils/JSON"
)

func Socket(params martini.Params, receiver <-chan *webSocket.Message, sender chan<- *webSocket.Message, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
	return webSocket.Listen(params, receiver, sender, done, disconnect, err)
}

func init() {

	webSocket.OnAppend(func(clientLength int) {
		client.StartTask()
	})

	webSocket.OnOut(func(clientLength int) {
		if clientLength == 0 {
			client.StopTask()
		}
	})

	webSocket.OnEmit("heartbeat", func() JSON.Type {
		list := client.List()
		result := JSON.Type{}

		for _, c := range list {
			result[helper.Itoa64(c.Id)] = JSON.Type{
				"status":  c.Status,
				"message": c.Message,
				"error":   c.Error,
			}
		}
		return helper.Success(result)
	})

	webSocket.OnEmit("procstat", func() JSON.Type {
		list := client.GetAliveList()
		result := JSON.Type{}

		for _, c := range list {
			result[helper.Itoa64(c.Id)] = JSON.Parse(c.Proc)
		}

		return helper.Success(result)
	})

	webSocket.OnEmit("master", func() JSON.Type {
		return helper.Success(JSON.Type{
			"message": master.Message,
			"error":   master.Error,
			"status":  master.Status,
		})
	})
}
