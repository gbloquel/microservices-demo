package handler

import (
	"cart-service/reporitory"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteCart(ctx *gin.Context) {
	cartId := ctx.Param("cartId")

	err := reporitory.DeleteCart(ctx.Request.Context(), cartId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusAccepted)
}
