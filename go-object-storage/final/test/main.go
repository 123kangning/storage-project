package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.Open("/home/kangning/speak.txt")
	if err != nil {
		log.Println("open error , ", err)
		return
	}

	for i := 0; i < 10; i++ {
		buf := make([]byte, 1563)
		file.Read(buf)
		log.Println("file = ", string(buf))
		//info, _ := file.Stat()
		//log.Println("file = ", info)
		//log.Println("hash = ", utils.CalculateHash(file))
	}
}
