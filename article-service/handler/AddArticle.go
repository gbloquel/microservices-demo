package handler

import (
	"article-service/logger"
	"article-service/model"
	"article-service/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddArticle is a handler function for adding an article.
// It expects a JSON payload in the request body, representing an AddArticleRequest.
// If the request body is invalid or does not match the expected format, it returns a Bad Request response.
// If the AddArticleRequest fails validation, it returns a Bad Request response with validation error details.
// Otherwise, it adds the article to the repository, and returns a success response.
func AddArticle(ctx *gin.Context) {
	// Parse the JSON request body into an AddArticleRequest struct
	var addArticleRequest AddArticleRequest
	if err := ctx.ShouldBindJSON(&addArticleRequest); err != nil {
		logger.Logger.Warnf(errorMsgFormat, "AddArticle", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Infof(
		"AddArticle - Adding % (%s) to database",
		addArticleRequest.ArticleName,
		addArticleRequest.ArticleDescription)

	// Validate the AddArticleRequest
	if err := addArticleRequest.validate(); err != nil {
		logger.Logger.Warnf(errorMsgFormat, "AddArticle", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the validated article to the repository
	err := repository.AddArticle(ctx, addArticleRequest.toArticle())
	if err != nil {
		logger.Logger.Errorf(errorMsgFormat, "AddArticle", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Infof(successfulMsgFormat, "AddArticle", addArticleRequest.ArticleName)
}

// AddArticleRequest represents the structure of a JSON request
// expected when adding an article. It includes fields for the
// article's name and description. The struct tags define the
// corresponding JSON key names for serialization and deserialization.
type AddArticleRequest struct {

	// ArticleName is the name of the article.
	// It is expected to be provided in the JSON request as "name".
	ArticleName string `json:"name"`

	// ArticleDescription is the description of the article.
	// It is expected to be provided in the JSON request as "description".
	ArticleDescription string `json:"description"`
}

func (r AddArticleRequest) validate() error {
	if r.ArticleName == "" {
		return fmt.Errorf("articleName is required")
	}
	if r.ArticleDescription == "" {
		return fmt.Errorf("articleDescription is required")
	}
	return nil
}

func (r AddArticleRequest) toArticle() *model.Article {
	return &model.Article{
		Title:       r.ArticleName,
		Description: r.ArticleDescription,
	}
}
