package main

import (
	"net/http"

	"github.com/hasanm95/go-express/goexpress"
)

func main (){
	g := goexpress.NEW()

	g.GET("/", func(c *goexpress.Context) {
		c.String(http.StatusOK, "Welcome to Go Framework Lab 1!\n")
	})

	g.GET("/hello", func(c *goexpress.Context) {
		c.String(http.StatusOK, "hello world\n")
	})
	
	g.GET("/api/info", func(c *goexpress.Context) {
		c.JSON(200, map[string]interface{}{
			"framework": "GoExpress",
			"version":   "2.0",
			"author":    "Your Name",
		})
	})
	g.Run(":8080")
}