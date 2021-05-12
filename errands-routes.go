package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	gin "github.com/gin-gonic/gin"

	schemas "github.com/polygon-io/errands-server/schemas"
	utils "github.com/polygon-io/errands-server/utils"
)

func (s *ErrandsServer) errandNotifications(c *gin.Context) {
	client, err := s.NewClient(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Creating Subscription",
			"error":   err.Error(),
		})

		return
	}

	w := client.Gin.Writer

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")

	clientGone := client.Gin.Writer.CloseNotify()

	client.Gin.Stream(func(wr io.Writer) bool {
		for {
			select {
			case <-clientGone:
				client.Gone()
				return false
			case t, ok := <-client.Notifications:
				if ok {
					// If we are subscribed to this event type:
					if utils.Contains(client.EventSubs, t.Event) || client.EventSubs[0] == "*" {
						jsonData, _ := json.Marshal(t)
						client.Gin.SSEvent("message", string(jsonData))
						w.Flush()
					}
					return true
				}
				return false
			}
		}
	})
}

func (s *ErrandsServer) createErrand(c *gin.Context) {
	log.Println("creating errand")

	var item schemas.Errand

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Errand validation failed!",
			"error":   err.Error(),
		})

		return
	}

	item.SetDefaults()
	s.ErrandStore.SetDefault(item.ID, item)
	s.AddNotification("created", &item)
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": item,
	})
}

func (s *ErrandsServer) saveErrand(errand *schemas.Errand) error {
	if !utils.ContainsStatus(schemas.ErrandStatuses, errand.Status) {
		return errors.New("invalid errand status state")
	}

	s.ErrandStore.SetDefault(errand.ID, *errand)

	return nil
}

func (s *ErrandsServer) getAllErrands(c *gin.Context) {
	errands := s.GetErrandsBy(func(errand *schemas.Errand) bool {
		return true
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": errands,
	})
}

func (s *ErrandsServer) getFilteredErrands(c *gin.Context) {
	key := c.Param("key")
	value := c.Param("val")

	errands := s.GetErrandsBy(func(errand *schemas.Errand) bool {
		switch key {
		case "status":
			return string(errand.Status) == value
		case "type":
			return errand.Type == value
		default:
			return false
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": errands,
	})
}

type filteredUpdateReq struct {
	Status string `json:"status"`
	Delete bool   `json:"delete"`
}

func (s *ErrandsServer) updateFilteredErrands(c *gin.Context) {
	key := c.Param("key")
	value := c.Param("val")

	var updateReq filteredUpdateReq

	if err := c.ShouldBind(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error":   err.Error(),
		})

		return
	}

	errands := s.GetErrandsBy(func(errand *schemas.Errand) bool {
		switch key {
		case "status":
			return string(errand.Status) == value
		case "type":
			return errand.Type == value
		default:
			return false
		}
	})

	var err error

	for _, errand := range errands {
		if updateReq.Delete {
			s.deleteErrandByID(errand.ID)
		} else if updateReq.Status != "" {
			_, err = s.UpdateErrandByID(errand.ID, func(e *schemas.Errand) error {
				// TODO: update errand in pipeline
				e.Status = schemas.Status(updateReq.Status)
				return nil
			})
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"count":  len(errands),
	})
}

func (s *ErrandsServer) processErrand(c *gin.Context) {
	var procErrand schemas.Errand

	errands := make([]schemas.Errand, 0)
	typeFilter := c.Param("type")

	for _, itemObj := range s.ErrandStore.Items() {
		item := itemObj.Object.(schemas.Errand)

		if item.Status != schemas.StatusInactive {
			continue
		}

		if item.Type != typeFilter {
			continue
		}
		// Add to list of errands we could possibly process:
		errands = append(errands, item)
	}

	if len(errands) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No jobs",
		})

		return
	}

	// Of the possible errands to process, sort them by date & priority:
	sort.SliceStable(errands, func(i, j int) bool {
		return errands[i].Created < errands[j].Created
	})
	sort.SliceStable(errands, func(i, j int) bool {
		return errands[i].Options.Priority > errands[j].Options.Priority
	})

	procErrand = errands[0]

	// We are processing this errand:
	procErrand.Started = utils.GetTimestamp()
	procErrand.Attempts += 1
	procErrand.Status = "active"
	procErrand.Progress = 0.0
	_ = procErrand.AddToLogs("INFO", "Started!")
	_ = s.saveErrand(&procErrand)

	s.AddNotification("processing", &procErrand)
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": procErrand,
	})
}

func (s *ErrandsServer) GetErrandsBy(fn func(*schemas.Errand) bool) []schemas.Errand {
	errands := make([]schemas.Errand, 0)

	for _, itemObj := range s.ErrandStore.Items() {
		errand := itemObj.Object.(schemas.Errand)
		if fn(&errand) {
			errands = append(errands, errand)
		}
	}

	return errands
}

func (s *ErrandsServer) clearErrands(c *gin.Context) {
	duration, err := time.ParseDuration(c.Param("duration"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Duration",
			"error":   err.Error(),
		})

		return
	}

	threshold := time.Now().Add(-duration)
	errands := make([]schemas.Errand, 0)

	for _, itemObj := range s.ErrandStore.Items() {
		item := itemObj.Object.(schemas.Errand)

		var stoppedRunning int64

		switch item.Status {
		case "completed":
			stoppedRunning = item.Completed
		case "failed":
			stoppedRunning = item.Failed
		default:
			continue
		}

		if stoppedRunning >= threshold.UnixNano() {
			continue
		}

		errands = append(errands, item)
	}

	for _, errand := range errands {
		s.deleteErrandByID(errand.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": errands,
	})
}
