package gps

func StartGPS(c GPSConnection){
	println(c.Port,c.Baud)
}

func StopGPS(){
	
}

type GPSConnection struct{
	Port string
	Baud int
}
