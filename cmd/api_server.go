package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Depado/ginprom"
	"github.com/aurowora/compress"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "github.com/map-services/company-data-api/docs"
	"github.com/map-services/company-data-api/internal"
	"github.com/map-services/company-data-api/internal/middleware"
	"github.com/map-services/company-data-api/pkg/logger"
	repo "github.com/map-services/company-data-api/internal/repositories"
	"github.com/map-services/company-data-api/internal/routes"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rm-hull/godx"
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
	logger.SetupLogger()

	godx.GitVersion()
	godx.EnvironmentVars()
	godx.UserInfo()

	db, err := internal.Connect(dbPath)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("error closing database", "error", err)
		}
	}()

	repo, err := repo.NewSqliteDbRepository(db)
	if err != nil {
		slog.Error("failed to initialize repository", "error", err)
		os.Exit(1)
	}

	r := gin.New()

	prometheus := ginprom.New(
		ginprom.Engine(r),
		ginprom.Path("/metrics"),
		ginprom.Ignore("/healthz"),
	)

	r.Use(
		gin.Recovery(),
		middleware.RequestLogger(),
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
		slog.Warn("pprof endpoints are enabled and exposed. Do not run with this flag in production.")
		pprof.Register(r)
	}

	err = healthcheck.New(r, hc_config.DefaultConfig(), []checks.Check{
		checks.SqlCheck{Sql: db},
	})
	if err != nil {
		slog.Error("failed to initialize healthcheck", "error", err)
		os.Exit(1)
	}

	v1 := r.Group("/v1/company-data")
	v1.GET("/search", routes.Search(repo))
	v1.GET("/search/by-postcode", routes.GroupByPostcode(repo))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", port)
	slog.Info("Starting HTTP API Server", "port", port)
	if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP API Server failed to start", "port", port, "error", err)
		os.Exit(1)
	}
}
