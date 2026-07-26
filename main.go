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

// Global Middleware
func GlobalLogger() goexpress.RouteHandler {
	return func(c *goexpress.Context) {
		fmt.Printf("[GLOBAL OMNIPRESENT LOG] Intercepted: %s %s\n", c.Method, c.Path)
		c.Next()
	}
}
func AdminGuard() goexpress.RouteHandler {
	return func(c *goexpress.Context) {
		token := c.Req.Header.Get("X-Admin-Token")
		fmt.Println("[GUARD] Token: ", token)
		if token != "super-secret-admin-pass" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"status": "Rejected",
				"reason": "Administrative clearance token missing or invalid.",
			})
			fmt.Println("[GUARD] Aborting...")
			c.Abort() // Overwrites c.index to 3, stopping the loop machine
			return    // Safely exits this function frame
		}
		fmt.Println("[GUARD] Clearance confirmed. Transitioning inward...")
		c.Next()
	}
}

// Middleware for /something-else Group
func ContextualTracker() goexpress.RouteHandler {
	return func(c *goexpress.Context) {
		fmt.Println("[TRACKER] Request routed directly into the 'Something-Else' cluster.")
		c.Next()
	}
}


func main (){
	g := goexpress.NEW()

	// Apply Global Middleware
	g.Use(GlobalLogger())
	g.Use(goexpress.CORS()) 

	// THE WILDCARD PREFLIGHT CATCHER
    // Any OPTIONS request will match this wildcard. The request will flow through:
    // [GlobalLogger -> CORS -> Empty Handler]
    // But because CORS calls c.Abort(), the empty handler is safely ignored!
    g.OPTIONS("/*cors", func(c *goexpress.Context) {
        // Intentionally left blank. c.Abort() in the CORS middleware stops execution before reaching here.
    })

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

	// ==========================================
	// GROUP 1: /admin
	// ==========================================
	admin := g.Group("/admin")
	admin.Use(AdminGuard())
	{
		// Resolves to: PUT /admin/users/:id
		admin.PUT("/users/:id", func(c *goexpress.Context) {
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

			c.JSON(http.StatusOK, map[string]string{
				"action":  "PUT",
				"message": fmt.Sprintf("User structure for ID %s has been completely overwritten.", id),
			})
		})

		// Resolves to: DELETE /admin/users/:id
		admin.DELETE("/users/:id", func(c *goexpress.Context) {
			id := c.Param("id")
			
			if _, exists := userDB[id]; exists {
				delete(userDB, id)
				c.JSON(http.StatusOK, map[string]string{
					"action":  "DELETE",
					"message": fmt.Sprintf("User record with ID %s purged permanently from system storage.", id),
				})
			} else {
				c.JSON(http.StatusNotFound, map[string]string{
					"action":  "DELETE",
					"message": fmt.Sprintf("User not found with ID %s", id),
				})
			}

		})
	}

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

	// (Optional) Keep the wildcard route from earlier to show multiple features coexisting
    g.GET("/static/*filepath", func(c *goexpress.Context) {
        file := c.Param("filepath")
        c.String(http.StatusOK, "Simulating serving file: "+file)
    })

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