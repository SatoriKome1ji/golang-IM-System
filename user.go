package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		addr,
		addr,
		make(chan string),
		conn,
		server,
	}
	go user.ListenMessage()
	return user
}

func (this *User) Online() {
	this.server.MapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.MapLock.Unlock()
	this.server._BroadCast(this, "is on line")
}

func (this *User) Offline() {
	close(this.C)
	this.server.MapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.MapLock.Unlock()
	this.server._BroadCast(this, "is offline")
}

func (this *User) _SendMSG(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) _CheckOnlineUsers() {
	this.server.MapLock.Lock()
	for _, user := range this.server.OnlineMap {
		onlineMsg := "[" + user.Name + "] is online\n"
		this._SendMSG(onlineMsg)
	}
	this.server.MapLock.Unlock()
}

func (this *User) _ChangeName(newName string) {
	if this.Name == newName {
		this._SendMSG("please input a new name different to current name")
		return
	}
	this.server.MapLock.Lock()
	_, ok := this.server.OnlineMap[newName]
	this.server.MapLock.Unlock()
	if ok {
		this._SendMSG("\"" + newName + "\" already used")
		return
	}
	this.server.MapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.OnlineMap[newName] = this
	this.server.MapLock.Unlock()
	this.Name = newName
	this._SendMSG("your name has changed to\"" + newName + "\"")
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this._CheckOnlineUsers()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		this._ChangeName(strings.Split(msg, "|")[1])
	} else {
		this.server._BroadCast(this, msg)
	}
}
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
