package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"stash-go/config"
	"stash-go/redis"
	"stash-go/server"
	"syscall"
)

func main() {

	fmt.Print(`         __             __  
   _____/ /_____ ______/ /_ 
  / ___/ __/ __ ` + "`" + `/ ___/ __ \
 (__  ) /_/ /_/ (__  ) / / /
/____/\__/\__,_/____/_/ /_/   by: zgr`)
	fmt.Println()

	addr := flag.String("redis_addr", "127.0.0.1:6379", "Redis server address")
	redisUsername := flag.String("redis_username", "", "Redis username")
	redisPassword := flag.String("redis_password", "", "Redis password")
	redisDB := flag.Int("redis_db", 0, "Redis database number")
	httpPort := flag.Int("http_port", 8080, "HTTP server port")

	flag.Parse()

	// config
	config.SetUp()

	// redis
	redis.Init(*addr, *redisUsername, *redisPassword, *redisDB)

	// server
	go server.Run(*httpPort)

	// 等待终端信号来优雅关闭服务器
	quit := make(chan os.Signal, 1)                      // 创建一个接受信号的通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞此处，当接受到上述两种信号时，才继续往下执行
}
