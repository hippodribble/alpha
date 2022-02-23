package main

import (
	// "time"
	
	"github.com/hippodribble/alpha/utils/gps"
)

func main(){

	server:=*new(gps.GPSDevice)
	server.Baud=57600
	server.Port="COM8"
	server.TCPOut=":5544"

	server.StartGPS()
}