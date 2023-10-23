package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/iokr/coder/core/filex"
	"github.com/iokr/coder/core/golang"
)

var dirPath = flag.String("dir", "./CodeStudy", "please input dir path")

func main() {
	flag.Parse()

	if *dirPath == "" || *dirPath == "./" {
		return
	}

	tmpl, err := template.ParseFiles("tmpl/coder.html")
	if err != nil {
		panic(err)
	}

	err = getAllFiles(tmpl, *dirPath)
	if err != nil {
		panic(err)
	}
}

func getAllFiles(tmpl *template.Template, dirPath string) error {
	// 获取目录下所有文件和子目录
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			// 如果是子目录，则递归获取子目录下的文件
			if err := getAllFiles(tmpl, fullPath); err != nil {
				return err
			}
		} else {
			// 如果是文件，则进行文件内容转换
			if err := convertToHTMLFile(tmpl, entry.Name(), fullPath); err != nil {
				return err
			}

			if err := os.Remove(fullPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func convertToHTMLFile(tmpl *template.Template, fileName, filePath string) error {
	contentByte, err := filex.ReadAll(filePath)
	if err != nil {
		return err
	}

	fd, err := filex.CreateIfNotExist(fmt.Sprintf("%s.html", filePath))
	if err != nil {
		return err
	}
	defer fd.Close()

	return tmpl.Execute(fd, map[string]string{
		"fileName":    fileName,
		"codeContent": golang.FormatCode(string(contentByte)),
	})
}
