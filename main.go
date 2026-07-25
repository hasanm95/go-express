package main

import (
	"net/http"

	"github.com/hasanm95/go-express/goexpress"
)

// Data model
type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

// The In-Memory Datastore
// Maps a string ID (like "1") to a User struct
var userDB = make(map[string]User)

func main (){
	g := goexpress.NEW()

	// Welcome route
	g.GET("/", func(c *goexpress.Context) {
		c.String(http.StatusOK, "Welcome to Go Framework Lab 1!\n")
	})
	
	// Info route
	g.GET("/api/info", func(c *goexpress.Context) {
		c.JSON(200, map[string]interface{}{
			"framework": "GoExpress",
			"version":   "2.0",
			"author":    "Your Name",
		})
	})

	// Create: Add new user
	g.POST("/user", func(c *goexpress.Context) {
		var newUser User
		if err := c.BindJSON(&newUser); err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON data")
			return;
		}

		// Save to our fake database (Hardcoding ID "1" for simplicity)
		userDB["1"] = newUser
		c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
	})

	// Read: Get the user
	g.GET("/user", func(c *goexpress.Context) {
		if user, exists := userDB["1"]; exists {
			c.JSON(http.StatusOK, user)
		} else {
			c.String(http.StatusNotFound, "user not found")
		}
	})

	// Update: modify user data
	g.PUT("/user", func(c *goexpress.Context) {
		var updatedUser User
		if err := c.BindJSON(&updatedUser); err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON data")
			return
		}

		// Update db
		userDB["1"] = updatedUser
		c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
	})

	// Delete: Delete user
	g.DELETE("/user", func(c *goexpress.Context) {
		delete(userDB, "1")
		c.String(http.StatusOK, "User deleted")
	})


	g.Run(":8080")
}