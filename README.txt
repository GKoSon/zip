功能:调用github.com/alexmullins/zip完成文件压缩/解压
用途:可执行文件的OTA升级
操作记录:
本地
1---新建文件夹
C:\Users\Koson.Gong\Desktop\pi_ota_tool
2---文件夹内启动VSCODE 新建文件main.go
package main
import "fmt"

func main() {
    fmt.Println("HELLO WORLD")
}
3---模仿fmt 引入包
"github.com/alexmullins/zip"
它的格式是【github.com】+【作者】+【仓库名称】
此后是可以使用仓库名称-包名的函数方法的  
注意:不能直接使用 zip_test.ExampleWrite_Encrypy()
因为它不是zip这个包里面的 需要参考它的逻辑自己写

4---开始操作mod
go mod init github.com/GKoSon/zip   #【github.com】+【作者】+【仓库名称】
go mod tidy

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

因为 mod里面是 github.com/GKoSon/zip
所以我建立这个zip的仓库

5--写代码
package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/looklzj/zip"
)

//COMMON
var PASSWD string = "KOSON"

//ZIP-->FILE
var ZIP_PATH string = "./rpi_gw.zip"

//FILE->ZIP
var FILE_NAME string = "rpi_gw"

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
func Out() {
	r, err := zip.OpenReader(ZIP_PATH)
	must("zip.OpenReader", err)
	defer r.Close()
	NEW_DIR_PATH := ZIP_PATH[0:strings.Index(ZIP_PATH, ".zip")]
	fmt.Printf("NEW_DIR_PATH = %s\n", NEW_DIR_PATH)

again:
	err = os.Mkdir(NEW_DIR_PATH, 0777)
	if err != nil {
		must("os.RemoveAll", os.RemoveAll(ZIP_PATH))
		goto again
	}

	for _, f := range r.File {
		fmt.Printf("f.Name = %s\n", f.Name)
		NEW_FILE_PATH := NEW_DIR_PATH + "/" + f.Name
		fmt.Printf("NEW_FILE_PATH = %s\n", NEW_FILE_PATH)

		of, err := os.Create(NEW_FILE_PATH)
		must("os.Create", err)

		f.SetPassword(PASSWD)
		rr, err := f.Open()
		must("io.Open", err)
		//_, err = io.Copy(os.Stdout, rr)
		_, err = io.Copy(of, rr)

		must("io.Copy", err)
		rr.Close()
		of.Close()
	}
}
func checkFileMd5() bool {
	FILE_PATH := "./" + FILE_NAME + "/" + FILE_NAME
	MD5_FILE_PATH := "./" + FILE_NAME + "/" + FILE_NAME + ".md5"
	//fmt.Println(FILE_PATH)
	//fmt.Println(MD5_FILE_PATH)
	pFile, err := os.Open(FILE_PATH)
	if err != nil {
		fmt.Printf("打开文件失败,filename=%v, err=%v", FILE_PATH, err)
		return false
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)
	mymd5 := hex.EncodeToString(md5h.Sum(nil))

	youmd5, err := os.ReadFile(MD5_FILE_PATH)
	if err != nil {
		fmt.Printf("打开文件失败,filename=%v, err=%v", MD5_FILE_PATH, err)
		return false
	}
	fmt.Println(hex.EncodeToString(youmd5))
	fmt.Println(mymd5)
	return mymd5 == hex.EncodeToString(youmd5)
}

// 获取文件的md5码
func getFileMd5(filename string) string {
	// 文件全路径名
	path := fmt.Sprintf("./%s", filename)
	pFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("打开文件失败,filename=%v, err=%v", filename, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil))
}

func In() {
	//contents := []byte("Hello World")
	FILE_PATH := "./" + FILE_NAME
	contents, err := os.ReadFile(FILE_PATH) //打开文件
	must("os.ReadFile", err)

	ZIP_PATH := FILE_PATH + ".zip"
	fzip, err := os.Create(ZIP_PATH) //建立ZIP文件
	must("os.Create", err)

	zipw := zip.NewWriter(fzip)
	defer zipw.Close()
	w, err := zipw.Encrypt(FILE_NAME, PASSWD) //ZIP里面的文件名称是FILE_NAME 密码是KOSON
	must("zipw.Encrypt", err)

	_, err = io.Copy(w, bytes.NewReader(contents)) //把文件A写道文件B
	must("io.Copy", err)
	zipw.Flush()

	FILE_MD5 := FILE_NAME + "." + "md5"
	w, err = zipw.Encrypt(FILE_MD5, PASSWD)
	must("zipw.Encrypt", err)
	md5, _ := hex.DecodeString(getFileMd5(FILE_PATH))
	_, err = io.Copy(w, bytes.NewReader(md5))
	must("io.Copy", err)
	zipw.Flush()
}

func main() {
	Out()
	fmt.Println(checkFileMd5())
}


此时测试完成了对应只有2个接口
一个是IN 一个是OUT 
它是从ZIP角度出发的 
如果你有.ZIP文件   那你应该使用OUT函数  获得解压缩文件
如果你没有.ZIP文件 那你应该使用IN函数   获得.ZIP文件

测试效果:
提供一个文件 名字是rpi_gw
执行IN可以获得ZIP文件 
执行OUT可以获得文件夹

6--文件拆分
前面只在一个mian.go完成所有事情
现在抽象一点 
在外面建立文件夹 zip
在文件夹zip里面建立文件 名字无所谓X.go但是包是 package zip
此时修改main.go
package main

import (
	"fmt"

	x "github.com/GKoSon/zip/zip"
)

func main() {
	x.In()
	fmt.Println("HELLO WORLD")
}

++++++++++++++++++++++++++
++++++++++++++++++++++++++
        
git clone git@github.com:GKoSon/zip.git
删除所有文件 只保留main.go
通过VSCODE打开
直接run main出错
开始
PS C:\Users\Koson.Gong\3D Objects\zip> go mod init test
PS C:\Users\Koson.Gong\3D Objects\zip> go mod tidy

此时可以正常WORK了 一行代码都没有修改
当然也可以修改为主流风格
package main

import (
	"fmt"

	"github.com/GKoSon/zip/zip"
)

func main() {
	zip.In()
	fmt.Println("HELLO WORLD")
}

******************对比***************
package main

import (
	"fmt"

	x "github.com/GKoSon/zip/zip"
)

func main() {
	x.In()
	fmt.Println("HELLO WORLD")
}


到这里全部流程结束



增加新功能:
1---
win10执行避免黑框框一闪而过 需要一个类似 百度 go的getchar函数
我放在must里面 比如来料加工 没有来料 就会悬停CMD
2--
加压以后文件MD5的比较 需要大写 后面才可以开发API
也就现在有3个API了 IN+OUT+CHECK