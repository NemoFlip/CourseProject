package middleware

import (
	"CourseProject/auth_service/pkg/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func CheckAuthorization(tm *auth.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if len(authHeader) == 0 {
			log.Printf("authroization token is absent")
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if bearerToken[0] != "Bearer" || len(bearerToken) != 2 {
			log.Printf("invalid format of auth token")
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]
		err := tm.ValidateToken(tokenString)
		if err != nil {
			log.Println(err)
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

}
