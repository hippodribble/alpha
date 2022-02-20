package main

import(
	
	g "github.com/hippodribble/alpha/utils/geometry"
	gps "github.com/hippodribble/alpha/utils/gps"
)

func main(){
	
	p:=g.Point{3,4,"Smith"}
	println(p.Label)

	conn:=gps.GPSConnection{":5444",57600}
	println(conn.Port)
}