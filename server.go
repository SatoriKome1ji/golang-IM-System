package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	MapLock   sync.RWMutex

	//消息广播channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (this *Server) _BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Name + "]:" + msg
	this.Message <- sendMsg
}

func (this *Server) _ListenMsg() {
	for {
		msg := <-this.Message
		this.MapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.MapLock.Unlock()
	}
}

func (this *Server) _Handler(conn net.Conn) {
	user := NewUser(conn, this)
	isLive := make(chan bool)

	user.Online()

	go func() {

		for {
			buf := make([]byte, 4096)
			n, connReadErr := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			} else if connReadErr != nil && connReadErr != io.EOF {
				fmt.Println("connection read error")
				return
			}
			user.DoMessage(string(buf[:n-1]))
			isLive <- true
		}
	}()

	select {
	case <-isLive:
	case <-time.After(time.Second * 10):
		user._SendMSG("long time no operation,you should be offline")
		conn.Close()
		return
	}
}

func (this *Server) Start() {
	listen, listenErr := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if listenErr != nil {
		fmt.Println("listen error")
		return
	}
	defer listen.Close()
	go this._ListenMsg()
	for {
		conn, accErr := listen.Accept()
		if accErr != nil {
			fmt.Println("accept error")
			continue
		}
		go this._Handler(conn)
	}
}
