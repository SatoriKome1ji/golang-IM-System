package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		ip,
		port,
	}
}

func (this *Server) _Handler(conn net.Conn) {
	fmt.Println("connection sucess, handling...")
}

func (this *Server) Start() {
	listen, listen_err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if listen_err != nil {
		fmt.Println("listen error")
		return
	}
	defer listen.Close()

	for {
		conn, acc_err := listen.Accept()
		if acc_err != nil {
			fmt.Println("accept error")
			continue
		}
		go this._Handler(conn)
	}
}
