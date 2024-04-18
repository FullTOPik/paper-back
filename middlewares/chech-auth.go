package middlewares

import (
	"net/http"
	"paper_back/exceptions"
	token_service "paper_back/services/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func parseToken(token string) (*token_service.CustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &token_service.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return token_service.SecretAccess, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*token_service.CustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code int

		code = http.StatusOK
		token, err := ctx.Cookie("token")

		if token == "" || err != nil {
			code = http.StatusUnauthorized
		} else {
			claims, err := parseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = http.StatusForbidden
				default:
					code = http.StatusUnauthorized
				}
			}

			if claims != nil {
				ctx.Set("user_id", claims.Id)
				ctx.Set("username", claims.Username)
				ctx.Set("role", claims.Role)
			}
		}

		if code != http.StatusOK {
			ctx.JSON(code, exceptions.ErrorWithStatus{
				Code:    int32(code),
				Message: "Unauthorized",
			})

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
