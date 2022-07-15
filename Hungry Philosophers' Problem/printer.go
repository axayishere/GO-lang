package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"

	"fmt"
)

type Printer struct {
	n       int
	port    string
	updates []string
}
type UpdateMessage struct {
	PId    int
	Status string
}

func NewPrinter(port string, size int) *Printer {
	p := &Printer{
		n:       size,
		port:    port,
		updates: make([]string, size),
	}
	return p
}
func (p *Printer) handleConnection(c net.UDPConn) { //it will update the message for philosoper
	for {
		var buf [2048]byte
		msg := new(UpdateMessage)
		n, _, err := c.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println("Udp error")
			return
		}

		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			log.Println(err)
		}
		var str string
		for i := 0; i < p.n; i++ {
			if i == msg.PId {
				str = str + " " + msg.Status
			} else {
				str = str + " " + "......."
			}
		}
		log.Println(str)
	}

}
func main() {
	udpPort := flag.String("host", "localhost:19983", "Here is the default Udp Address")
	noOfusers := flag.Int("n", 5, "No.of Philospers")
	flag.Parse()
	Addrs, cerr := net.ResolveUDPAddr("udp", *udpPort)
	if cerr != nil {
		log.Println("Invalid Port no")
		return
	}
	connection, err := net.ListenUDP("udp", Addrs)
	if err != nil {
		log.Panic(err)
	}
	p := NewPrinter(*udpPort, *noOfusers)
	fmt.Println("Listening on udp port", *udpPort)

	p.handleConnection(*connection)

}
