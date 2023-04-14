package tool

import (
	"io"
	"log"
	"os"
	"time"
)

const (
	//FORMAT .
	FORMAT = "20060102"
	//LineFeed 换行
	LineFeed = "\r\n"
)

var (
	staticPath string
	err        error
)

func init() {
	staticPath, err = os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
}

// WriteLog return error
func WriteLog(msg string) error {
	var (
		err      error
		f        *os.File
		fileName = time.Now().Format("200601021504") + ".log"
	)
	newPath := staticPath + "/filelog/"
	if !IsExist(newPath) {
		err = CreateDir(newPath)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	f, err = os.OpenFile(newPath+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	_, err = io.WriteString(f, LineFeed+time.Now().Format("2006-01-02 15:04:05")+" "+msg)

	defer f.Close()
	return err
}

// CreateDir  文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

// IsExist  判断文件夹/文件是否存在  存在返回 true
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
