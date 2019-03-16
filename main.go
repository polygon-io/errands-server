
package main

import (
	"net/http"
	"reflect"
	// badger "github.com/dgraph-io/badger"
	gin "github.com/gin-gonic/gin"
	binding "github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
)


var server *ErrandsServer
func main(){

	server = NewErrandsServer()

}





type ErrandsServer struct {

	API 				*gin.Engine
	ErrandsRoutes 		*gin.RouterGroup
	ErrandRoutes 		*gin.RouterGroup
}



func NewErrandsServer() *ErrandsServer {
	obj := &ErrandsServer{}


	obj.createAPI()

	return obj
}


/*

Name
Type
TTL
Retries
Priority
Data

 */

type Errand struct {
	Name 			string 		`json:"name" binding:"required"`
	Type 			string 		`json:"type" binding:"required"`
	TTL 			int 		`json:"ttl"`
	Retries 		int 		`json:"retries"`
	Priority 		int 		`json:"priority"`
	Data 			*gin.H 		`json:"data"`
}







func ( s *ErrandsServer ) createErrand( c *gin.Context ){
	println("creating errand")
	var item Errand
	if err := c.ShouldBindJSON(&item); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Errand validation successful."})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Errand validation failed!",
			"error":   err.Error(),
		})
	}
}



func UserStructLevelValidation(v *validator.Validate, structLevel *validator.StructLevel) {
	errand := structLevel.CurrentStruct.Interface().(Errand)

	if errand.TTL < 5 && errand.TTL != 0 {
		structLevel.ReportError(
			reflect.ValueOf(errand.TTL), "ttl", "ttl", "must be positive, and more than 5",
		)
	}

	// plus can to more, even with different tag than "fnameorlname"
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
		s.ErrandsRoutes.GET("/", s.createErrand)
		// Get all errands in a current type or state:
		s.ErrandsRoutes.GET("/list/:key/:val", s.createErrand)
		// Update all errands in this state:
		s.ErrandsRoutes.POST("/update/:type", s.createErrand)
	}

	s.API.Run(":5555")

}

