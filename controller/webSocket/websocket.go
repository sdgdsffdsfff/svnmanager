package webSocket

import(
	"github.com/go-martini/martini"
	"king/service/webSocket"
	"king/utils/JSON"
	"king/service/client"
	"king/helper"
	"king/service/group"
)

func Socket( params martini.Params, receiver <-chan *webSocket.Message, sender chan<- *webSocket.Message, done <-chan bool, disconnect chan<- int, err <-chan error ) (int, string) {
	return webSocket.Listen(params, receiver, sender, done, disconnect, err)
}

func init(){
	webSocket.OnEmit("heartbeat", func() JSON.Type {
		list := client.List()
		result := JSON.Type{}

		for _, c := range list {
			result[helper.Itoa64(c.Id)] = c.Status
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
	});

	webSocket.OnEmit("getClientList", func() JSON.Type{
		client.SetHeartEnable(true)
		client.SetProcMonitorEnable(true)

		result := map[string]JSON.Type{}
		clientsDict := map[string]JSON.Type{}

		groups := group.List()
		clients := client.List()

		for id, g := range groups {
			result[helper.Itoa64(id)] = JSON.Parse(g)
			clientsDict[helper.Itoa64(id)] = JSON.Type{}
		}

		for id, c := range clients {
			if list, found := clientsDict[helper.Itoa64(c.Group)]; found {
				list[helper.Itoa64(id)] = c
			}
		}

		for id, g := range result {
			g["Clients"] = clientsDict[id]
		}

		return helper.Success(result)
	})
}