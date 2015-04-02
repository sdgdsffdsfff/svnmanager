package config

const NGINX_CONF = "nginx.conf"
const INT_MAX = 65536

var masterIp string
var masterRpc string
var resIp string

func init() {
	masterIp = GetString("master")
	masterRpc = masterIp + "/rpc"
	resIp = GetString("resServer")
}

func ResServer() string {
	return resIp
}

func MasterIp() string {
	return masterIp
}

func MasterRpc() string {
	return masterRpc
}
