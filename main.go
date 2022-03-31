package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/GKoSon/zip/zip"
)

func CreateBash() {
	os.Remove("./git.bash")
	f, err := os.Create("./git.bash")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	defer f.Close()

	git_bash := "echo \"--------GitStart--------\"\n" +
		"git init\n" +
		"git add .\n" +
		"git status\n" +
		"git config --global user.email \"you@example.com\"\n" +
		"git config --global user.name \"Your Name\"\n" +
		"git commit -m $1\n" +
		"echo \"--------GitEnd--------\"\n"

	//fmt.Println(git_bash)
	_, err = f.WriteString(git_bash)
	if err != nil {
		fmt.Println("err = ", err)
	}
	os.Chmod("./git.bash", 0777)
}

func main() {

	//步骤1 git
	os.RemoveAll("./.git")
	CreateBash()
	cmd := exec.Command("bash", "./git.bash")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("[pack]exec.Command fail %v\r\n", err)
		os.Remove("./git.bash")
		return
	}
	fmt.Printf("[pack]exec.Command ok %s\r\n", stdout)
	os.Remove("./git.bash")
	//步骤2 zip

	/*俄罗斯套娃做法*/
	/*正向*/
	zip.ZipDir("./.git", "./TEMP.ZIP")
	zip.ZipFile("./TEMP.ZIP", "./IMG.ZIP", "IMG") //把压缩包TEMP.ZIP再次压缩 压后名字是IMG.ZIP 它里面有一个文件+一个MD5文件 其实这个文件是一个压缩包
	os.Remove("./TEMP.ZIP")
	fmt.Println("DONE")

	/*逆向*/
	/*
		zip.UNZipFile("./IMG.ZIP") //解压 出来一个文件夹 名字是上面的参数3 "IMG"
		if zip.CheckFileMd5("./IMG/IMG", "./IMG/IMG.md5") {
			//zip.UNZipDir("./IMG/IMG", "./IMG/.git")//来源是当前目录的IMG文件夹里面的IMG文件 产出的当前目前的IMG文件夹 它里面在间一个.git文件夹 里面放置解压出来的文件夹
			zip.UNZipDir("./IMG/IMG", "./IMG") //产出是当前的IMG文件夹 直接出现压缩包里面的东西 比上面少一层
			os.Remove("./IMG/IMG")
			os.Remove("./IMG/IMG.md5")
			fmt.Println("DONE")
		} else {
			fmt.Println("FAIL")
		}
	*/
}
