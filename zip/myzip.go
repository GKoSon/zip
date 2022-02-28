package zip

import (
	"bufio"
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
		//panic("failed to " + action + ": " + err.Error())
		fmt.Printf("failed to " + action + ": " + err.Error())
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		fmt.Printf("Input char is:%v", string([]byte(input)))
		panic("failed")
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
func CheckFileMd5() bool {
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
