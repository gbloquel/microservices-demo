package main

import (
	"article-service/handler"
	"article-service/logger"
	"article-service/middlewares"
	"article-service/repository"
	"context"
	"github.com/Depado/ginprom"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func main() {
	loadConfig()
	logger.SetupLogging()
	logger.Logger.Infoln("-= Article Service =-")
	if os.Getenv("ENABLE_TRACING") == "1" {
		logger.Logger.Info("Tracing enabled.")
		cleanup := initTracer()
		defer cleanup(context.Background())
	} else {
		logger.Logger.Info("Tracing disabled.")
	}
	initDatabase()
	loadAPIServer()
}

func initTracer() func(context.Context) error {

	var secureOption otlptracegrpc.Option

	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		logger.Logger.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Logger.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return exporter.Shutdown
}

// loadConfig define the default values and loads the user configuration from config.yaml
func loadConfig() {
	viper.SetDefault("listen", ":8080")
	viper.SetDefault("mongodbUri", "mongodb://localhost:27017/alpha-articles")
	err := viper.BindEnv("mongodbUri", "MONGODB_URI")
	if err != nil {
		logger.Logger.Warnln(err)
	}

	viper.SetDefault("logLevel", "info")
	err = viper.BindEnv("logLevel", "LOG_LEVEL")
	if err != nil {
		logger.Logger.Warnln(err)
	}

	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("/etc/article-service/")
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Warnln(err)
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
	prometheus := ginprom.New(
		ginprom.Engine(Router),
		ginprom.Namespace("article_srv"),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
		ginprom.Ignore("/metrics", "/healthz"),
	)

	Router.Use(otelgin.Middleware(serviceName))
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
	Router.Use(prometheus.Instrument())
	Router.Use(requestid.New())
	Router.Use(gin.Recovery())

	Router.GET("/", handler.HealthZ)
	Router.GET("/healthz", handler.HealthZ)

	Router.GET("/article/", handler.GetArticle)
	Router.POST("/article/", handler.AddArticle)
	Router.DELETE("/article/:articleId/", handler.DeleteArticle)

	listenAddress := viper.GetString("listen")
	err := Router.Run(listenAddress)
	logger.Logger.Panicln(err)
}
