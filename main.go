package main

import(
	g "github.com/hippodribble/testignore/geometry"
	K "github.com/hippodribble/alpha/utils"
)

func main(){
	println(K.C)
	p:=g.Point{3,4,"Smith"}
	println(p.Label)
}