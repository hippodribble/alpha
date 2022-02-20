package gps

import (
	"bufio"
	"log"
	"net"

	"github.com/tarm/serial"
)

type GPSDevice struct {
	Port    string
	Baud    int
	TCPOut  string
	active  bool
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
	g.active=true
	go g.Addconnections()
	gpshandle, err := serial.OpenPort(config)
	if err != nil {
		println("\nFailed to open GPS unit! Check port and speed details")
		g.active = false
		return g.active
	}
	defer gpshandle.Close()
	scanner := bufio.NewScanner(gpshandle)
	for scanner.Scan() && g.active {
		g.HandleString(scanner.Text())
	}
	println("BACK")
	return true
}

func (g *GPSDevice) HandleString(text string) {
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
	println("Opening server at", g.TCPOut,tcpAddr.Port)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	println("Listening")
	if err != nil {
		println("Error!")
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		print("Customer!")
		println(conn.RemoteAddr().String())
		
		if err != nil {
			println("Rejected conn")
			continue
		}
		g.conns = append(g.conns, conn)
		println(len(g.conns), "connections")
	}
}
