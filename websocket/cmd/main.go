package main

import (
	"chatroom/websocket/server"
	"fmt"
	"log"
	"net/http"
)

var (
	addr   = ":2002"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

       Go语言编程之旅 —— 一起用Go做项目：ChatRoom，start on：%s
    `
)

func main() {
	fmt.Printf(banner+"\n", addr)
	server.RegisterHandle()
	log.Fatal(http.ListenAndServe(addr, nil))

}
