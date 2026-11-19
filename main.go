package main

import (
	"pr-5-go/database"
	"pr-5-go/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		incidents := v1.Group("/incidents")
		{
			incidents.POST("", handlers.CreateIncident)
			incidents.GET("", handlers.GetAllIncidents)
			incidents.GET("/:id", handlers.GetIncidentByID)
			incidents.PATCH("/:id/status", handlers.UpdateIncidentStatus)
			incidents.DELETE("/:id", handlers.DeleteIncident)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}
	}

	r.Run(":8080")
}