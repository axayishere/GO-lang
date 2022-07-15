package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"sync"
)

var (
	port = flag.String("host", "localhost:19980", "Fork Ip address")
	host = flag.String("manager", "localhost:19982", "Address of the manager process")
	fid  = flag.Int("id", 1, "Fork ID")
)

type fork struct { // registratioin fork struction
	safe   sync.Mutex
	id     int
	pid    int
	host   string
	status string
	conn   net.Conn
}
type Response struct {  // registratioin responce structure
	Id     int
	Status bool
}
type Register struct { // registratioin structure
	Addrs string
	Id    int
	Type  int
}
type Request struct { // registratioin request structure
	Id   int
	PId  int
	Type int
}

func NewFork(id int, host string, c net.Conn) *fork { 
	f := &fork{
		status: "free",
		host:   host,
		id:     id,
		conn:   c,
	}
	return f
}
func readMessage(c net.Conn) []byte {
	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		//log.Println(err)
		return nil
	}
	return buffer[:n]
}
func (f *fork) handleNewConnection(c net.Conn) {
	for {
		buf := readMessage(c)
		if buf != nil {
			var req Request
			if err := json.Unmarshal(buf, &req); err != nil {
				panic(err)
			}
			f.handle(req.Type, req.Id, req.PId, c)
		}
	}
}
func (f *fork) start(port string) {  //start fork process
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, lerr := listener.Accept()
		if lerr != nil {
			log.Panic(lerr)
		}

		go f.handleNewConnection(conn)
	}
}
func (f *fork) registerFork() error { // fork registration
	Msg := Register{
		Type:  1,
		Id:    f.id,
		Addrs: f.host,
	}
	buf, err := json.Marshal(Msg)
	if err != nil {
		log.Panic(err)
	}
	if _, werr := f.conn.Write(buf); werr != nil {
		return werr
	}
	return nil
}
func (f *fork) handle(reqType int, fid int, pid int, c net.Conn) {
	log.Println("In handle", f.status)
	f.safe.Lock()
	defer f.safe.Unlock()
	if reqType == 1 {
		if f.status == "free" && fid == f.id {
			f.pid = pid
			f.status = "busy"
			resp := Response{
				Id:     f.id,
				Status: true,
			}
			msg, err := json.Marshal(resp)
			if err != nil {
				log.Panic(err)
			}
			if _, err := c.Write(msg); err != nil {
				log.Println("Message sending failed")
			}
		} else {
			resp := Response{
				Id:     f.id,
				Status: false,
			}
			msg, err := json.Marshal(resp)
			if err != nil {
				log.Panic(err)
			}
			if _, err := c.Write(msg); err != nil {
				log.Println("Message sending failed")
			}
		}

	} else if reqType == 2 && f.pid == pid && f.id == fid && f.status == "busy" {
		f.pid = -1
		f.status = "free"

	}

}
func main() {
	flag.Parse()
	con, err := net.Dial("tcp", *host)
	if err != nil {
		log.Panic(err)
	}
	f := NewFork(*fid, *port, con)
	f.registerFork()
	f.start(*port)
}
