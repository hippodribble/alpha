package main

import (
	// "time"

	g "github.com/hippodribble/alpha/utils/geometry"
	"github.com/hippodribble/alpha/utils/gps"
)

func main(){
	a:=g.Point{X:3,Y:4,
		Label:"f"}
	println(a.Label)

	server:=*new(gps.GPSDevice)
	server.Baud=57600
	server.Port="COM8"
	server.TCPOut=":5544"

	server.StartGPS()
}