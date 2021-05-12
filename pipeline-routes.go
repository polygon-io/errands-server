package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/polygon-io/errands-server/schemas"
	"github.com/polygon-io/errands-server/utils"
)

const dryRunParam = "dryRun"

func (s *ErrandsServer) createPipeline(c *gin.Context) {
	var pipeline schemas.Pipeline
	if err := c.ShouldBindJSON(&pipeline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "pipeline validation failed",
			"error":   err.Error(),
		})

		return
	}

	if err := pipeline.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "pipeline validation failed",
			"error":   err.Error(),
		})

		return
	}

	// If this was just a dry run, we're done
	dryRun, _ := strconv.ParseBool(c.Query(dryRunParam))
	if dryRun {
		c.JSON(http.StatusOK, gin.H{"message": "pipeline validated successfully"})
		return
	}

	// Set the ID and status of the pipeline
	pipeline.ID = uuid.New().String()
	pipeline.Status = "inactive"
	pipeline.StartedMillis = utils.GetTimestamp()

	// Initialize all the errands in the pipeline
	for _, errand := range pipeline.Errands {
		errand.SetDefaults()
		errand.Status = schemas.StatusBlocked

		if err := s.saveErrand(errand); err != nil {
			// This should never really happen
			log.WithError(err).Error("error saving errand")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal store error"})
			return
		}
	}

	// Kick off all the unblocked errands now
	for _, unblockedErrand := range pipeline.GetUnblockedErrands() {
		unblockedErrand.Status = schemas.StatusInactive
		s.ErrandStore.SetDefault(unblockedErrand.ID, *unblockedErrand)
		s.AddNotification("created", unblockedErrand)
	}

	s.PipelineStore.SetDefault(pipeline.ID, pipeline)

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": pipeline,
	})
}

func (s *ErrandsServer) getPipeline(pipelineID string) (schemas.Pipeline, bool) {
	pipeline, exists := s.PipelineStore.Get(pipelineID)
	if !exists {
		return schemas.Pipeline{}, false
	}

	return pipeline.(schemas.Pipeline), true
}

func (s *ErrandsServer) updateErrandInPipeline(errand *schemas.Errand) {
	pipeline, exists := s.getPipeline(errand.PipelineID)
	if !exists {
		return
	}

	// Update the pipeline's internal representation of the errand
	for i, pipelineErrand := range pipeline.Errands {
		if errand.ID == pipelineErrand.ID {
			pipeline.Errands[i] = errand
			break
		}
	}

	// Check for any newly unblocked errands
	for _, unblockedErrand := range pipeline.GetUnblockedErrands() {
		// If this unblocked errand is already in progress, just continue
		if unblockedErrand.Status != schemas.StatusBlocked {
			continue
		}

		unblockedErrand.Status = schemas.StatusInactive
		if err := unblockedErrand.AddToLogs("INFO", "errand unblocked"); err != nil {
			// Log this but continue
			log.WithError(err).Error("unable to add to errand logs")
		}

		s.ErrandStore.SetDefault(unblockedErrand.ID, *unblockedErrand)
	}

	// Save the updated pipeline
	s.PipelineStore.SetDefault(pipeline.ID, pipeline)
}