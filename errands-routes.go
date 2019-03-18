
package main


import (
	// "fmt"
	"log"
	"sort"
	"bytes"
	"net/http"
	"encoding/gob"
	gin "github.com/gin-gonic/gin"
	badger "github.com/dgraph-io/badger"
)






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
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": item,
	})
}



func ( s *ErrandsServer ) saveErrand( txn *badger.Txn, errand *Errand ) error {
	var bytesBuffer bytes.Buffer
	enc := gob.NewEncoder(&bytesBuffer)
	err := enc.Encode( errand ); if err != nil {
		return err
	}
	return txn.Set([]byte(errand.ID), bytesBuffer.Bytes())
}




func ( s *ErrandsServer ) getErrands( c *gin.Context ){
	errands, err := s.GetErrandsBy()
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






func ( s *ErrandsServer ) processErrand( c *gin.Context ){
	var procErrand *Errand
	errands := make([]*Errand, 0)
	hasFound := false
	typeFilter := c.Query("type")
	err := s.DB.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 50
		it := txn.NewIterator( opts )
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func( v []byte ) error {
				var bytesBuffer bytes.Buffer
				dec := gob.NewDecoder( &bytesBuffer )
				_, err := bytesBuffer.Write( v ); if err != nil {
					return err
				}
				var errand Errand
				err = dec.Decode( &errand ); if err != nil {
					return err
				}
				if errand.Status != "inactive" {
					return nil
				}
				if errand.Type != typeFilter {
					return nil
				}
				// Add to list of errands we could possibly process:
				errands = append( errands, &errand )
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
			procErrand.FailedReason = ""
			procErrand.Status = "active"
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
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"results": procErrand,
	})
}








func ( s *ErrandsServer ) GetErrandsBy() ( []Errand, error ) {

	errands := []Errand{}
	err := s.DB.View(func( txn *badger.Txn ) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 50
		it := txn.NewIterator( opts )
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func( v []byte ) error {
				var bytesBuffer bytes.Buffer
				dec := gob.NewDecoder( &bytesBuffer )
				_, err := bytesBuffer.Write( v ); if err != nil {
					return err
				}
				var errand Errand
				err = dec.Decode( &errand ); if err != nil {
					return err
				}
				errands = append( errands, errand )
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







