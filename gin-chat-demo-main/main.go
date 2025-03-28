package main

import (
	"chat/conf"
	"chat/router"
	"chat/service"
	"log"
)

func main() {
	conf.Init()
	go service.Manager.Start()
	go service.StartGroupChatService()
	r := router.NewRouter()
	err := r.Run(conf.HttpPort)
	if err != nil {
		// 不调用err.error()方法，减少逃逸到堆的情况
		log.Fatal(err)
	}
}
