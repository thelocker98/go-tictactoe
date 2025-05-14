package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func loadGamePage(context *gin.Context) {
	context.HTML(http.StatusOK, "game.html", gin.H{"title": "Best Move"})
}
