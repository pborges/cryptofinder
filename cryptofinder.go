package main

import (
	"path/filepath"
	"os"
	"flag"
	"log"
	"bytes"
	"fmt"
	"encoding/hex"
)

var foundBanners []string
var foundFiles []string
var infectedHeader []byte
var clean bool

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		data := make([]byte, len(infectedHeader))
		filename := f.Name()[0:len(f.Name()) - len(filepath.Ext(f.Name()))]
		if filename == "HELP_DECRYPT" {
			f, _ := filepath.Abs(path)
			foundBanners = append(foundBanners, filepath.Dir(f))
			log.Println("FOUND BANNER:", f)
			if clean {
				log.Println("Delete:", f)
				os.Remove(f)
			}
		} else {
			f, _ := filepath.Abs(path)
			file, err := os.Open(f)
			if err != nil {
				log.Println("ERR: file open:", f, " ", err)
			} else {
				defer file.Close()
				fs, err := file.Stat();
				if err != nil {
					log.Println("ERR: file stat:", f, " ", err)
				}
				if fs.Size() >= int64(len(infectedHeader)) {
					count, err := file.Read(data)

					if err != nil {
						log.Println("ERR: file read:", f, " ", err)
					} else if count == len(infectedHeader) && bytes.Compare(infectedHeader, data) == 0 {
						foundFiles = append(foundFiles, f)
						log.Println("FOUND INFECTED FILE:", path)
					}
				}
			}
		}
	}
	return nil
}

func uniq(list []string) []string {
	unique_set := make(map[string]bool, len(list))
	for _, x := range list {
		unique_set[x] = true
	}
	result := make([]string, 0, len(unique_set))
	for x := range unique_set {
		result = append(result, x)
	}
	return result
}

func main() {
	infectedFileName := "infected_files.txt"
	infectedDirectoryFileName := "infected_directories.txt"

	flag.Parse()
	if len(os.Args) < 3 {
		fmt.Println("usage: cryptofinder <start directory> <header> [clean]")
		return
	}
	var err error
	infectedHeader, err = hex.DecodeString(flag.Arg(1))
	if err != nil {
		log.Fatal("Unable to decode bytes %v\n", err)
	}

	log.Println("Searching for", hex.EncodeToString(infectedHeader))


	root := flag.Arg(0)
	log.Println("Starting at", root)

	if len(os.Args) >= 4 && flag.Arg(2) == "clean" {
		clean = true
		log.Println("Deleteing banners")
	}

	err = filepath.Walk(root, visit)
	if err != nil {
		log.Fatal("filepath.Walk() returned %v\n", err)
	}

	os.Remove(infectedDirectoryFileName)
	dirLog, err := os.OpenFile(infectedDirectoryFileName, os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer dirLog.Close()

	os.Remove(infectedFileName)
	fileLog, err := os.OpenFile(infectedFileName, os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer fileLog.Close()

	uniqueBanners := uniq(foundBanners)

	for _, banner := range uniqueBanners {
		dirLog.WriteString(banner + "\n")
	}
	for _, file := range foundFiles {
		fileLog.WriteString(file + "\n")
	}
	log.Println("Done!!!")
}