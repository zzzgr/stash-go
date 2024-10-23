package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"stash-go/config"
	activityController "stash-go/controller/activity"
	pkgController "stash-go/controller/pkg"
	"stash-go/repository"
	"stash-go/server"
	"stash-go/service/activity"
	pkg "stash-go/service/pkg"
	"syscall"
)

func main() {

	fmt.Print(`         __             __  
   _____/ /_____ ______/ /_ 
  / ___/ __/ __ ` + "`" + `/ ___/ __ \
 (__  ) /_/ /_/ (__  ) / / /
/____/\__/\__,_/____/_/ /_/   by: zgr`)
	fmt.Println()

	httpPort := flag.Int("http_port", 9090, "HTTP server port")

	flag.Parse()

	// config
	config.SetUp()

	// dao
	repo := repository.NewRepository(repository.NewDB())

	// service
	activityService := activity.NewService(repo)
	pkgService := pkg.NewService(repo, activityService)

	// controller
	activityController := activityController.New(activityService, pkgService)
	pkgController := pkgController.New(activityService, pkgService)

	// server
	s := server.New(activityController, pkgController)
	go s.Run(*httpPort)

	// 测试
	//activityService.Save(&activity_dto.SaveRequestDTO{
	//	Id:           1,
	//	Name:         "京东测试 newBabelAwardCollection",
	//	Cron:         "40 * * * * ?",
	//	Code:         "jd",
	//	UrlPattern:   "https://api.m.jd.com/client.action",
	//	QueryPattern: "functionId=newBabelAwardCollection",
	//	Type:         1,
	//	Field:        "pt_pin,pin",
	//	Advance:      200,
	//	Interval:     10,
	//})

	//activityService.Save(&entity.Controller{
	//	Name:         "美团外卖refId测试",
	//	Cron:         "40 * * * * ?",
	//	Code:         "mt_waimai",
	//	UrlPattern:   "*promotion.waimai.meituan.com*",
	//	QueryPattern: "couponReferId=4EEA7E4B166946D9B83A57361E01687E",
	//	Type:         1,
	//	Field:        "token",
	//	Advance:      160,
	//	Interval:     20,
	//})

	// 等待终端信号来优雅关闭服务器
	quit := make(chan os.Signal, 1)                      // 创建一个接受信号的通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞此处，当接受到上述两种信号时，才继续往下执行
}
