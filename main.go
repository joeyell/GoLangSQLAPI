package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func get_port() string {
	port := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = ":" + val
	}
	return port
}

func main() {
	router := gin.Default()

	// Default home page
	router.GET("/api/home", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Home page!",
		})
	})

	// Retrieve count of crew members
	router.GET("/api/count", handleCount)

	// Retrieve crew member information
	router.GET("/api/crew", handleEntireCrew)

	// Retrieve crew member information
	router.GET("/api/crew/:id", handleCrewMember)

	port := get_port()
	router.Run(port)
}
