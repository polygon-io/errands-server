
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
	StorageDir 			string
	Port 				string
	DB 					*badger.DB
	Server 				*http.Server
	API 				*gin.Engine
	ErrandsRoutes 		*gin.RouterGroup
	ErrandRoutes 		*gin.RouterGroup
}



func NewErrandsServer( cfg *Config ) *ErrandsServer {
	obj := &ErrandsServer{
		StorageDir: cfg.Storage,
		Port: cfg.Port,
	}
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
	opts.Dir = s.StorageDir
	opts.ValueDir = s.StorageDir
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
		s.ErrandRoutes.DELETE("/:id", s.deleteErrand)
		// Update an errand by id:
		s.ErrandRoutes.PUT("/:id", s.updateErrand)
		s.ErrandRoutes.PUT("/:id/failed", s.failedErrand)
		s.ErrandRoutes.PUT("/:id/completed", s.completeErrand)
		s.ErrandRoutes.POST("/:id/log", s.logToErrand)
		s.ErrandRoutes.POST("/:id/retry", s.retryErrand)
	}


	// Errands Routes
	s.ErrandsRoutes = s.API.Group("/v1/errands")
	{
		// Create a new errand:
		s.ErrandsRoutes.POST("/", s.createErrand)
		// Get all errands:
		s.ErrandsRoutes.GET("/", s.getAllErrands)
		// Ready to process an errand:
		s.ErrandsRoutes.POST("/process", s.processErrand)
		// Get all errands in a current type or state:
		s.ErrandsRoutes.GET("/list/:key/:val", s.getFilteredErrands)
		// Update all errands in this state:
		s.ErrandsRoutes.POST("/update/:key/:val", s.updateFilteredErrands)
	}

	s.Server = &http.Server{
		Addr: 		s.Port,
		Handler: 	s.API,
	}

	log.Println("Starting server on port:", s.Port)
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

}




