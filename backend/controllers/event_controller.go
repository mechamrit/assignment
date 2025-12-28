package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"backend/realtime"

	"github.com/gin-gonic/gin"
)

type Event struct {
	broadcaster realtime.Broadcaster
}

func NewEvent(broadcaster realtime.Broadcaster) *Event {
	return &Event{
		broadcaster: broadcaster,
	}
}

func (ctrl *Event) StreamEvents(c *gin.Context) {
	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project_id"})
		return
	}

	clientChan := ctrl.broadcaster.Subscribe(uint(projectID))
	defer ctrl.broadcaster.Unsubscribe(uint(projectID), clientChan)

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-clientChan:
			if !ok {
				return false
			}
			c.SSEvent("message", msg)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})

	fmt.Printf("Client disconnected from SSE for project %d\n", projectID)
}
