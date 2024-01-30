package main

import (
	"cart-service/handler"
	"cart-service/logger"
	"cart-service/middlewares"
	"cart-service/repository"
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
	logger.Logger.Infoln("-= Cart service =-")
	initDatabase()
	loadAPIServer()
}

// loadConfig define the default values and loads the user configuration from config.yaml
func loadConfig() {
	viper.SetDefault("listen", ":8081")
	viper.SetDefault("redisUri", "redis://localhost:6379")
	err := viper.BindEnv("redisUri", "REDIS_URI")
	if err != nil {
		logger.Logger.Warnln(err)
	}

	viper.SetDefault("logLevel", "info")
	err = viper.BindEnv("logLevel", "LOG_LEVEL")
	if err != nil {
		log.Warnln(err)
	}

	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("/etc/article-cart/")
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Warnln(err)
	}
}

// initDatabase initialize the database connection
func initDatabase() {
	redisURI := viper.GetString("redisURI")
	if err := repository.Initialize(redisURI); err != nil {
		logger.Logger.Errorf("Failed to connect to %s", redisURI)
		logger.Logger.Panicln(err)
	}
	logger.Logger.Infof("Connected to %s", redisURI)
}

// loadAPIServer initialize the API server with a cors middleware and define routes to be served.
// This function is blocking: it will wait until the server returns an error
func loadAPIServer() {
	Router := gin.New()
	Router.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true //return origin == "http://localhost:3001"
		},
		MaxAge: 12 * time.Hour,
	}))

	Router.Use(
		middlewares.LoggingMiddleware(logger.Logger, "/", "/healthz"),
		requestid.New(),
		gin.Recovery(),
	)

	Router.GET("/", handler.HealthZ)
	Router.GET("/healthz", handler.HealthZ)
	Router.GET("/cart/:cartId/", handler.GetCart)
	Router.PUT("/cart/:cartId/", handler.UpdateCart)
	Router.DELETE("/cart/:cartId/", handler.DeleteCart)

	listenAddress := viper.GetString("listen")
	err := Router.Run(listenAddress)
	log.Panicln(err)
}
