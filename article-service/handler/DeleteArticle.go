package handler

import (
	"article-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func DeleteArticle(ctx *gin.Context) {
	articleId := ctx.Param("articleId")

	if err := repository.DeleteArticle(ctx.Request.Context(), articleId); err != nil {
		log.Warnf("DeleteArticle Error: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":     err.Error(),
			"articleId": articleId,
		})
		return
	}

	ctx.Status(http.StatusAccepted)
}
