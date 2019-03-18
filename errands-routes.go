
package main


import (
	"log"
	"fmt"
	"time"
	"bytes"
	"net/http"
	"encoding/gob"
	uuid "github.com/google/uuid"
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

	err := s.DB.Update(func( txn *badger.Txn ) error {
		itemUUID := uuid.New()
		itemUUIDText, err := itemUUID.MarshalText(); if err != nil {
			return err
		}
		// Set some additional generated attributes:
		item.ID = string(itemUUIDText)
		item.Created = ( time.Now().UnixNano() / 1000000 )
		item.Status = "inactive"

		var bytesBuffer bytes.Buffer
		enc := gob.NewEncoder(&bytesBuffer)
		err = enc.Encode( item ); if err != nil {
			return err
		}
		return txn.Set(itemUUIDText, bytesBuffer.Bytes())
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error!",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})

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
				println("written")
				var errand Errand
				err = dec.Decode( &errand ); if err != nil {
					return err
				}
				println("parsed")
				errands = append( errands, errand )
				fmt.Println("errands:", errands, errand)
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







