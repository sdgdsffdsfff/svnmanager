package rpc

import (
	"bytes"
	"github.com/gorilla/rpc"
	rpcJSON "github.com/gorilla/rpc/json"
	"io/ioutil"
	"king/helper"
	"king/utils/JSON"
	"net"
	"net/http"
	"time"
)

type RpcArgs interface {
	String() string
}

type RpcParams map[string]string

type RpcReply struct {
	Response interface{}
}

var rpcCtrlList []interface{}

func AddCtrl(ctrl interface{}) {
	rpcCtrlList = append(rpcCtrlList, ctrl)
}

func GetServer() *rpc.Server {
	s := rpc.NewServer()
	s.RegisterCodec(rpcJSON.NewCodec(), "application/json")

	for _, v := range rpcCtrlList {
		s.RegisterService(v, "")
	}
	return s
}

func dial(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(5*time.Second))
}

func Send(url string, method string, params interface{}) (interface{}, error) {
	contentString := `{"method": "` + method + `", "params":[` + JSON.Stringify(params) + `], "id":"` + helper.RandString(10) + `"}`
	contentBody := []byte(contentString)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(contentBody))
	defer req.Body.Close()

	if err != nil {
		return nil, helper.NewError("create request error:", err)
	}

	req.Header.Set("Content-Type", "application/json")

	transport := http.Transport{
		Dial: dial,
	}

	client := &http.Client{
		Transport: &transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, helper.NewError("rpc request error", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := JSON.Parse(string(body))

	if result == nil {
		return nil, nil
	}

	if result["error"] != nil {
		return nil, helper.NewError(method + ":" + result["error"].(string))
	}
	//无返回内容
	if result["result"] == nil {
		return nil, nil
	}
	return result["result"].(map[string]interface{})["Response"], nil
}
