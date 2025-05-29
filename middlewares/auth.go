package middlewares

import (
	"net/http"

	"example.com/tictactoe/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {

	cookie, err := context.Cookie("session")
	if err != nil {
		// No Cookie Found so try header
		token := context.Request.Header.Get("Authorization")

		// check if it is empty
		if token == "" {
			context.Redirect(http.StatusFound, "/login?status=pleaselogin")
			context.Abort()
			return
		}

		cookie = token
	}

	userId, err := utils.VerifyToken(cookie)
	if err != nil {
		context.Redirect(http.StatusFound, "/login?status=pleaselogin")
		context.Abort()
		return
	}

	context.Set("userId", userId)

	context.Next()
}
