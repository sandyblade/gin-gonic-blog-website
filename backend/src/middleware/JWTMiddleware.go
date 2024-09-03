/**
 * This file is part of the Sandy Andryanto Blog Applicatione.
 *
 * @author     Sandy Andryanto <sandy.andryanto.blade@gmail.com>
 * @copyright  2024
 *
 * For the full copyright and license information,
 * please view the LICENSE.md file that was distributed
 * with this source code.
 */

package middleware

import (
	service "api/backend/src/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(strings.TrimSpace(authHeader)) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			tokenString := authHeader[len(BEARER_SCHEMA):]
			token, err := service.JWTAuthService().ValidateToken(tokenString)
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				c.Set("claims", claims)
				fmt.Println(claims)
			} else {
				fmt.Println(err)
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		}

	}
}
