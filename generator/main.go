package main

import (
	"fmt"
	"github.com/acikkaynak/musahit-harita-backend/aws/s3"
)

func main() {
	//conf := ParseConfig()
	//pool := NewDB(conf)
	//UpdateGeolocation(pool)
	//Migrate(pool)
	ob := s3.DownloadMostRecentObject("secim-ovo/sonuc/")
	fmt.Print(string(ob))
}
