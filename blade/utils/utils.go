package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func WriteFile(userFile string, content string) {
	fout, err := os.Create(userFile)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	fout.WriteString(content)
}

func WriteTemplate(tplName, filePath, tplContent string, data interface{}) {
	tpl, err := template.New(tplName).Parse(tplContent)
	if err != nil {
		panic(err)
	}
	out, err := os.Create(filePath)
	defer out.Close()
	err = tpl.Execute(out, data)
	if err != nil {
		panic(err)
	}
}

func ReadFile(filename string) (map[string]string, error) {
	var xxx = map[string]string{}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(bytes, &xxx); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}
	return xxx, nil
}

func ReadFileAsByte(path string) []byte {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return fd
}

func GetFileModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}
