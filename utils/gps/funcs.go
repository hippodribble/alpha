package gps

func StartGPS(c GPSConnection){
	println(c.Port,c.Baud)
}

type GPSConnection struct{
	Port string
	Baud int
}
