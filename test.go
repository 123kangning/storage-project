package main

import (
	"fmt"
	"net/http"
)

func main() {
	r, err := http.Get("http://127.0.0.1:8080/ttms/user/all?Current=0&PageSize=10&Token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTUsIlVzZXJUeXBlIjoxLCJleHAiOjI0MDA4NjcyOTcsImlzcyI6ImthbmduaW5nIn0.TyEkYcxhbRVgH5XoZK_llNpB2m2IdIaqBGEUBTxQoCA")
	if err != nil {
		panic(err)
	}
	data := make([]byte, 100)
	r.Body.Read(data)
	fmt.Println("r = ", r)
	fmt.Println("r.Body = ", string(data))

}
