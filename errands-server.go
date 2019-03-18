
package main


import (
	"log"
	"reflect"
	"net/http"
	gin "github.com/gin-gonic/gin"
	cors "github.com/gin-contrib/cors"
	gzip "github.com/gin-contrib/gzip"
	badger "github.com/dgraph-io/badger"
	binding "github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
)








type ErrandsServer struct {
	DB 					*badger.DB
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
	s.DB.Close()
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
}


func ( s *ErrandsServer) createAPI(){

	s.API = gin.Default()

	CORSconfig := cors.DefaultConfig()
	CORSconfig.AllowOriginFunc = func( origin string ) bool {
		// fmt.Println("Connection from", origin)
		return true
	}
	s.API.Use(cors.New(CORSconfig))
	s.API.Use(gzip.Gzip(gzip.DefaultCompression))

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
		// Ready to process an errand:
		s.ErrandsRoutes.POST("/process", s.processErrand)
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




