package main

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	cors "github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"

	binding "github.com/gin-gonic/gin/binding"
	store "github.com/polygon-io/errands-server/memorydb"
	schemas "github.com/polygon-io/errands-server/schemas"
	validator "gopkg.in/go-playground/validator.v8"
)

const (
	errandsDBPathSuffix   = "errands/"
	pipelinesDBPathSuffix = "pipelines/"
)

//easyjson:json
type Notification struct {
	Event  string         `json:"event"`
	Errand schemas.Errand `json:"errand,omitempty"`
}

type ErrandsServer struct {
	StorageDir       string
	Port             string
	ErrandStore      *store.MemoryStore
	PipelineStore    *store.MemoryStore
	Server           *http.Server
	API              *gin.Engine
	ErrandsRoutes    *gin.RouterGroup
	ErrandRoutes     *gin.RouterGroup
	Notifications    chan *Notification
	StreamClients    []*Client
	RegisterClient   chan *Client
	UnregisterClient chan *Client
	periodicSave     bool
}

func NewErrandsServer(cfg *Config) *ErrandsServer {
	obj := &ErrandsServer{
		StorageDir:       cfg.Storage,
		Port:             cfg.Port,
		StreamClients:    make([]*Client, 0),
		RegisterClient:   make(chan *Client, 10),
		UnregisterClient: make(chan *Client, 10),
		Notifications:    make(chan *Notification, 100),
		ErrandStore:      store.New(),
		PipelineStore:    store.New(),
		periodicSave:     true,
	}

	go obj.createAPI()
	go obj.broadcastLoop()

	if err := obj.ErrandStore.LoadFile(path.Join(cfg.Storage, errandsDBPathSuffix)); err != nil {
		log.WithError(err).Error("unable to load errands db")
		log.Warning("Could not load errand data from previous DB file.")
		log.Warning("If this is your first time running, this is normal.")
		log.Warning("If not please check the contents of your file: ", cfg.Storage)
	}

	if err := obj.ErrandStore.LoadFile(path.Join(cfg.Storage, pipelinesDBPathSuffix)); err != nil {
		log.WithError(err).Error("unable to load pipelines db")
		log.Warning("Could not load pipeline data from previous DB file.")
		log.Warning("If this is your first time running, this is normal.")
		log.Warning("If not please check the contents of your file: ", cfg.Storage)
	}

	go obj.periodicallySaveDB()

	go obj.periodicallyCheckTTLs()

	return obj
}

func (s *ErrandsServer) periodicallySaveDB() {
	for {
		time.Sleep(60 * time.Second)

		if !s.periodicSave {
			return
		}

		log.Info("Checkpoint saving DB to file...")
		s.saveDBs()
	}
}

func (s *ErrandsServer) periodicallyCheckTTLs() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	filter := func(errand *schemas.Errand) bool {
		if errand.Options.TTL <= 0 {
			return false
		}

		started := time.Unix(0, errand.Started * 1_000_000) // ms to ns
		ttlDuration := time.Duration(errand.Options.TTL)

		return time.Since(started) > ttlDuration
	}

	update := func(errand *schemas.Errand) error {
		if err := failErrand(errand, FailedRequest{Reason: "TTL Expired"}); err != nil {
			return fmt.Errorf("unable to fail errand: %s; %w", errand.ID, err)
		}

		s.AddNotification("failed", errand)
		return nil
	}

	for range t.C  {
		if err := s.UpdateErrandsByFilter(filter, update); err != nil {
			log.WithError(err).Error("Unable to fail errand(s)")
		}
	}
}

func (s *ErrandsServer) saveDBs() {
	if err := s.ErrandStore.SaveFile(path.Join(s.StorageDir, errandsDBPathSuffix)); err != nil {
		log.Error("----- Error checkpoint saving the errand DB to file -----")
		log.Error(err)
	}

	if err := s.PipelineStore.SaveFile(path.Join(s.StorageDir, pipelinesDBPathSuffix)); err != nil {
		log.Error("----- Error checkpoint saving the pipeline DB to file -----")
		log.Error(err)
	}
}

func (s *ErrandsServer) AddNotification(event string, errand *schemas.Errand) {
	obj := &Notification{
		Event:  event,
		Errand: *errand,
	}
	s.Notifications <- obj
}

func (s *ErrandsServer) broadcastLoop() {
	for {
		select {
		case client := <-s.RegisterClient:
			s.StreamClients = append(s.StreamClients, client)
		case client := <-s.UnregisterClient:
			for i, c := range s.StreamClients {
				if c == client {
					s.StreamClients = append(s.StreamClients[:i], s.StreamClients[i+1:]...)
				}
			}
		case not := <-s.Notifications:
			for _, client := range s.StreamClients {
				notificationCopy := &Notification{}
				*notificationCopy = *not
				client.Notifications <- notificationCopy
			}
		}
	}
}

func (s *ErrandsServer) kill() {
	s.killAPI()

	for _, client := range s.StreamClients {
		client.Gone()
	}

	s.killDB()
}

func (s *ErrandsServer) killAPI() {
	log.Println("Closing the HTTP Server")

	_ = s.Server.Close()
}

func (s *ErrandsServer) killDB() {
	log.Println("Closing the DB")
	s.saveDBs()
}

func UserStructLevelValidation(v *validator.Validate, structLevel *validator.StructLevel) {
	errand := structLevel.CurrentStruct.Interface().(schemas.Errand)
	if errand.Options.TTL < 5 && errand.Options.TTL != 0 {
		structLevel.ReportError(
			reflect.ValueOf(errand.Options.TTL), "ttl", "ttl", "must be positive, and more than 5",
		)
	}
}

func (s *ErrandsServer) createAPI() {
	s.API = gin.Default()

	CORSconfig := cors.DefaultConfig()
	CORSconfig.AllowCredentials = true
	CORSconfig.AllowOriginFunc = func(origin string) bool {
		return true
	}
	s.API.Use(cors.New(CORSconfig))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(UserStructLevelValidation, schemas.Errand{})
	}

	// Singular errand Routes:
	s.ErrandRoutes = s.API.Group("/v1/errand")
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

	// Errands Routes
	s.ErrandsRoutes = s.API.Group("/v1/errands")
	// Create a new errand:
	s.ErrandsRoutes.POST("/", s.createErrand)
	// Get all errands:
	s.ErrandsRoutes.GET("/", s.getAllErrands)
	// Notifications:
	s.ErrandsRoutes.GET("/notifications", s.errandNotifications)
	// Ready to process an errand:
	s.ErrandsRoutes.POST("/process/:type", s.processErrand)
	// Get all errands in a current type or state:
	s.ErrandsRoutes.GET("/list/:key/:val", s.getFilteredErrands)
	// Update all errands in this state:
	s.ErrandsRoutes.POST("/update/:key/:val", s.updateFilteredErrands)
	// Clear all finished errands older than.
	s.ErrandsRoutes.POST("/clear/:duration", s.clearErrands)

	// Singular pipeline routes
	pipelineRoutes := s.API.Group("/v1/pipeline")
	pipelineRoutes.POST("/", s.createPipeline)
	pipelineRoutes.GET("/:id", s.getPipeline)
	pipelineRoutes.DELETE("/:id", s.deletePipeline)

	// Plural pipeline routes
	pipelinesRoutes := s.API.Group("/v1/pipelines")
	pipelinesRoutes.GET("/", s.listPipelines)

	// Prometheus metrics
	s.API.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	s.Server = &http.Server{
		Addr:    s.Port,
		Handler: s.API,
	}

	log.Println("Starting server on port:", s.Port)

	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
