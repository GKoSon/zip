package zip

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alexmullins/zip"
)

//COMMON
var PASSWD string = "KOSON"

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
func UNZipFile(zipFile string) {
	r, err := zip.OpenReader(zipFile)
	must("zip.OpenReader", err)
	defer r.Close()

	fmt.Printf("zipFile = %s \n", zipFile)

	NEW_DIR_PATH := zipFile[0 : len(zipFile)-4] //move 后缀.zip
again:
	err = os.Mkdir(NEW_DIR_PATH, 0777)
	if err != nil {
		must("os.RemoveAll", os.RemoveAll(NEW_DIR_PATH))
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
func CheckFileMd5(afile, md5file string) bool {
	FILE_PATH := afile
	MD5_FILE_PATH := md5file
	fmt.Println(FILE_PATH)
	fmt.Println(MD5_FILE_PATH)
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
func GetFileMd5(filename string) string {
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

/*
aFile一个客观存在的文件的路径全名
zipfile 需要生成的压缩文件的全名
filenameinzip 压缩文件内文件的名字 建议和aFile最后一样
*/
func ZipFile(aFile, zipFile, filenameinzip string) {
	//contents := []byte("Hello World")
	FILE_PATH := aFile
	contents, err := os.ReadFile(FILE_PATH) //打开文件
	must("os.ReadFile", err)

	ZIP_PATH := zipFile
	fzip, err := os.Create(ZIP_PATH) //建立ZIP文件
	must("os.Create", err)

	zipw := zip.NewWriter(fzip)
	defer zipw.Close()
	w, err := zipw.Encrypt(filenameinzip, PASSWD) //ZIP里面的文件名称是FILE_NAME 密码是KOSON
	must("zipw.Encrypt", err)

	_, err = io.Copy(w, bytes.NewReader(contents)) //把文件A写到文件B
	must("io.Copy", err)
	zipw.Flush()

	FILE_MD5 := filenameinzip + "." + "md5"
	w, err = zipw.Encrypt(FILE_MD5, PASSWD)
	must("zipw.Encrypt", err)
	md5, _ := hex.DecodeString(GetFileMd5(FILE_PATH))
	_, err = io.Copy(w, bytes.NewReader(md5))
	must("io.Copy", err)
	zipw.Flush()
}

/*文件夹操作*/
/*
Dirzip("C:\\Users\\Koson.Gong\\Desktop\\XX\\IMG", "C:\\Users\\Koson.Gong\\Desktop\\XX\\IMG.zip")
桌面一个文件夹IMG压缩为参数2指定的文件
*/
func ZipDir(dir, zipFile string) {
	fz, err := os.Create(zipFile)
	must("os.Create", err)
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			//fDest, err := w.Create(path[len(dir)+1:])
			fDest, err := w.Create(path)
			must("w.Create", err)
			fSrc, err := os.Open(path)
			must("os.Open", err)
			defer fSrc.Close()
			_, err = io.Copy(fDest, fSrc)
			must("io.Copy", err)
		}
		return nil
	})
}

/*unzipDir("C:\\Users\\Koson.Gong\\Desktop\\XX\\IMG.zip", "D:\\adb\\.git")
桌面一个IMG.zip 加压为参数2指定的文件
*/
func UNZipDir(zipFile, dir string) {
	r, err := zip.OpenReader(zipFile)
	must("zip.OpenReader", err)
	defer r.Close()

	for _, f := range r.File {
		func() {
			path := dir + string(filepath.Separator) + f.Name
			os.MkdirAll(filepath.Dir(path), 0755)
			fDest, err := os.Create(path)
			must(" os.Create", err)
			defer fDest.Close()

			fSrc, err := f.Open()
			must(" f.Open", err)
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			must("io.Copy", err)
		}()
	}
}
