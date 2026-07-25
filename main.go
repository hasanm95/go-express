package main

import (
	"net/http"

	"github.com/hasanm95/go-express/goexpress"
)

func main (){
	g := goexpress.NEW()

	g.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Go Framework Lab 1!\n"))
	})

	g.Run(":8080")
}