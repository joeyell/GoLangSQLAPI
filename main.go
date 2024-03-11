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

	// GET count of crew members
	router.GET("/api/count", getCount)

	// GET and POST crew member information
	router.GET("/api/crew", getEntireCrew)
	router.POST("/api/crew", postCrew)

	// GET crew member information
	router.GET("/api/crew/:id", getCrewMember)

	port := get_port()
	router.Run(port)
}
