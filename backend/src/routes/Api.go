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

package routes

import (
	middleware "api/backend/src/middleware"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RouteSource struct {
	Name   string
	Method string
	Auth   bool
	Result func(c *gin.Context)
}

func SetupRoutes(db *gorm.DB) *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	var AppRoutes = []RouteSource{}
	AppRoutes = append(AppRoutes, AuthRoutes()...)
	AppRoutes = append(AppRoutes, AccountRoutes()...)
	AppRoutes = append(AppRoutes, NotificationRoutes()...)
	AppRoutes = append(AppRoutes, ArticleRoutes()...)
	AppRoutes = append(AppRoutes, CommentRoutes()...)

	for _, element := range AppRoutes {
		if element.Method == "POST" {
			if element.Auth == true {
				r.POST(element.Name, middleware.AuthorizeJWT(), element.Result)
			} else {
				r.POST(element.Name, element.Result)
			}
		} else if element.Method == "DELETE" {
			if element.Auth == true {
				r.DELETE(element.Name, middleware.AuthorizeJWT(), element.Result)
			} else {
				r.DELETE(element.Name, element.Result)
			}
		} else if element.Method == "PATCH" {
			if element.Auth == true {
				r.PATCH(element.Name, middleware.AuthorizeJWT(), element.Result)
			} else {
				r.PATCH(element.Name, element.Result)
			}
		} else {
			if element.Auth == true {
				r.GET(element.Name, middleware.AuthorizeJWT(), element.Result)
			} else {
				r.GET(element.Name, element.Result)
			}
		}
	}

	r.MaxMultipartMemory = 8 << 20
	r.Static("uploads", os.Getenv("UPLOAD_PATH"))
	return r
}
