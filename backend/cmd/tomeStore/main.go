package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/diconico07/TomeTracker/pkg/ent"
	"github.com/diconico07/TomeTracker/pkg/tomeStore"
	"github.com/gin-gonic/gin"

	"github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	sqlCfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Net:                  "tcp",
		Addr:                 dbHost + ":" + dbPort,
		DBName:               dbName,
		ParseTime:            true,
		Loc:                  time.UTC,
		AllowNativePasswords: true,
	}

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	client, err := ent.Open("mysql", sqlCfg.FormatDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to mysql")
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed creating schema resources")
	}

	svc := tomeStore.NewService(client)
	bookSvc := tomeStore.NewBookService(client)
	router := gin.Default()

	api := router.Group("/api")

	tomeStore.DefineRoutes(api, svc, bookSvc)

	router.Static("/assets", "./dist/assets")
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/assets") {
			c.Status(http.StatusNotFound)
			return
		}
		c.File("./dist/index.html")
	})

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	router.Run(listenAddr)
}
