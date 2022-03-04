package main

import (
	"fmt"

	x "github.com/GKoSon/zip/zip"
)

func main() {
	x.ZipFile()
	//x.ZipDir("./rpi_gw", "./rr.zip")//压缩文件夹
	//x.UNZipDir("./rr.zip", "./rpi_gw2")
	//x.UNZipDir("./rr.zip", "./rpi_gw2")//解压文件夹回去
	fmt.Println("HELLO WORLD")
}
