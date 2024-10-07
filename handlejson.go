package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func replyWithError(msg string, code int, w http.ResponseWriter) {
	errMsg := struct {
		Msg string `json:"msg"`
	}{
		Msg: msg,
	}

	if code > 499 {
		log.Printf("Server Error -> %v", code)
	}

	replyWithJson(w, errMsg, code)

}

func replyWithJson(w http.ResponseWriter, msg interface{}, code int) {
	dat, err := json.Marshal(msg)
	if err != nil {
		log.Printf("can not marshal json > %v", err)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
