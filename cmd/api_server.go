package cmd

import (
	_ "company-data-api/docs"
	"company-data-api/internal"
	"company-data-api/repositories"
	"time"

	"fmt"
	"log"

	"github.com/Depado/ginprom"
	"github.com/aurowora/compress"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	hc_config "github.com/tavsec/gin-healthcheck/config"
	cachecontrol "go.eigsys.de/gin-cachecontrol/v2"
)

// @title Company Data API
// @version 1.0
// @description A fast REST API for querying UK company data by geographic bounding box, built with Go, SQLite, and Gin. It imports official datasets from Companies House and Ordnance Survey CodePoint Open, providing spatial search capabilities for company records.
// @BasePath /v1/company-data
func ApiServer(dbPath string, port int, debug bool) {

	db, err := internal.Connect(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	repo, err := repositories.NewSqliteDbRepository(db)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}

	r := gin.New()

	prometheus := ginprom.New(
		ginprom.Engine(r),
		ginprom.Path("/metrics"),
		ginprom.Ignore("/healthz"),
	)

	r.Use(
		gin.Recovery(),
		gin.LoggerWithWriter(gin.DefaultWriter, "/healthz", "/metrics"),
		prometheus.Instrument(),
		compress.Compress(),
		cachecontrol.New(cachecontrol.Config{
			// Set the max-age to 28 days: Companies House data is generally updated monthly
			MaxAge:    cachecontrol.Duration(28 * 24 * time.Hour),
			Immutable: true,
			Public:    true,
		}),
		cors.Default(),
	)

	if debug {
		log.Println("WARNING: pprof endpoints are enabled and exposed. Do not run with this flag in production.")
		pprof.Register(r)
	}

	err = healthcheck.New(r, hc_config.DefaultConfig(), []checks.Check{
		checks.SqlCheck{Sql: db},
	})
	if err != nil {
		log.Fatalf("failed to initialize healthcheck: %v", err)
	}

	v1 := r.Group("/v1/company-data")
	v1.GET("/search", internal.Search(repo))
	v1.GET("/search/by-postcode", internal.GroupByPostcode(repo))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting HTTP API Server on port %d...", port)
	err = r.Run(addr)
	log.Fatalf("HTTP API Server failed to start on port %d: %v", port, err)
}
