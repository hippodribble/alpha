package gps

type GPSDevice struct{
	Port string
	Baud int
	TCPOut string
}

func(g *GPSDevice) StartGPS(){
	println(g.Port,g.Baud, g.TCPOut)
}

func(g *GPSDevice) StopGPS(){
	println("Stopping GPS...")
}

func(g *GPSDevice) HandleGPS(){
	println("GPS Running...")
}
