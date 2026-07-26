package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
var idCounter = 1



// Logger is a custom middleware function
func Logger() goexpress.RouteHandler {
	return func(c *goexpress.Context) {
		// 1. PRE-PROCESSING
		t := time.Now()
		fmt.Printf("[START] Request incoming: %s %s\n", c.Method, c.Path)

		// 2. PASS CONTROL to the next middleware or handler
		c.Next()

		// 3. POST-PROCESSING (Executes after the entire request is handled)
		latency := time.Since(t)
		fmt.Printf("[END] Request completed: %s %s in %v\n", c.Method, c.Path, latency)
	}
}

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
	g.POST("/users", func(c *goexpress.Context) {
		var newUser User
		if err := c.BindJSON(&newUser); err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON data")
			return;
		}

		// Generate a new ID as a string
        id := strconv.Itoa(idCounter)
        idCounter++

		// Save to our fake database (Hardcoding ID "1" for simplicity)
		userDB[id] = newUser
		c.JSON(http.StatusCreated, map[string]string{"id": id, "message": "User created successfully"})
	})

    // READ ALL: Get the entire collection of users
    g.GET("/users", func(c *goexpress.Context) {
        if len(userDB) == 0 {
            c.JSON(http.StatusOK, map[string]string{"message": "No users found"})
            return
        }
        c.JSON(http.StatusOK, userDB)
    })

	// READ: Get a specific user by dynamic ID (Targeting an item: /users/:id)
    g.GET("/users/:id", func(c *goexpress.Context) {
        id := c.Param("id")
        
        if user, exists := userDB[id]; exists {
            c.JSON(http.StatusOK, user)
        } else {
            c.String(http.StatusNotFound, "User not found")
        }
    })

    // UPDATE: Modify a specific user by dynamic ID
    g.PUT("/users/:id", func(c *goexpress.Context) {
        id := c.Param("id") // Extract ID from the URL
        
        if _, exists := userDB[id]; !exists {
            c.String(http.StatusNotFound, "User not found")
            return
        }

        var updatedUser User
        if err := c.BindJSON(&updatedUser); err != nil {
            c.String(http.StatusBadRequest, "Invalid JSON data")
            return
        }
        
        userDB[id] = updatedUser
        c.JSON(http.StatusOK, map[string]string{"message": "User " + id + " updated"})
    })

    // DELETE: Remove a specific user by dynamic ID
    g.DELETE("/users/:id", func(c *goexpress.Context) {
        id := c.Param("id")
        
        if _, exists := userDB[id]; exists {
            delete(userDB, id)
            c.String(http.StatusOK, "User " + id + " deleted")
        } else {
            c.String(http.StatusNotFound, "User not found")
        }
    })

	// (Optional) Keep the wildcard route from earlier to show multiple features coexisting
    g.GET("/static/*filepath", func(c *goexpress.Context) {
        file := c.Param("filepath")
        c.String(http.StatusOK, "Simulating serving file: "+file)
    })

	// Register global middleware using .Use()
	g.Use(Logger())

	// A simple route to test the middleware
	g.GET("/heavy", func(c *goexpress.Context) {
		// Simulate some heavy processing
		time.Sleep(200 * time.Millisecond)
		c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to the Onion Model",
		})
	})

	g.GET("/crash", func(c *goexpress.Context) {
		panic("Simulated database failure!")
	})


	g.Run(":8080")
}