package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/polygon-io/errands-server/schemas"
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

	// TODO: actually save it
}
