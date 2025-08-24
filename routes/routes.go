package routes

import (
	"log"
	"net/http"

	"github.com/Rahul-RB/go-jobqueue/jobs"
	"github.com/Rahul-RB/go-jobqueue/stream"
	"github.com/gin-gonic/gin"
)

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
