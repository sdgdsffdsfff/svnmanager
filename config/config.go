package config

import (
	"github.com/dlintw/goconf"
	"king/utils/JSON"
	"king/utils"
)

const (
	Server = iota
	Client
)

var configFile *goconf.ConfigFile

func init(){
	configFile, _ = goconf.ReadConfigFile(utils.GetRuntimeDir("config.conf"))
}

func GetString(key string) string {
	str, _ := configFile.GetString("default", key)
	return str
}

func GetInt(key string) int {
	i, _ := configFile.GetInt("default", key)
	return i
}

var options JSON.Type = make(JSON.Type)

func Set(key string, value interface{}) {
	options[key] = value
}

func Get(key string) interface{} {
	if value, found := options[key]; found {
		return value
	}
	return nil
}

func Env() interface{} {
	if value, found := options["env"]; found {
		return value
	}
	return -1
}

