package service

import(
	"time"
	"sync"
	"github.com/go-martini/martini"
	"king/helper"
	"king/utils/JSON"
)

var callbacks = map[string]func()JSON.Type{}
func BindWebSocketMethod(name string, callback func()JSON.Type){
	if _, found := callbacks[name]; found {
		return
	}
	callbacks[name] = callback
}

type Message struct {
	Method string `json:"method"`
	Params interface {} `json:"data"`
}

type socketClient struct {
	in         <-chan *Message
	out        chan<- *Message
	done       <-chan bool
	err        <-chan error
	disconnect chan<- int
	message *Message // 为无限通知预留的调用状态，为最后一次请求内容
}

type webSocket struct {
	sync.Mutex
	clients []*socketClient
	loopEnable bool
}

var WebSocket = &webSocket{sync.Mutex{}, make([]*socketClient, 0), true}

func (r *webSocket) AppendClient(client *socketClient) {
	r.Lock()
	r.clients = append(r.clients, client)
	r.BroadCast(client, &Message{"online", client})
	r.Unlock()
}

func (r *webSocket) removeClient(client *socketClient) {
	r.Lock()
	defer r.Unlock()

	for index, xc := range r.clients {
		if xc == client {
			r.clients = append(r.clients[:index], r.clients[(index+1):]...)
		}
	}
}

func (r *webSocket) Emit(client *socketClient, msg *Message) {
	r.Lock()
	method := msg.Method

	if method == "broadcast" {
		r.BroadCast(client, msg)
	} else if _, found := callbacks[method]; method != "" && found {
		client.out <- &Message{method, callbacks[method]()}
	} else {
		client.out <- &Message{method, helper.Error("method undefined")}
	}
	r.Unlock()
}

func (r *webSocket) NotifyAll(msg *Message){
	for _, xc := range r.clients {
		xc.out <- msg
	}
}

func (r *webSocket) BroadCast(client *socketClient, msg *Message){
	for _, xc := range r.clients {
		if xc != client {
			xc.out <- msg
		}
	}
}

/*
主动推送方法，将会根据method不停向客户端推送内容
 */
func (r *webSocket) loopPushFrame(){
	for {
		r.Lock()
		helper.AsyncMap(r.clients, func(i int) bool{
			client := r.clients[i]
			if client != nil && client.message != nil {
				if method := client.message.Method; method != "" {
					client.out <- &Message{method, callbacks[method]()}
				}
			}
			return false
		})
		r.Unlock()
		time.Sleep(3*time.Second)
	}
}

func (c *webSocket) Listen( params martini.Params, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error ) (int, string){
	client := &socketClient{receiver, sender, done, err, disconnect, nil}
	c.AppendClient(client)
	// 暂停使用无限通知，让客户端自己通过循环调用
	//	if chat.loopEnable {
	//		go chat.loopPushFrame()
	//		chat.loopEnable = false
	//	}
	for {
		select {
		case <-client.err:
			// Don't try to do this:
			// client.out <- &Message{"system", "system", "There has been an error with your connection"}
			// The socket connection is already long gone.
			// Use the error for statistics etc
		case msg := <-client.in:
			c.Emit(client, msg)
		case <-client.done:
			c.removeClient(client)
			return 200, "OK"
		}
	}

}
