package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iokr/coder/core/filex"
	"github.com/iokr/coder/core/golang"
)

var dirPath = flag.String("dir", "./CodeStudy/src", "please input dir path")

const (
	destDirPath = "./CodeStudy/dest"
)

func main() {
	flag.Parse()

	if *dirPath == "" || *dirPath == "./" {
		return
	}

	tmpl, err := template.ParseFiles("tmpl/coder.html")
	if err != nil {
		panic(err)
	}

	err = getAllFiles(tmpl, *dirPath, destDirPath, 0)
	if err != nil {
		panic(err)
	}
}

func getAllFiles(tmpl *template.Template, srcDirPath, destDirPath string, dirLevel int) error {
	// 获取目录下所有文件和子目录
	entries, err := os.ReadDir(srcDirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcFullPath := filepath.Join(srcDirPath, entry.Name())
		destFullPath := filepath.Join(destDirPath, entry.Name())
		if entry.IsDir() {
			// 目标文件夹创建目录
			if err := os.Mkdir(destFullPath, 0777); err != nil {
				return err
			}

			// 如果是子目录，则递归获取子目录下的文件
			currentDirLevel := dirLevel + 1
			if err := getAllFiles(tmpl, srcFullPath, destFullPath, currentDirLevel); err != nil {
				return err
			}
		} else {
			// 如果是文件，则进行文件内容转换
			if err := convertToHTMLFile(tmpl, entry.Name(), srcFullPath, destFullPath, dirLevel); err != nil {
				return err
			}
		}
	}
	return nil
}

func convertToHTMLFile(tmpl *template.Template, fileName, srcFilePath, destFilePath string, dirLevel int) error {
	contentByte, err := filex.ReadAll(srcFilePath)
	if err != nil {
		return err
	}

	fd, err := filex.CreateIfNotExist(fmt.Sprintf("%s.html", destFilePath))
	if err != nil {
		return err
	}
	defer fd.Close()

	staticPath := make([]string, 0, dirLevel)
	for i := 0; i < dirLevel; i++ {
		staticPath = append(staticPath, "..")
	}

	staticFullPath := strings.Join(staticPath, "/")
	return tmpl.Execute(fd, map[string]string{
		"fileName":    fileName,
		"codeContent": golang.FormatCode(string(contentByte)),
		"staticPath":  staticFullPath,
	})
}
