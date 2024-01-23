package handler

import (
	"article-service/logger"
	"article-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteArticle is a handler function for deleting an article.
// It retrieves the article ID from the request parameters and attempts to delete the corresponding article.
// If an error occurs during the deletion process, it logs the error and returns an Internal Server Error response.
// Otherwise, it returns an Accepted response.
func DeleteArticle(ctx *gin.Context) {
	// Retrieve the article ID from the request parameters
	articleID := ctx.Param("articleId")

	logger.Logger.Infof(
		"DeleteArticle - Deleting article with ID %s",
		articleID)

	// Attempt to delete the article from the repository
	if err := repository.DeleteArticle(ctx, articleID); err != nil {
		// Handle repository error, log the error, and return Internal Server Error response
		logger.Logger.Errorf(errorMsgFormat, "DeleteArticle", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":     err.Error(),
			"articleID": articleID,
		})
		return
	}

	logger.Logger.Infof(successfulMsgFormat, "DeleteArticle", articleID)
	ctx.Status(http.StatusAccepted)
}
