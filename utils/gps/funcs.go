package gps

func StartGPS(c GPSConnection){
	println(c.port,c.baud)
}

type GPSConnection struct{
	port string
	baud int
}
