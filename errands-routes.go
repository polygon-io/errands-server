
package main


import (
	"io"
	// "fmt"
	"log"
	"sort"
	// "time"
	"errors"
	"net/http"
	"encoding/json"
	gin "github.com/gin-gonic/gin"
	badger "github.com/dgraph-io/badger"
)





func ( s *ErrandsServer ) errandNotifications( c *gin.Context ){
	client, err := s.NewClient( c ); if err != nil {
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
			case <- clientGone:
				client.Gone()
				return false
			case t, ok := <-client.Notifications:
				if ok {
					// If we are subscribed to this event type:
					if contains(client.EventSubs, t.Event) || client.EventSubs[0] == "*" {
						jsonData, _ := json.Marshal( t )
						client.Gin.SSEvent("message", string( jsonData ))
						w.Flush()
					}
					return true
				}
				return false
			}
		}
		return false
	})
}










func ( s *ErrandsServer ) createErrand( c *gin.Context ){
	log.Println("creating errand")
	var item Errand
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Errand validation failed!",
			"error":   err.Error(),
		})
		return	
	}
	item.setDefaults()
	err := s.DB.Update(func( txn *badger.Txn ) error {
		return s.saveErrand( txn, &item )
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}
	s.AddNotification( "created", &item )
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": item,
	})
}



func ( s *ErrandsServer ) saveErrand( txn *badger.Txn, errand *Errand ) error {
	if !contains(ErrandStatuses, errand.Status) {
		return errors.New("Invalid errand status state")
	}
	bytes, err := errand.MarshalJSON(); if err != nil {
		return err
	}
	return txn.Set([]byte(errand.ID), bytes)
}




func ( s *ErrandsServer ) getAllErrands( c *gin.Context ){
	errands, err := s.GetErrandsBy(func( errand *Errand ) bool {
		return true
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
		"results": errands,
	})
}



func ( s *ErrandsServer ) getFilteredErrands( c *gin.Context ){
	key := c.Param("key")
	value := c.Param("val")
	errands, err := s.GetErrandsBy(func( errand *Errand ) bool {
		switch key {
		case "status":
			return ( errand.Status == value )
		case "type":
			return ( errand.Type == value )
		default:
			return false
		}
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
		"results": errands,
	})
}




type filteredUpdateReq struct {
	Status 			string 		`json:"status"`
	Delete 			bool 		`json:"delete"`
}
func ( s *ErrandsServer ) updateFilteredErrands( c *gin.Context ){
	key := c.Param("key")
	value := c.Param("val")
	var updateReq filteredUpdateReq
	if err := c.ShouldBind(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Parameters",
			"error": err.Error(),
		})
		return
	}
	errands, err := s.GetErrandsBy(func( errand *Errand ) bool {
		switch key {
		case "status":
			return ( errand.Status == value )
		case "type":
			return ( errand.Type == value )
		default:
			return false
		}
	})
	if err == nil {
		for _, errand := range errands {
			if updateReq.Delete == true {
				err = s.deleteErrandByID( errand.ID ); if err != nil {
					break
				}
			}else {
				if updateReq.Status != "" {
					_, err = s.UpdateErrandByID( errand.ID, func( e *Errand ) error {
						e.Status = updateReq.Status
						return nil
					})
					if err != nil {
						break
					}
				}
			}
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"count": len( errands ),
	})
}






func ( s *ErrandsServer ) processErrand( c *gin.Context ){
	var procErrand *Errand
	errands := make([]*Errand, 0)
	hasFound := false
	typeFilter := c.Param("type")
	err := s.DB.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 50
		it := txn.NewIterator( opts )
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func( v []byte ) error {

				errand := &Errand{}
				err := errand.UnmarshalJSON( v ); if err != nil {
					return err
				}

				if errand.Status != "inactive" {
					return nil
				}
				if errand.Type != typeFilter {
					return nil
				}
				// Add to list of errands we could possibly process:
				errands = append( errands, errand )
				return nil
			})
			if err != nil {
				return err
			}
		}

		// Of the possible errands to process, sort them by date & priority:
		if len( errands ) > 0 {
			sort.SliceStable(errands, func(i, j int) bool {
				return errands[i].Created < errands[j].Created 
			})
			sort.SliceStable(errands, func(i, j int) bool {
				return errands[i].Options.Priority > errands[j].Options.Priority 
			})
			procErrand = errands[0]
			// We are processing this errand:
			procErrand.Started = getTimestamp()
			procErrand.Attempts += 1
			procErrand.Status = "active"
			procErrand.Progress = 0.0
			if err := procErrand.addToLogs("INFO", "Started!"); err != nil {
				return err
			}
			err := s.saveErrand( txn, procErrand ); if err != nil {
				return err
			}
			hasFound = true
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
	if !hasFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No jobs",
		})
		return
	}
	s.AddNotification( "processing", procErrand )
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": procErrand,
	})
}








func ( s *ErrandsServer ) GetErrandsBy( fn func ( *Errand ) bool ) ( []*Errand, error ) {
	errands := make([]*Errand, 0)
	err := s.DB.View(func( txn *badger.Txn ) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 50
		it := txn.NewIterator( opts )
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func( v []byte ) error {
				errand := &Errand{}
				err := errand.UnmarshalJSON( v ); if err != nil {
					return err
				}
				if( fn( errand ) ){
					errands = append( errands, errand )
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return errands, err
}







