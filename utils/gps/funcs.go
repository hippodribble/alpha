package gps

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"github.com/tarm/serial"
)

// Represents a serial GPS device
type GPSDevice struct {
	Port   string
	Baud   int
	TCPOut string
	active bool
	conns  []net.Conn
	hub Hub
}

type Hub struct{
	fields map[string]*GPSFloatField
}
func (m *Hub) Set(key string,val *GPSFloatField){

}

type GPSFloatField struct {
	Label        string
	timetag      time.Time
	value        float64
	stringvalue  string
	filterlength int `default:"1"`
	history      []float64
}

func (f *GPSFloatField) SetFilterLength(n int) {
	f.history = make([]float64, n)
}

func (f *GPSFloatField) GetFilterLength() int {
	return f.filterlength
}

func (f *GPSFloatField) SetValue(v float64, t time.Time) {
	if f.filterlength == 1 {
		f.value = v
		return
	}
	f.history = append(f.history[1:], v)
	sum := 0.0
	for i := 0; i < len(f.history); i++ {
		sum += f.history[i]
	}
	f.value = sum / float64(len(f.history))
	f.timetag = t
}

func (f *GPSFloatField) GetValue() float64 {
	return f.value
}

// Start the GPS server (serial to TCP)
func (g *GPSDevice) StartGPS()  {
	config := &serial.Config{
		Name:        g.Port,
		Baud:        g.Baud,
		ReadTimeout: 1,
		Size:        8,
	}
	g.hub = Hub{ make(map[string]*GPSFloatField,10)}
	println("Made Field map")
	g.active = true
	go g.addconnections()
	gpshandle, err := serial.OpenPort(config)
	if err != nil {
		println("\nFailed to open GPS unit! Check port and speed details")
		g.active = false
	}
	defer gpshandle.Close()
	scanner := bufio.NewScanner(gpshandle)
	for scanner.Scan() && g.active {
		g.handleString(scanner.Text())
	}
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

//manages connections to network clients
func (g *GPSDevice) addconnections() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", g.TCPOut)
	if err != nil {
		log.Fatal(err)
	}
	println("Opening server at", g.TCPOut, tcpAddr.Port)
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

func (g *GPSDevice) processNMEA(text string) {
	if len(text) < 6 {
		return
	}
	// println("Process",text[3:6])
	arr := strings.Split(text, ",")
	if len(arr) < 3 {
		println("Short")
		return
	}
	switch text[3:6] {
	case "GGA":
		println(text)
		if arr[1] == "" {
			return
		}
		field := new(GPSFloatField)
		field.Label = "UTC"
		field.SetFilterLength(5)
		floatval, err := strconv.ParseFloat(arr[1], 64)
		if err != nil {
			return
		}
		field.timetag = time.Now()
		field.stringvalue = arr[1]
		field.value = floatval
		println(field.value, field.stringvalue)
		g.hub.fields["UTC"] = field

		if arr[2] == "" {
			return
		}
		field = new(GPSFloatField)
		field.Label = "UTC"
		field.SetFilterLength(5)
		floatval, err = strconv.ParseFloat(arr[2], 64)
		if err != nil {
			return
		}
		field.timetag = time.Now()
		field.stringvalue = arr[1]
		field.value = floatval
		println(field.value, field.stringvalue)
		g.hub.fields["Latitude"] = field
		
		if arr[4] == "" {
			return
		}
		field = new(GPSFloatField)
		field.Label = "UTC"
		field.SetFilterLength(5)
		floatval, err = strconv.ParseFloat(arr[4], 64)
		if err != nil {
			return
		}
		field.timetag = time.Now()
		field.stringvalue = arr[4]
		field.value = floatval
		println(field.value, field.stringvalue)
		g.hub.fields["Longitude"] = field
	}
}