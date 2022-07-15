package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
)

var (
	port = flag.String("host", "localhost:19982", "Ip address ")
	n    = flag.Int("n", 5, "No.of Philospers & Folks")
)

type manager struct {
	n          int
	forks      []string
	philospers []string
}
type Register struct {
	Addrs string
	Type  int
	Id    int
}
type Forks struct {
	Id    int
	Addrs string
}

func readMessage(c net.Conn) []byte {
	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		return nil
	}
	return buffer[:n]
}
func (m *manager) handleNewConnection(c net.Conn) {
	log.Println("New connection")
	for {
		buffer := readMessage(c)
		if buffer == nil {
			continue
		}
		var msg Register

		if err := json.Unmarshal(buffer, &msg); err != nil {
			panic(err)
		}
		log.Println(msg)
		if msg.Type == 1 {
			m.forks[msg.Id] = msg.Addrs
			log.Println(m.forks)
		} else if msg.Type == 2 {
			arr := make([]Forks, m.n)
			m.philospers[msg.Id] = msg.Addrs
			leftFork := Forks{
				Id:    msg.Id,
				Addrs: m.forks[msg.Id],
			}
			rightFork := Forks{
				Id:    (msg.Id + 1) % m.n,
				Addrs: m.forks[(msg.Id+1)%m.n],
			}
			arr[0] = leftFork
			arr[1] = rightFork

			buffer, err := json.Marshal(arr)
			if err != nil {
				log.Panicln(err)
			}
			if _, err := c.Write(buffer); err != nil {
				log.Println(err)
			}
		}

	}

}
func New(n int) *manager {
	m := &manager{
		n:          n,
		forks:      make([]string, n),
		philospers: make([]string, n),
	}
	return m
}
func (m *manager) start(port string) {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Server started")
	for {
		conn, lerr := listener.Accept()
		if lerr != nil {
			log.Panic(lerr)
		}

		go m.handleNewConnection(conn)
	}
}
func main() {
	flag.Parse()
	m := New(*n)
	m.start(*port)

}
