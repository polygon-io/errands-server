
package main


import (
	"log"
	"fmt"
	"time"
	"bytes"
	"reflect"
	"net/http"
	"encoding/gob"
	uuid "github.com/google/uuid"
	gin "github.com/gin-gonic/gin"
	badger "github.com/dgraph-io/badger"
	binding "github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
)








type ErrandsServer struct {
	DB 					*badger.DB
	DBSeq 				*badger.Sequence
	Server 				*http.Server
	API 				*gin.Engine
	ErrandsRoutes 		*gin.RouterGroup
	ErrandRoutes 		*gin.RouterGroup
}



func NewErrandsServer() *ErrandsServer {
	obj := &ErrandsServer{}
	go obj.createAPI()
	obj.createDB()
	return obj
}


func ( s *ErrandsServer ) kill(){
	s.killAPI()
	s.killDB()
}
func ( s *ErrandsServer) killAPI(){
	log.Println("Closing the HTTP Server")
	s.Server.Close()
}
func ( s *ErrandsServer ) killDB(){
	log.Println("Closing the DB")
	// s.DBSeq.Release()
	s.DB.Close()
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












func UserStructLevelValidation(v *validator.Validate, structLevel *validator.StructLevel) {
	errand := structLevel.CurrentStruct.Interface().(Errand)
	if errand.Options.TTL < 5 && errand.Options.TTL != 0 {
		structLevel.ReportError(
			reflect.ValueOf(errand.Options.TTL), "ttl", "ttl", "must be positive, and more than 5",
		)
	}
}



func ( s *ErrandsServer ) createDB(){
	opts := badger.DefaultOptions
	opts.Dir = "./badger"
	opts.ValueDir = "./badger"
	var err error
	s.DB, err = badger.Open( opts ); if err != nil {
		log.Fatal( err )
	}
	// s.DBSeq, err = s.DB.GetSequence([]byte("errands"), 1000); if err != nil {
	// 	log.Fatal( err )
	// }
}


func ( s *ErrandsServer) createAPI(){

	s.API = gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(UserStructLevelValidation, Errand{})
	}

	// Singular errand Routes:
	s.ErrandRoutes = s.API.Group("/v1/errand")
	{
		// Get an errand by id:
		s.ErrandRoutes.GET("/:id", s.createErrand)
		// Delete an errand by id:
		s.ErrandRoutes.DELETE("/:id", s.createErrand)
		// Update an errand by id:
		s.ErrandRoutes.PATCH("/:id", s.createErrand)
	}


	// Errands Routes
	s.ErrandsRoutes = s.API.Group("/v1/errands")
	{
		// Create a new errand:
		s.ErrandsRoutes.POST("/", s.createErrand)
		// Get all errands:
		s.ErrandsRoutes.GET("/", s.getErrands)
		// Get all errands in a current type or state:
		s.ErrandsRoutes.GET("/list/:key/:val", s.createErrand)
		// Update all errands in this state:
		s.ErrandsRoutes.POST("/update/:type", s.createErrand)
	}

	s.Server = &http.Server{
		Addr: 		":5555",
		Handler: 	s.API,
	}

	log.Println("Starting server on port: 5555")
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

}




