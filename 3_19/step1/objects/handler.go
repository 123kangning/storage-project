package objects

import (
	"encoding/json"
	"log"
	"net/http"
)

type BaseResp struct {
	StatusCode    int
	StatusMessage string
}

// Handler
/**
ResponseWriter，Request我就不解释了
*/
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		resp := BaseResp{}
		err := put(w, r)
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMessage = err.Error()
		} else {
			resp.StatusMessage = "success"
		}
		jsonByte, err := json.Marshal(resp)
		if err != nil {
			log.Println("json.Marshal error = ", err)
		}
		w.Write(jsonByte)
		return
	}
	if m == http.MethodGet {
		resp := BaseResp{}
		err := get(w, r)
		if err != nil {
			resp.StatusCode = 1
			resp.StatusMessage = err.Error()
		} else {
			resp.StatusMessage = "success"
		}
		jsonByte, err := json.Marshal(resp)
		if err != nil {
			log.Println("json.Marshal error = ", err)
		}
		w.Write(jsonByte)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
