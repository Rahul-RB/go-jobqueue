package main

import (
	"github.com/Rahul-RB/go-jobqueue/routes"
	"github.com/Rahul-RB/go-jobqueue/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(utils.InjectNats(utils.CreateStreamAndConsumer()))

	v1Router := router.Group("/v1")
	v1Router.POST("/job", routes.PostJob)   // Create a job
	v1Router.GET("/job/:id", routes.GetJob) // Get job metadata
	// v1Router.GET("/job/:id/output")       // Get job metadata
	// v1Router.DELETE("/job/:id")           // Delete the job

	router.Run(":3000")

}
