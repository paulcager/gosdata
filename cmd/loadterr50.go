package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const APP_VERSION = "0.1"

const (
	sourceDir = "/home/paul/os-terrain-50/"
)

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		os.Exit(0)
	}

var count int
	err := filepath.Walk(sourceDir, func(path string, f os.FileInfo, err error) error {
		if err == nil && !f.IsDir() && strings.HasSuffix(f.Name(), ".zip") {
			count++
			err = readZip(path)
		}
		return err
	})

	if err != nil {
		panic(err)
	}
	
	fmt.Println("Count", count)
}

func readZip(path string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	
	return nil
}
