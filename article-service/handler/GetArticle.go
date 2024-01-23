package handler

import (
	"article-service/logger"
	"article-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetArticle is a handler function for retrieving articles.
// It retrieves a list of articles from the repository using GetArticles function.
// If an error occurs during the retrieval process, it returns an Internal Server Error response.
// Otherwise, it returns a JSON response containing the retrieved articles.
func GetArticle(ctx *gin.Context) {
	// Retrieve a list of articles from the repository
	logger.Logger.Infof("GetArticle - Retrieve list of article from database")

	articles, err := repository.GetArticles(ctx, nil)
	if err != nil {
		// Handle repository error, log the error, and return Internal Server Error response
		logger.Logger.Errorf(errorMsgFormat, "GetArticle", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Infof(successfulMsgFormat, "GetArticle", articles)

	ctx.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}
