package bootstrap

import (
	"flag"
	"github.com/go-martini/martini"
	"github.com/golang/glog"
	"github.com/martini-contrib/render"
	"king/config"
	"king/utils/db"
	_ "king/controller/common"
	"king/rpc"
	"net/http"
)

var methods = []func(){}

func Register(fn func()){
	methods = append( methods, fn)
}

func Start(port string, onStart func()) {

	// Logging init
	flag.Set("log_dir", config.GetString("log_dir"))
	flag.Set("alsologtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		Delims:  render.Delims{"${", "}"},
	}))

	config.MappingController(m)

	http.Handle("/rpc", rpc.GetServer())
	http.Handle("/", m)

	if db.IsConnected() {
		defer db.Close()
	}

	onStart()

	for _, fn := range methods {
		go fn()
	}

	http.ListenAndServe(":"+port, nil)
}
