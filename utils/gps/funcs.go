package gps

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/tarm/serial"
)

type GPSDevice struct {
	Port    string
	Baud    int
	TCPOut  string
	active  bool
	logfile string
	conns   []net.Conn
}

func (g *GPSDevice) StopGPS() {
	println("Stopping GPS...")
	g.active = false
}

func (g *GPSDevice) HandleGPS() {
	println("GPS Running...")
}

func (g *GPSDevice) StartGPS() bool {
	config := &serial.Config{
		Name:        g.Port,
		Baud:        g.Baud,
		ReadTimeout: 1,
		Size:        8,
	}
	gpshandle, err := serial.OpenPort(config)
	if err != nil {
		println("\nFailed to open GPS unit! Check port and speed details")
		g.active = false
		return g.active
	}
	defer gpshandle.Close()
	f, err := os.OpenFile(g.logfile, os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(gpshandle)
	for scanner.Scan() && g.active {
		// fmt.Println(scanner.Text()) // Println will add back the final '\n'
		if g.logfile != "" {
			f.WriteString(scanner.Text() + "\n")
		}
		g.HandleString(scanner.Text())

		// time.Sleep(time.Millisecond*10)
	}
	return true
}

func(g *GPSDevice) HandleString(text string) {
	var delconn int = -1
	for i, conn := range g.conns {
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			delconn = i
		}
	}
	if delconn > -1 {
		g.conns[delconn] = g.conns[len(g.conns)-1]
		g.conns = g.conns[:len(g.conns)-1]
	}

	
}

func(g *GPSDevice) Addconnections() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", g.TCPOut)
	if err != nil {
		log.Fatal(err)
	}
	println("Opening server at", g.TCPOut)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		print("Customer!")
		println(len(g.conns), "connections")
		if err != nil {
			println("Rejected conn")
			continue
		}
		g.conns = append(g.conns, conn)
	}
}