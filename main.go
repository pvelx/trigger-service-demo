package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pvelx/triggerHook"
	"net/http"
)

var router = gin.Default()
var scheduler = triggerHook.Default()

func main() {

	scheduler.SetTransport(NewTransportAmqp())
	go scheduler.Run()

	router.POST("/task", func(c *gin.Context) {
		var taskRequest taskRequest
		if err := c.ShouldBindJSON(&taskRequest); err != nil {
			c.JSON(http.StatusInternalServerError, "Server error")
			return
		}
		if err := taskRequest.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, "Validation error")
			return
		}

		task, e := scheduler.Create(taskRequest.NextExecTime)
		if e != nil {
			c.JSON(http.StatusInternalServerError, "Something wrong")
			return
		}

		c.JSON(http.StatusOK, task)
	})
	router.Run(":8083")
}
