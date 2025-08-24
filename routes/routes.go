package routes

import (
	"log"
	"net/http"

	"github.com/Rahul-RB/go-jobqueue/jobs"
	"github.com/Rahul-RB/go-jobqueue/stream"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func PostJob(c *gin.Context) {
	j := jobs.NewJob(c.MustGet("stream").(*stream.Stream))
	go j.Run()

	c.IndentedJSON(http.StatusOK, j)
}

func GetJob(c *gin.Context) {
	_id := c.Param("id")
	log.Println("ID:", _id)
	if j, err := jobs.GetJob(_id); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	} else {
		c.IndentedJSON(http.StatusOK, j)
	}
}

func StreamOutput(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket:", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, "Failed to upgrade to Websocket connection: "+err.Error())
		return
	}

	// get id
	jobId := c.Param("id")
	j, err := jobs.GetJob(jobId)

	if err != nil {
		log.Println("Failed to find job:", err.Error())
		c.Abort()
		return
	}

	// Create new session for this job
	if err := j.StartConsumer(conn); err != nil {
		log.Println("Failed to start consumer:", err.Error())
		c.Abort()
		return
	}
}
