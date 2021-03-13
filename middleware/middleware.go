package middleware

import (
	"net/http"
	"todo/auth"

	"github.com/gin-gonic/gin"
)

//Middleware function to authorize the user
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//check for token validity
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "You need to be authorized to access this route")
			c.Abort()
			return
		}
		//extract user id from token
		id, err := auth.ExtractTokenID(c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Internal server error")
			c.Abort()
			return
		}
		// set the user-id in gin context
		c.Set("user-id", id)
		c.Next()
	}

}
