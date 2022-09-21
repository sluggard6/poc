package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sluggard/poc/application"
	"github.com/sluggard/poc/config"

	log "github.com/sirupsen/logrus"
)

var (
	configFile = flag.String("c", config.DefaultConfigPath, "设置项目配置文件地址`<path>`，可选")
	version    = flag.Bool("v", false, "打印版本号并退出")
	command    string
)

// Version 程序版本号
const Version = "0.0.1-snapshot"

func main() {
	log.SetLevel(log.TraceLevel)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [-options] [command]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  start 启动项目\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		// fmt.Fprintf(os.Stderr, "  -c <path>\n")
		// fmt.Fprintf(os.Stderr, "    设置项目配置文件路径，可选\n")
		// fmt.Fprintf(os.Stderr, "  -v 打印项目版本号，默认为: false\n")
		// fmt.Fprintf(os.Stderr, "    打印版本号\n")
		// fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Parse()

	log.Debugln(os.Args)

	if (os.Args[len(os.Args)-1])[0] != '-' {
		command = os.Args[len(os.Args)-1]
	}

	if *version {
		fmt.Printf("当前版本：%s\n", Version)
		return
	}

	switch command {
	case "start":
		fmt.Println("server starting...")
	default:
		flag.Usage()
		return
	}

	irisServer := application.NewServer(config.LoadConfing(*configFile))
	log.Info(config.GetConfig())
	if irisServer == nil {
		panic("http server 初始化失败")
	}

	// if libs.IsPortInUse(libs.Config.Port) {
	// 	if !irisServer.Status {
	// 		panic(fmt.Sprintf("端口 %d 已被使用\n", libs.Config.Port))
	// 	}
	// 	irisServer.Stop() // 停止
	// }

	err := irisServer.Start()
	if err != nil {
		panic(fmt.Sprintf("http server 启动失败: %+v", err))
	}

}
