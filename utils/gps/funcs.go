package gps

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// Represents a serial GPS device
type GPSDevice struct {
	Port    string
	Baud    int
	TCPOut  string
	active  bool
	conns   []net.Conn
}

type GPSFloatField struct{
	Label string
	Timetag time.Time
	value float64
	filterlength int `default:"1"`
	history []float64
}

func(f *GPSFloatField) SetFilterLength(n int){
	f.history = make([]float64,n)
}

func(f *GPSFloatField) GetFilterLength() int{
	return f.filterlength
}

func(f *GPSFloatField) SetValue(v float64, t time.Time) {
	if f.filterlength==1{
		f.value = v
		return
	}
	f.history = append(f.history[1:],v)
	sum:=0.0
	for i:=0;i<len(f.history);i++{
		sum += f.history[i]
	}
	f.value = sum/float64(len(f.history))
	f.Timetag = t
}
func(f *GPSFloatField) GetValue() float64{
	return f.value
}


// Start the GPS server (serial to TCP)
func (g *GPSDevice) StartGPS() bool {
	config := &serial.Config{
		Name:        g.Port,
		Baud:        g.Baud,
		ReadTimeout: 1,
		Size:        8,
	}
	g.active=true
	go g.addconnections()
	gpshandle, err := serial.OpenPort(config)
	if err != nil {
		println("\nFailed to open GPS unit! Check port and speed details")
		g.active = false
		return g.active
	}
	defer gpshandle.Close()
	scanner := bufio.NewScanner(gpshandle)
	for scanner.Scan() && g.active {
		g.handleString(scanner.Text())
	}
	println("BACK")
	return true
}
// Stop the GPS server
func (g *GPSDevice) StopGPS() {
	println("Stopping GPS...")
	g.active = false
}

// serves NMEA data to clients
func (g *GPSDevice) handleString(text string) {
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
	g.processNMEA(text)
}

func (g *GPSDevice) processNMEA(text string){
	if len(text)<6{return}
	// println("Process",text[3:6])
	arr:=strings.Split(text,",")
	if len(arr)<3{
		println("Short")
		return
	}
	switch text[3:6]{
	case "RMC":

	}


}

//manages connections to network clients
func(g *GPSDevice) addconnections() {

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