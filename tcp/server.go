//package main  //執行要切換
package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (u *User) String() string {
	return u.Addr + ", UID:" + strconv.Itoa(u.ID) + ", Enter At:" +
		u.EnterAt.Format("2006-01-02 15:04:05+8000")
}

// 生成用户 ID
var (
	globalID int
	idLocker sync.Mutex
)

type Message struct {
	OwnerID int
	Content string
}

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++
	return globalID
}

var (
	enteringChannel = make(chan *User)
	leavingChannel  = make(chan *User)
	//messageChannel = make(chan string, 8)
	messageChannel = make(chan Message, 8)
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	go sendMessage(conn, user.MessageChannel)
	user.MessageChannel <- "Welcome" + user.String()
	//messageChannel <- "user:`" + strconv.Itoa(user.ID) +"` has enter"

	msg := Message{
		OwnerID: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has enter",
	}
	messageChannel <- msg

	//記錄到全局用戶列表中，避免用鎖
	enteringChannel <- user

	// 迴圈讀取用戶輸入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChannel <- msg
	}
	if err := input.Err(); err != nil {
		log.Println("讀取錯誤：", err)
	}
	leavingChannel <- user
	msg.Content = "user:`" + strconv.Itoa(user.ID) + "` has left"
	messageChannel <- msg
}

func sendMessage(conn net.Conn, channel <-chan string) {
	for msg := range channel {
		fmt.Fprintln(conn, msg)
	}
}

func broadcaster() {
	log.Println("進入broadcaster")
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			//新用戶進入
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)
		case msg := <-messageChannel:
			//向所有在線用戶發送訊息
			for user := range users {
				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}
	}
}
