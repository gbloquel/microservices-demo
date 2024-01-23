package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware instance a Logger middleware with the specified logrus instance.
func LoggingMiddleware(logger *logrus.Logger, notLoggedPaths ...string) gin.HandlerFunc {

	var skippedPaths map[string]bool
	if length := len(notLoggedPaths); length > 0 {
		skippedPaths = make(map[string]bool, length)
		for _, currentPath := range notLoggedPaths {
			skippedPaths[currentPath] = true
		}
	}

	return func(ctx *gin.Context) {
		startTime := time.Now()
		path := ctx.Request.URL.Path

		// Processing request
		ctx.Next()

		if _, ok := skippedPaths[path]; !ok {
			endTime := time.Now()
			latencyTime := endTime.Sub(startTime)
			reqMethod := ctx.Request.Method
			requestURI := ctx.Request.RequestURI
			statusCode := ctx.Writer.Status()
			clientIP := ctx.ClientIP()
			clientUserAgent := ctx.Request.UserAgent()
			referer := ctx.Request.Referer()

			msg := logger.WithFields(logrus.Fields{
				"method":     reqMethod,
				"path":       requestURI,
				"statusCode": statusCode,
				"latency":    latencyTime,
				"clientIP":   clientIP,
				"userAgent":  clientUserAgent,
				"referer":    referer,
				"requestId":  requestid.Get(ctx),
			})

			switch {
			case statusCode >= http.StatusInternalServerError:
				msg.Error("Error when processing request")
			case statusCode >= http.StatusBadRequest:
				msg.Warn("Invalid request")
			default:
				msg.Info("processed Request with success")
			}
		}
	}
}
