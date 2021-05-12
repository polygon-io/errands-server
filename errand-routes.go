package main

import (
	"errors"
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/polygon-io/errands-server/metrics"
	schemas "github.com/polygon-io/errands-server/schemas"
	utils "github.com/polygon-io/errands-server/utils"
)

const inactive = "inactive"
const active = "active"

type UpdateRequest struct {
	Progress float64  `json:"progress"`
	Logs     []string `json:"logs"`
}

func (s *ErrandsServer) updateErrand(c *gin.Context) {
	var (
		updatedErrand *schemas.Errand
		updateReq     UpdateRequest
	)

	if err := c.ShouldBind(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error":   err.Error(),
		})

		return
	}

	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func(errand *schemas.Errand) error {
		if errand.Status != active {
			return errors.New("errand must be in active state to update progress")
		}
		// Update this errand attributes:
		if updateReq.Progress != 0 {
			if updateReq.Progress < 0 || updateReq.Progress >= 101 {
				return errors.New("progress must be between 0 - 100")
			}
			errand.Progress = updateReq.Progress
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error":   err.Error(),
		})

		return
	}

	s.AddNotification("updated", updatedErrand)
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": updatedErrand,
	})
}

type FailedRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func (s *ErrandsServer) failedErrand(c *gin.Context) {
	var (
		updatedErrand *schemas.Errand
		failedReq     FailedRequest
	)

	if err := c.ShouldBind(&failedReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error":   err.Error(),
		})

		return
	}

	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func(errand *schemas.Errand) error {
		// Update this errand attributes:
		if err := errand.AddToLogs("ERROR", failedReq.Reason); err != nil {
			return err
		}
		errand.Failed = utils.GetTimestamp()
		errand.Status = "failed"
		errand.Progress = 0
		if errand.Options.Retries > 0 {
			// If we should retry this errand:
			if errand.Attempts <= errand.Options.Retries {
				errand.Status = inactive
			} else {
				// If this errand is out of retries
				metrics.ErrandFailed(errand.Type)
			}
		} else {
			// If this errand was not configured with retries
			metrics.ErrandFailed(errand.Type)
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error":   err.Error(),
		})

		return
	}

	s.AddNotification("failed", updatedErrand)
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": updatedErrand,
	})
}

type CompletedRequest struct {
	Results *gin.H `json:"results"`
}

func (s *ErrandsServer) completeErrand(c *gin.Context) {
	var (
		updatedErrand *schemas.Errand
		compReq       CompletedRequest
	)

	if err := c.ShouldBind(&compReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error":   err.Error(),
		})

		return
	}

	shouldDelete := false

	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func(errand *schemas.Errand) error {
		// Update this errand attributes:
		if err := errand.AddToLogs("INFO", "Completed!"); err != nil {
			return err
		}
		errand.Completed = utils.GetTimestamp()
		errand.Status = "completed"
		errand.Progress = 100
		// errand.Results = compReq.Results
		// If we should delete this errand upon completion:
		if errand.Options.DeleteOnCompleted {
			shouldDelete = true
		}

		metrics.ErrandCompleted(errand.Type)
		return nil
	})
	if err == nil && shouldDelete && updatedErrand.ID != "" {
		s.deleteErrandByID(updatedErrand.ID)
	}

	if shouldDelete && updatedErrand.ID != "" {
		s.deleteErrandByID(updatedErrand.ID)
	}

	s.AddNotification("completed", updatedErrand)

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": updatedErrand,
	})
}

func (s *ErrandsServer) retryErrand(c *gin.Context) {
	var updatedErrand *schemas.Errand

	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func(errand *schemas.Errand) error {
		if errand.Status == inactive {
			return errors.New("cannot retry errand which is in inactive state")
		}
		// Update this errand attributes:
		if err := errand.AddToLogs("INFO", "Retrying!"); err != nil {
			return err
		}
		errand.Status = inactive
		errand.Progress = 0
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error":   err.Error(),
		})

		return
	}

	s.AddNotification("retry", updatedErrand)
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": updatedErrand,
	})
}

func (s *ErrandsServer) logToErrand(c *gin.Context) {
	var logReq schemas.Log
	if err := c.ShouldBind(&logReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error":   err.Error(),
		})

		return
	}

	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func(errand *schemas.Errand) error {
		if errand.Status != active {
			return errors.New("errand must be in active state to log to")
		}

		// Update this errand attributes:
		return errand.AddToLogs(logReq.Severity, logReq.Message)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": updatedErrand,
	})
}

func (s *ErrandsServer) deleteErrand(c *gin.Context) {
	s.Store.Delete(c.Param("id"))

	s.deleteErrandByID(c.Param("id"))

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (s *ErrandsServer) deleteErrandByID(id string) {
	s.Store.Delete(id)
}

// UpdateErrandByID Lets you pass in a function which will be called allowing you to update the errand. If no error is returned, the errand will be saved in the DB with the new attributes.
func (s *ErrandsServer) UpdateErrandByID(id string, fn func(*schemas.Errand) error) (*schemas.Errand, error) {
	errandObj, found := s.Store.Get(id)
	if !found {
		return nil, errors.New("errand with this ID not found")
	}

	errand := errandObj.(schemas.Errand)
	if err := fn(&errand); err != nil {
		return nil, errors.New("error in given update function (fn)")
	}

	s.Store.SetDefault(id, errand)

	return &errand, nil
}
