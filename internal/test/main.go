package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func main() {
	zip()
	ziP()
	//unzip()
	//unZiP()
}
func unZiP() {
	zip, _ := os.Open("/home/kangning/go-project/speak.unzip")
	buf := make([]byte, 10000)
	zip.Read(buf)
	fmt.Println("unzip = ", string(buf))
	zip.Close()
}
func unzip() {
	zip, _ := os.Open("/home/kangning/go-project/speak.zip")
	unzip, _ := os.Create("/home/kangning/go-project/speak.unzip")
	r, _ := gzip.NewReader(zip)
	io.Copy(unzip, r)
	unzip.Close()
	zip.Close()
}
func ziP() {
	zip, _ := os.Open("/home/kangning/go-project/speak.zip")
	buf := make([]byte, 10000)
	zip.Read(buf)
	fmt.Println("zip = ", string(buf))
	zip.Close()
}
func zip() {
	file, _ := os.Open("/home/kangning/go-project/speak.txt")
	zip, _ := os.Create("/home/kangning/go-project/speak.zip")
	w := gzip.NewWriter(zip)
	io.Copy(w, file)
	w.Close()
	zip.Close()
	file.Close()
}
