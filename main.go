package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

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

func findFile(originDir string) string {
	fmt.Println(originDir)
	files, err := ioutil.ReadDir(originDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := file.Name()
		if len(name) > 4 {
			if string(name[len(name)-4:len(name)]) == ".gpj" {
				return originDir + name
			}
		}
	}
	return ""
}

func cleanup(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	// After 8 hours running
	var modTime time.Time
	var names []string
	if len(files) > 32 {
		for _, fi := range files {
			if fi.Mode().IsRegular() && fi.Name() != ".DS_STORE" {
				if !fi.ModTime().Before(modTime) {
					if fi.ModTime().After(modTime) {
						modTime = fi.ModTime()
						names = names[:0]
					}
					names = append(names, fi.Name())
				}
			}
		}
		// Move to folder tagged with modTime
		destDir := path + modTime.Format("2006-01-02 15:04:05") + "/"
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			err = os.Mkdir(destDir, 0700)

			if err != nil {
				panic(err)
			}
		}
		_, err = copy(path+names[0], destDir+names[0])
		for _, fi := range files {
			os.Remove(path + fi.Name())
		}

	}
}

func xtract() {
	// addPath()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	originDir := dir
	if originDir == "\\" {
		panic("Empty path")
	}
	if runtime.GOOS == "windows" {
		originDir += "\\"
	}
	originFile := findFile(originDir)
	_, originFileName := filepath.Split(originFile)

	destDir := originDir + "00_SS/"

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.Mkdir(destDir, 0700)

		if err != nil {
			panic(err)
		}
	}

	cleanup(destDir)
	ext := filepath.Ext(originFileName)
	cleanFile := strings.Replace(originFileName, ext, "", -1)

	currentTime := time.Now()
	destFile := destDir + cleanFile + "_" + currentTime.Format("20060102150405") + ext

	_, err = copy(originFile, destFile)
	if err != nil {
		panic(err)
	}

}

func main() {
	for true {
		go xtract()
		time.Sleep(15 * time.Minute)
	}
}
