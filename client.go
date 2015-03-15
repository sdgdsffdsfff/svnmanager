package main

import (
	"king/bootstrap"
	"king/config"
	_ "king/routes/client"
	_ "king/rpc/client"
	"log"
	"os"
	"king/rpc"
	"king/model"
	"time"
	"github.com/golang/glog"
)

func active(port string) error{
	_, err := rpc.Send(
		config.MasterRpc(),
		"RpcServer.Active",
		model.WebServer{
			Name: "Unknown",
			Ip: "127.0.0.1",
			Port: port,
		},
	)
	if err != nil {
		time.Sleep(time.Second * 2)
		glog.Errorln(err)
		return active(port)
	}
	return nil
}

func main() {
	port := config.GetString("clientPort")
	argLen := len(os.Args)

	if argLen > 1 {
		port = os.Args[1]
	}

	log.Println("Running client on port", port)

	//设置Rpc地址
	config.Set("rpc", config.GetString("master")+"/rpc")


	bootstrap.Start(port, func(){
		//在服务器上添加自己，必须确定唯一属性
		//active(port);
		log.Println("already connect to server")
	})
}
