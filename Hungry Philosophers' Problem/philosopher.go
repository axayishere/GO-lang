package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"time"
)

var (
	n    = flag.Int("n", 5, "No.of Philospers")
	port = flag.String("host", "localhost:19981", "Philosper adress")
	host = flag.String("manager", "localhost:19982", "ip Address of the manager ")
	udp  = flag.String("printer", "localhost:19983", "printer ip address")
	pid  = flag.Int("id", 1, "Philosper ID")
)

type philosper struct { // registration philosoper structure
	id          int
	host        string
	leftFork    *Fork
	rightFork   *Fork
	printerConn net.Conn
	total       int
}
type Fork struct {
	id    int
	conn  net.Conn
	addrs string
}
type Forks struct {
	Id    int
	Addrs string
}
type Request struct {
	Id   int
	PId  int
	Type int
}
type Response struct {
	Id     int
	Status bool
}
type UpdateMessage struct {
	PId    int
	Status string
}
type Register struct {
	Addrs string
	Id    int
	Type  int
}

func New(id int, h string, udp net.Conn, n int) *philosper {
	p := &philosper{
		id:          id,
		host:        h,
		leftFork:    new(Fork),
		rightFork:   new(Fork),
		printerConn: udp,
		total:       n,
	}
	return p
}
func readMessage(c net.Conn) (int, []byte) {
	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		//log.Println(err)
		return 0, nil
	}
	return n, buffer
}
func (p *philosper) register(con net.Conn) error {
	var req Register
	req.Id = p.id
	req.Type = 2
	req.Addrs = p.host
	buf, err := json.Marshal(req)
	if err != nil {
		log.Println("Marshaling Error")
	}
	if _, rerr := con.Write(buf); rerr != nil {
		log.Println("Message sending has failed", rerr)
		return rerr
	}
	n, buffer := readMessage(con)
	var frks []Forks
	if uerr := json.Unmarshal(buffer[:n], &frks); uerr != nil {
		panic(uerr)
	}
	log.Println(frks)
	p.leftFork.id = frks[0].Id
	p.leftFork.addrs = frks[0].Addrs

	p.leftFork.conn, err = net.Dial("tcp", p.leftFork.addrs)
	if err != nil {
		return err
	}
	p.rightFork.id = frks[1].Id
	p.rightFork.addrs = frks[1].Addrs

	p.rightFork.conn, err = net.Dial("tcp", p.rightFork.addrs)
	if err != nil {
		return err
	}
	return nil
}
func (p *philosper) pickRightFork() bool {
	msg := Request{
		Id:   (p.id + 1) % 5,
		PId:  p.id,
		Type: 1}
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}
	p.rightFork.conn.Write(bytes)
	n, buffer := readMessage(p.rightFork.conn)
	var resp Response
	json.Unmarshal(buffer[:n], &resp)
	return resp.Status
}
func (p *philosper) pickLeftFork() bool {
	req := Request{
		Id:   p.id,
		PId:  p.id,
		Type: 1}
	bytes, err := json.Marshal(req)
	if err != nil {
		log.Panic(err)
	}
	p.leftFork.conn.Write(bytes)
	n, buffer := readMessage(p.leftFork.conn)
	// buffer = make([]byte, 1024)
	var resp Response
	if err := json.Unmarshal(buffer[:n], &resp); err != nil {
		log.Panicln(err)
	}
	return resp.Status
}
func (p *philosper) dropRightFork() bool {
	req := Request{
		Id:   (p.id + 1) % 5,
		PId:  p.id,
		Type: 2,
	}
	bytes, err := json.Marshal(req)
	if err != nil {
		log.Panic(err)
	}
	if _, err := p.rightFork.conn.Write(bytes); err != nil {
		return false
	}
	return true
}
func (p *philosper) dropLeftFork() bool {
	req := Request{
		Id:   p.id,
		PId:  p.id,
		Type: 2,
	}
	bytes, err := json.Marshal(req)
	if err != nil {
		log.Panic(err)
	}
	if _, err := p.leftFork.conn.Write(bytes); err != nil {
		return false
	}
	return true
}
func (p *philosper) stateChange(str string) error {
	msg := UpdateMessage{
		PId:    p.id,
		Status: str,
	}
	buffer, err := json.Marshal(msg)
	if err != nil {
		log.Panicln(err)
	}
	if _, werr := p.printerConn.Write(buffer); werr != nil {
		return werr
	}
	return nil
}
func (p *philosper) Run() {
	for {
		p.stateChange("Thinking")
		log.Println("Thinking")
		time.Sleep(3 * time.Second)

		p.stateChange("Waiting")
		log.Println("waiting")
		//	var leftFork, rightFork bool
		leftFork := p.pickLeftFork()
		log.Println("leftfork", leftFork)
		if !leftFork {
			continue
		}
		rightFork := p.pickRightFork()
		if !rightFork {
			p.dropLeftFork()
			continue
		}
		if leftFork && rightFork {
			p.stateChange("Eating")
			log.Println("Eating")
			time.Sleep(3 * time.Second)
			p.dropLeftFork()
			p.dropRightFork()
		}
		time.Sleep(134 * time.Microsecond)
	}
}
func main() {
	flag.Parse()
	//ch := make(chan string, 1)
	addr, err := net.ResolveUDPAddr("udp4", *udp)
	conn, err2 := net.DialUDP("udp4", nil, addr)
	if err2 != nil {
		log.Panic("Coudn't connect to printer", err)
		return
	}
	con, err := net.Dial("tcp", *host)
	if err != nil {
		log.Panic(err)
	}
	p := New(*pid, *port, conn, *n)
	if err1 := p.register(con); err1 != nil {
		log.Panic(err1)
	}
	go p.Run()
	// select {
	// case <-ch:
	// 	break
	// default:
	// }
	for {

	}
}
