package filex

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// CreateIfNotExist creates a file if it is not exists.
func CreateIfNotExist(file string) (*os.File, error) {
	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exist", file)
	}

	return os.Create(file)
}

// WriteLineByAppend 按行写入文件.
//
//	1.文件不存在,会创建一个新的文件后写入.
//	2.文件存在,会在原文件内容后追加写入.
func WriteLineByAppend(fileName string, values []string) (err error) {
	fd, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	valueLength := len(values)
	bufWriter := bufio.NewWriter(fd)
	for idx := range values {
		_, _ = bufWriter.WriteString(values[idx])
		if idx < valueLength-1 {
			_, _ = bufWriter.WriteString("\n")
		}
	}
	return bufWriter.Flush()
}

// WriteLineByTruncate 按行写入文件.
//
//	1.文件不存在,会创建一个新的文件后写入.
//	2.文件存在,会清空原文件内容后写入.
func WriteLineByTruncate(fileName string, values []string) (err error) {
	fd, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	valueLength := len(values)
	bufWriter := bufio.NewWriter(fd)
	for idx := range values {
		_, _ = bufWriter.WriteString(values[idx])

		if idx < valueLength-1 {
			_, _ = bufWriter.WriteString("\n")
		}
	}
	return bufWriter.Flush()
}

// ReadAll 一次性读取文件中所有内容返回.
func ReadAll(fileName string) (values []byte, err error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return io.ReadAll(fd)
}

// ReadLine 按行读取文件内容并返回所有内容.
func ReadLine(fileName string) (values []string, err error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	values = make([]string, 0)
	bufRead := bufio.NewReader(fd)
	for {
		lineByte, _, err := bufRead.ReadLine()
		if err != nil {
			if err == io.EOF {
				return values, nil
			}
			return nil, err
		}
		values = append(values, string(lineByte))
	}
}
