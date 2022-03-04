package main

import (
	"fmt"

	x "github.com/GKoSon/zip/zip"
)

func main() {
	x.ZipFile("./rpi_gw", "./img.zip", "rpi_gw")
	//x.UNZipFile("./img.zip")

	fmt.Println("HELLO WORLD")
}
