package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"configLoad/g"
)

func main() {
	confFile := flag.String("c", "setting.yml", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println("0.0.1")
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.InitConfig(*confFile)
	// 初始化web接口
	g.WebAPIStart()
	fmt.Println("初始化完成")
	go fmt.Println("do someting")
	fmt.Println("任务启动")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	select {}
}
