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

const (
	dryRunQueryParam       = "dryRun"
	idParam                = "id"
	statusFilterQueryParam = "status"
)

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
	dryRun, _ := strconv.ParseBool(c.Query(dryRunQueryParam))
	if dryRun {
		c.JSON(http.StatusOK, gin.H{"message": "pipeline validated successfully"})
		return
	}

	// Set the ID and status of the pipeline
	pipeline.ID = uuid.New().String()
	pipeline.Status = schemas.StatusInactive
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

func (s *ErrandsServer) deletePipeline(c *gin.Context) {
	pipelineID := c.Param(idParam)
	pipeline, exists := s.getPipelineFromStore(pipelineID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"status": "not_found"})
		return
	}

	s.cascadeDeletePipeline(pipeline)
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *ErrandsServer) getPipeline(c *gin.Context) {
	pipelineID := c.Param(idParam)
	pipeline, exists := s.getPipelineFromStore(pipelineID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"status": "not_found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": pipeline,
	})
}

func (s *ErrandsServer) listPipelines(c *gin.Context) {
	statusFilter := c.Query(statusFilterQueryParam)
	filterFn := acceptAllPipelineFilter

	if statusFilter != "" {
		filterFn = statusPipelineFilter(schemas.Status(statusFilter))
	}

	pipelines := s.getFilteredPipelinesFromStore(filterFn)

	// We only want to return an overview of each pipeline in the list API.
	// Strip out errands and dependencies from the pipelines to make the response smaller.
	// Users can get pipeline details by ID
	for i := range pipelines {
		pipelines[i].Errands = nil
		pipelines[i].Dependencies = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"results": pipelines,
	})
}

func (s *ErrandsServer) cascadeDeletePipeline(pipeline schemas.Pipeline) {
	// Delete all of the errands in the pipeline
	for _, errand := range pipeline.Errands {
		s.deleteErrandByID(errand.ID)
	}

	// Delete the pipeline itself
	s.PipelineStore.Delete(pipeline.ID)
}

func (s *ErrandsServer) getPipelineFromStore(pipelineID string) (schemas.Pipeline, bool) {
	pipeline, exists := s.PipelineStore.Get(pipelineID)
	if !exists {
		return schemas.Pipeline{}, false
	}

	return pipeline.(schemas.Pipeline), true
}

func (s *ErrandsServer) getFilteredPipelinesFromStore(filterFn func(pipeline schemas.Pipeline) bool) []schemas.Pipeline {
	var results []schemas.Pipeline

	for _, pipelineItem := range s.PipelineStore.Items() {
		pipeline := pipelineItem.Object.(schemas.Pipeline)
		if filterFn(pipeline) {
			results = append(results, pipeline)
		}
	}

	return results
}

func statusPipelineFilter(status schemas.Status) func(pipeline schemas.Pipeline) bool {
	return func(pipeline schemas.Pipeline) bool {
		return pipeline.Status == status
	}
}

func acceptAllPipelineFilter(pipeline schemas.Pipeline) bool {
	return true
}

func (s *ErrandsServer) updateErrandInPipeline(errand *schemas.Errand) {
	pipeline, exists := s.getPipelineFromStore(errand.PipelineID)
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

	pipeline.RecalculateStatus()

	// If the pipeline just finished and is marked as deleteOnCompleted, delete it now
	if pipeline.Status == schemas.StatusCompleted && pipeline.DeleteOnCompleted {
		s.cascadeDeletePipeline(pipeline)
	} else {
		// Otherwise save the updated pipeline
		s.PipelineStore.SetDefault(pipeline.ID, pipeline)
	}
}
