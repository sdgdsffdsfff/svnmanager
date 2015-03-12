package webSocket

import(
	"time"
	"sync"
	"github.com/go-martini/martini"
	"king/helper"
	"king/utils/JSON"
)

var onEmitCallback = map[string]func()JSON.Type{}
func OnEmit(name string, callback func()JSON.Type){
	if _, found := onEmitCallback[name]; found {
		return
	}
	onEmitCallback[name] = callback
}

var onOutCallback = []func(int){}
func OnOut(callback func(int)) {
	onOutCallback = append( onOutCallback, callback)
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

var syncLock = sync.Mutex{}
var clients []*socketClient

func AppendClient(client *socketClient) {
	syncLock.Lock()
	clients = append(clients, client)
	BroadCast(client, &Message{"online", client})
	syncLock.Unlock()
}

func removeClient(client *socketClient) {
	syncLock.Lock()
	defer syncLock.Unlock()

	for index, xc := range clients {
		if xc == client {
			clients = append(clients[:index], clients[(index+1):]...)
		}
	}

	clientLength := len( clients )

	for _, callback := range onOutCallback {
		callback(clientLength)
	}
}

func Emit(client *socketClient, msg *Message) {
	syncLock.Lock()
	method := msg.Method

	if method == "broadcast" {
		BroadCast(client, msg)
	} else if _, found := onEmitCallback[method]; method != "" && found {
		client.out <- &Message{method, onEmitCallback[method]()}
	} else {
		client.out <- &Message{method, helper.Error("method undefined")}
	}
	syncLock.Unlock()
}

//服务器调用客户端方法
func BroadCastAll(msg *Message){
	for _, xc := range clients {
		xc.out <- msg
	}
}

//客户端调用其他客户端方法
func BroadCast(client *socketClient, msg *Message){
	for _, xc := range clients {
		if xc != client {
			xc.out <- msg
		}
	}
}

func Notify(text string) {
	BroadCastAll(&Message{"notify", text})
}

func NotifyOther(client *socketClient,text string) {
	Emit(client, &Message{
		"notify", text,
	})
}

/*
主动推送方法，将会根据method不停向客户端推送内容
 */
func loopPushFrame(){
	for {
		syncLock.Lock()
		helper.AsyncMap(clients, func(key, value interface{}) bool{
			client := value.(*socketClient)
			if client != nil && client.message != nil {
				if method := client.message.Method; method != "" {
					client.out <- &Message{method, onEmitCallback[method]()}
				}
			}
			return false
		})
		syncLock.Unlock()
		time.Sleep(3*time.Second)
	}
}

func Listen( params martini.Params, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error ) (int, string){
	client := &socketClient{receiver, sender, done, err, disconnect, nil}
	AppendClient(client)
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
			Emit(client, msg)
		case <-client.done:
			removeClient(client)
			return 200, "OK"
		}
	}

}
