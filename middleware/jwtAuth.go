package middleware

import (
	"dairy_service/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// The JWTAuthMiddleware returns a Gin HandlerFunc function. This function expects a context for which it tries to validate the
// JWT in the header. If it is invalid, an error response is returned. If not, the Next() function on the context is called.
// In our case, the controller function for the protected route is called.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := helper.ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		context.Next()
	}
}
