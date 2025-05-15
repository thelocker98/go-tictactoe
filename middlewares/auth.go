package middlewares

import (
	"net/http"

	"example.com/tictactoe/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.Redirect(http.StatusFound, "/login?status=pleaselogin")
		context.Abort()
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.Redirect(http.StatusFound, "/login?status=pleaselogin")
		context.Abort()
		return
	}

	context.Set("userId", userId)

	context.Next()
}
