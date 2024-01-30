package main

import (
	"article-service/handler"
	"article-service/logger"
	"article-service/middlewares"
	"article-service/repository"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	loadConfig()
	logger.SetupLogging()
	logger.Logger.Infoln("-= Article Service =-")
	initDatabase()
	loadAPIServer()
}

// loadConfig define the default values and loads the user configuration from config.yaml
func loadConfig() {
	viper.SetDefault("listen", ":8080")
	viper.SetDefault("mongodbUri", "mongodb://localhost:27017/alpha-articles")
	err := viper.BindEnv("mongodbUri", "MONGODB_URI")
	if err != nil {
		log.Warnln(err)
	}

	viper.SetDefault("logLevel", "info")
	err = viper.BindEnv("logLevel", "LOG_LEVEL")
	if err != nil {
		log.Warnln(err)
	}

	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("/etc/article-service/")
	if err := viper.ReadInConfig(); err != nil {
		log.Warnln(err)
	}
}

// initDatabase initialize the database connection
func initDatabase() {
	mongodbURI := viper.GetString("mongodbURI")
	if err := repository.Initialize(mongodbURI); err != nil {
		logger.Logger.Errorf("Failed to connect to %s", mongodbURI)
		logger.Logger.Panicln(err)
	}
	logger.Logger.Infof("Connected to %s", mongodbURI)
}

// loadAPIServer initialize the API server with a cors middleware and define routes to be served.
// This function is blocking: it will wait until the server returns an error
func loadAPIServer() {
	Router := gin.New()
	Router.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true //return origin == "http://localhost:3001"
		},
		MaxAge: 12 * time.Hour,
	}))

	Router.Use(middlewares.LoggingMiddleware(logger.Logger, "/", "/healthz"))
	Router.Use(requestid.New())
	Router.Use(gin.Recovery())

	Router.GET("/", handler.HealthZ)
	Router.GET("/healthz", handler.HealthZ)

	Router.GET("/article/", handler.GetArticle)
	Router.POST("/article/", handler.AddArticle)
	Router.DELETE("/article/:articleId/", handler.DeleteArticle)

	listenAddress := viper.GetString("listen")
	err := Router.Run(listenAddress)
	log.Panicln(err)
}
