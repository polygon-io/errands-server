
package main


import (
	"io"
	"encoding/json"
	gin "github.com/gin-gonic/gin"
)








type Client struct {
	ErrandsServer 		*ErrandsServer
	Notifications 		chan *Notification
	Gin 				*gin.Context
}

func ( s *ErrandsServer ) RemoveClient( c *Client ){
	close( c.Notifications )
	s.UnregisterClient <- c
}

func ( s *ErrandsServer ) NewClient( c *gin.Context ) *Client {
	obj := &Client{
		Notifications: make(chan *Notification, 10),
		ErrandsServer: s,
		Gin: c,
	}
	
	s.RegisterClient <- obj
	go obj.Broadcaster()
	return obj
}

func ( c *Client ) Gone(){
	c.ErrandsServer.RemoveClient( c )
}

func ( c *Client ) Broadcaster(){
	w := c.Gin.Writer
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")

	// clientGone := w.CloseNotify()
	c.Gin.Stream(func(wr io.Writer) bool {
		for {
			select {
			case t := <-c.Notifications:
				jsonData, _ := json.Marshal( t )
				c.Gin.SSEvent("message", string( jsonData ))
				w.Flush()
				return true
			}
		}
		return false
	})
}



