package main

import (
	"cart-service/handler"
	"cart-service/logger"
	"cart-service/middlewares"
	"cart-service/repository"
	"context"
	"github.com/Depado/ginprom"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	logger.Logger.Infoln("-= Cart service =-")
	if os.Getenv("ENABLE_TRACING") == "1" {
		logger.Logger.Infoln("Tracing enabled.")
		cleanup := initTracer()
		defer cleanup(context.Background())
	} else {
		logger.Logger.Infoln("Tracing disabled.")
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
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
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
	prometheus := ginprom.New(
		ginprom.Engine(Router),
		ginprom.Namespace("cart_srv"),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
		ginprom.Ignore("/metrics", "/healthz"),
	)

	Router.Use(otelgin.Middleware(serviceName))
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
		prometheus.Instrument(),
		otelgin.Middleware(serviceName),
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
