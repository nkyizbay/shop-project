package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var allowList = map[string]bool{
	"/register": true,
	"/login":    true,
}

func TokenMiddleware() gin.HandlerFunc {
	jwtSecretKey := viper.GetString("SHOP_GO_JWTKEY")

	return func(c *gin.Context) {
		if _, ok := allowList[c.Request.RequestURI]; ok {
			c.Next()
			return
		}

		cookie, err := c.Cookie("token")
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())

			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		// fmt.Println(jwtSecretKey)
		claim := Claims{}
		parsedTokenInfo, err := jwt.ParseWithClaims(cookie, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.String(http.StatusUnauthorized, "Please login again1")
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.String(http.StatusUnauthorized, "Please login again2")
			c.Abort()
			return
		}

		if !parsedTokenInfo.Valid {
			c.String(http.StatusForbidden, "Invalid token")

			c.Abort()
			return
		}

		c.Set("claim", claim)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// claim := c.Get("claim").(Claims)
		claim := c.Value("claim").(Claims)

		if claim.IsNotAdmin() {
			c.String(http.StatusForbidden, "You have no authority")

			c.Abort()
			return
		}

		c.Next()
	}
}
