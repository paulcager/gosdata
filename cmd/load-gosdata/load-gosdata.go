package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

const dataURL = `https://api.os.uk/downloads/v1/products/Terrain50/downloads?area=GB&format=ASCII+Grid+and+GML+%28Grid%29&redirect`

var (
	outDir string
)

func main() {
	flag.StringVarP(&outDir, "outDir", "o", "./data", "Directory to write to")
	flag.Parse()

	os.MkdirAll(outDir, os.ModeDir | os.ModePerm)

	zipName := filepath.Join(outDir, "terrain-50.zip")
	if err := download(zipName) ; err != nil {
		log.Printf("Could not download %q to %q: %s\n", dataURL, zipName, err)
		os.Exit(2)
	}

	if err := unzip(zipName); err != nil {
		log.Printf("Could not unzip %q: %s\n", zipName, err)
		os.Exit(2)
	}
}

func download(fileName string) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		log.Printf("Could not write to output directory %q\n", outDir)
		os.Exit(2)
	}
	defer outFile.Close()

	res, err := http.Get(dataURL)
	if err != nil {
		log.Printf("Could not load %q: %s\n", dataURL, err)
		os.Exit(2)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Could not load %q: %s\n", dataURL, res.Status)
		io.Copy(os.Stderr, res.Body)
		os.Stderr.WriteString("\n")
		os.Exit(2)
	}

	_, err = io.Copy(outFile, res.Body)
	return err
}

func unzip(zipFile string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		fileName := filepath.Clean(f.Name)
		if len(fileName) > 1 && fileName[:2] == ".." {
			// https://snyk.io/research/zip-slip-vulnerability
			return fmt.Errorf("invalid file name %q", f.Name)
		}
		fileName = filepath.Join(outDir, fileName)
		if err := copyEntry(fileName, f); err != nil {
			return err
		}
	}

	return nil
}

func copyEntry(fileName string, f *zip.File) error {
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(filepath.Dir(fileName), os.ModeDir | os.ModePerm)
	file, err := os.Create( fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}
