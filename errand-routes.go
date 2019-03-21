
package main


import (
	// "fmt"
	"errors"
	"net/http"
	gin "github.com/gin-gonic/gin"
	badger "github.com/dgraph-io/badger"
)




type UpdateRequest struct {
	Progress 			float64 	`json:"progress"`
	Logs 				[]string 	`json:"logs"`
}
func ( s *ErrandsServer ) updateErrand( c *gin.Context ){
	var updatedErrand *Errand
	var updateReq UpdateRequest
	if err := c.ShouldBind(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error": err.Error(),
		})
		return
	}
	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func( errand *Errand ) error {
		if errand.Status != "active" {
			return errors.New("Errand must be in active state to update progress")
		}
		// Update this errand attributes:
		if updateReq.Progress != 0 {
			if updateReq.Progress < 0 || updateReq.Progress >= 101 {
				return errors.New("Progress must be between 0 - 100")
			}
			errand.Progress = float64( updateReq.Progress )
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}
	s.AddNotification( "updated", updatedErrand )
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": updatedErrand,
	})
}






type FailedRequest struct {
	Reason 			string 	`json:"reason" binding:"required"`
}
func ( s *ErrandsServer ) failedErrand( c *gin.Context ){
	var updatedErrand *Errand
	var failedReq FailedRequest
	if err := c.ShouldBind(&failedReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error": err.Error(),
		})
		return
	}
	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func( errand *Errand ) error {
		// if errand.Status != "active" {
		// 	return errors.New("Errand must be in active state to fail")
		// }
		// Update this errand attributes:
		if err := errand.addToLogs("ERROR", failedReq.Reason); err != nil {
			return err
		}
		errand.Failed = getTimestamp()
		errand.Status = "failed"
		errand.Progress = 0
		if errand.Options.Retries > 0 {
			// If we should retry this errand:
			if errand.Attempts <= errand.Options.Retries {
				errand.Status = "inactive"
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}
	s.AddNotification( "failed", updatedErrand )
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": updatedErrand,
	})
}




type CompletedRequest struct {
	Results 			*gin.H 	`json:"results"`
}
func ( s *ErrandsServer ) completeErrand( c *gin.Context ){
	var updatedErrand *Errand
	var compReq CompletedRequest
	if err := c.ShouldBind(&compReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error": err.Error(),
		})
		return
	}
	shouldDelete := false
	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func( errand *Errand ) error {
		// if errand.Status != "active" {
		// 	return errors.New("Errand must be in active state to complete")
		// }
		// Update this errand attributes:
		if err := errand.addToLogs("INFO", "Completed!"); err != nil {
			return err
		}
		errand.Completed = getTimestamp()
		errand.Status = "completed"
		errand.Progress = 100
		errand.Results = compReq.Results
		// If we should delete this errand upon completion:
		if errand.Options.DeleteOnCompleted == true {
			shouldDelete = true
		}
		return nil
	})
	if err == nil && shouldDelete == true && updatedErrand.ID != "" {
		err = s.deleteErrandByID( updatedErrand.ID )
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}
	s.AddNotification( "completed", updatedErrand )
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": updatedErrand,
	})
}






func ( s *ErrandsServer ) retryErrand( c *gin.Context ){
	var updatedErrand *Errand
	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func( errand *Errand ) error {
		if errand.Status == "inactive" {
			return errors.New("Cannot retry errand which is in inactive state")
		}
		// Update this errand attributes:
		if err := errand.addToLogs("INFO", "Retrying!"); err != nil {
			return err
		}
		errand.Status = "inactive"
		errand.Progress = 0
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}
	s.AddNotification( "retry", updatedErrand )
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": updatedErrand,
	})
}




func ( s *ErrandsServer ) logToErrand( c *gin.Context ){
	var logReq Log
	if err := c.ShouldBind(&logReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error": err.Error(),
		})
		return
	}
	updatedErrand, err := s.UpdateErrandByID(c.Param("id"), func( errand *Errand ) error {
		if errand.Status != "active" {
			return errors.New("Errand must be in active state to log to")
		}
		// Update this errand attributes:
		return errand.addToLogs(logReq.Severity, logReq.Message)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": updatedErrand,
	})
}





func ( s *ErrandsServer ) deleteErrand( c *gin.Context ){
	err := s.deleteErrandByID( c.Param("id") )
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}








func ( s *ErrandsServer ) deleteErrandByID( id string ) error {
	return s.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte( id ))
	})
}



/*
	Pass in a function which will be called allowing you to update the errand. 
	If no error is returned, the errand will be saved in the DB with the new 
	attributes.
 */
func ( s *ErrandsServer ) UpdateErrandByID( id string, fn func( *Errand ) error ) ( *Errand, error ) {
	var updatedErrand *Errand
	err := s.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte( id )); if err != nil {
			return err
		}
		err = item.Value(func(v []byte) error {
			errand := &Errand{}
			err := errand.UnmarshalJSON( v ); if err != nil {
				return err
			}
			err = fn( errand ); if err != nil {
				return err
			}
			updatedErrand = errand
			err = s.saveErrand( txn, errand ); if err != nil {
				return err
			}
			return nil
		})
		return err
	})
	return updatedErrand, err
}

