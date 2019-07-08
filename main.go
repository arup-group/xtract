package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Accept path from standard args indicated by -p and add it to the local file
func addPath() {
	if len(os.Args) > 1 && os.Args[1] == "-p" {
		path := os.Args[2] + "\n"

		f, err := os.OpenFile("./paths", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		n, err := w.WriteString(path)
		fmt.Printf("wrote %d bytes\n", n)
		w.Flush()
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func xtract() {
	addPath()

	f, err := os.Open("./paths")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		originFile := scanner.Text()
		originDir, originFileName := filepath.Split(originFile)
		destDir := originDir + "00_SS/"

		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			err = os.Mkdir(destDir, 0700)

			if err != nil {
				panic(err)
			}
		}

		fmt.Print(originFileName, "\n")
		ext := filepath.Ext(originFileName)
		cleanFile := strings.Replace(originFileName, ext, "", -1)

		currentTime := time.Now()
		destFile := destDir + cleanFile + "_" + currentTime.Format("20060102150405") + ext

		copy(originFile, destFile)
	}
}

func main() {
	for true {
		go xtract()
		time.Sleep(15 * time.Minute)
	}
}
