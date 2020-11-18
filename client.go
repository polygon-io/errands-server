//nolint:golint // TODO
package main

import (
	"errors"
	"strings"

	gin "github.com/gin-gonic/gin"
)

type Client struct {
	ErrandsServer *ErrandsServer
	Notifications chan *Notification
	Gin           *gin.Context
	EventSubs     []string
}

func (s *ErrandsServer) RemoveClient(c *Client) {
	close(c.Notifications)
	s.UnregisterClient <- c
}

func (s *ErrandsServer) NewClient(c *gin.Context) (*Client, error) {
	obj := &Client{
		Notifications: make(chan *Notification, 10),
		ErrandsServer: s,
		Gin:           c,
	}
	events := c.DefaultQuery("events", "*")

	obj.EventSubs = strings.Split(events, ",")
	if len(obj.EventSubs) < 1 {
		return obj, errors.New("must have at least 1 event subscription")
	}
	s.RegisterClient <- obj

	return obj, nil
}

func (c *Client) Gone() {
	c.ErrandsServer.RemoveClient(c)
}
