package main

import (
	"fmt"
	"os"

	"github.com/GKoSon/zip/zip"
)

func main() {
	/*普通做法
	zip.ZipDir("./.git", "./IMG.ZIP") //把文件夹内的( .git文件夹) 压缩为 文件 (IMG.ZIP)
	zip.UNZipDir("./IMG.ZIP", "IMG")  //解压为IMG
	fmt.Println(zip.GetFileMd5("IMG.ZIP"))
	*/

	/*俄罗斯套娃做法*/
	/*正向*/
	zip.ZipDir("./.git", "./TEMP.ZIP")
	zip.ZipFile("./TEMP.ZIP", "./IMG.ZIP", "IMG") //把压缩包TEMP.ZIP再次压缩 压后名字是IMG.ZIP 它里面有一个文件+一个MD5文件 其实这个文件是一个压缩包
	os.Remove("./TEMP.ZIP")
	/*逆向*/
	zip.UNZipFile("./IMG.ZIP") //解压 出来一个文件夹 名字是上面的参数3 "IMG"
	fmt.Println(zip.CheckFileMd5("./IMG/IMG", "./IMG/IMG.md5"))
	zip.UNZipDir("./IMG/IMG", "./IMG/.git")
}
