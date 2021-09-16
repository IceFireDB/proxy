// Copyright 2014 Wandoujia Inc. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package main

import (
	"github.com/IceFireDB/proxy/pkg/proxy/router"
	"net/http"
	_ "net/http/pprof"
	"path"

	"github.com/docopt/docopt-go"
	"github.com/ledisdb/xcodis/utils"
	log "github.com/ngaut/logging"
	_ "go.uber.org/automaxprocs"
)

// build info
var (
	BuildVersion = "unknown"
	BuildDate    = "unknown"
)

var (
	addr       = ":9000"
	httpAddr   = ":9001"
	configFile = "config.ini"
)

var usage = `usage: proxy [-c <config_file>] [-L <log_file>] [--log-level=<loglevel>] [--cpu=<cpu_num>] [--addr=<proxy_listen_addr>] [--http-addr=<debug_http_server_addr>]

options:
   -c	set config file
   -L	set output log file, default is stdout
   --log-level=<loglevel>	set log level: info, warn, error, debug [default: info]
   --addr=<proxy_listen_addr>		proxy listen address, example: 0.0.0.0:9000
   --http-addr=<debug_http_server_addr>		debug vars http server
`

func handleSetLogLevel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	level := r.Form.Get("level")
	log.SetLevelByString(level)
	log.Info("set log level to", level)
}

func main() {
	log.SetLevelByString("info")

	args, err := docopt.Parse(usage, nil, true, "codis proxy v0.1", true)
	if err != nil {
		log.Error(err)
	}

	// set config file
	if args["-c"] != nil {
		configFile = args["-c"].(string)
	}

	// set output log file
	if args["-L"] != nil {
		log.SetOutputByName(args["-L"].(string))
	}

	// set log level
	if args["--log-level"] != nil {
		log.SetLevelByString(args["--log-level"].(string))
	}

	// set addr
	if args["--addr"] != nil {
		addr = args["--addr"].(string)
	}

	// set addr type
	if args["--type"] != nil {
		addr = args["--type"].(string)
	}

	// set http addr
	if args["--http-addr"] != nil {
		httpAddr = args["--http-addr"].(string)
	}

	dumppath := utils.GetExecutorPath()

	log.Info("dump file path:", dumppath)
	log.CrashLog(path.Join(dumppath, "codis-proxy.dump"))

	router.CheckUlimit(1024)

	http.HandleFunc("/setloglevel", handleSetLogLevel)
	go http.ListenAndServe(httpAddr, nil)
	log.Info("running on ", addr)
	conf, err := router.LoadConf(configFile)
	if err != nil {
		log.Fatal(err)
	}
	s := router.NewServer(addr, httpAddr, conf)
	s.Run()
	log.Warning("exit")
}
