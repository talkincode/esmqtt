package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	_ "time/tzdata"

	_ "github.com/joho/godotenv/autoload"
	"github.com/talkincode/esmqtt/app"
	"github.com/talkincode/esmqtt/assets"
	"github.com/talkincode/esmqtt/common/zaplog/log"
	"github.com/talkincode/esmqtt/config"
	"github.com/talkincode/esmqtt/installer"
	"github.com/talkincode/esmqtt/mqttc"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

// 命令行定义
var (
	h         = flag.Bool("h", false, "help usage")
	showVer   = flag.Bool("v", false, "show version")
	conffile  = flag.String("c", "", "config yaml file")
	install   = flag.Bool("install", false, "run install")
	uninstall = flag.Bool("uninstall", false, "run uninstall")
)

// PrintVersion Print version information
func PrintVersion() {
	fmt.Printf(assets.BuildInfo)
}

func printHelp() {
	if *h {
		ustr := fmt.Sprintf("version: %s, Usage:%s -h\nOptions:", assets.BuildVersion(), os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, ustr)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *showVer {
		PrintVersion()
		os.Exit(0)
	}

	printHelp()

	_config := config.LoadConfig(*conffile)

	// Install as a system service
	if *install {
		err := installer.Install()
		if err != nil {
			log.Error(err)
		}
		return
	}

	if *uninstall {
		installer.Uninstall()
		return
	}

	app.InitGlobalApplication(_config)
	defer app.Release()

	running := make(chan bool)

	err := mqttc.Start()
	if err != nil {
		panic(err)
	}

	<-running

}
