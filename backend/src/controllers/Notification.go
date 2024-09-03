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

package controllers

import (
	models "api/backend/src/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

func NotificationList(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)
	page := 1
	limit := 10
	order_by := "id"
	order_dir := "desc"
	offset := ((page - 1) * limit)
	var data []models.Notification

	if len(strings.TrimSpace(c.Query("limit"))) > 0 {
		limit_, err := strconv.ParseInt(c.Query("limit"), 0, 32)
		if err == nil {
			limit = int(limit_)
		}
	}

	if len(strings.TrimSpace(c.Query("page"))) > 0 {
		page_, err := strconv.ParseInt(c.Query("page"), 0, 32)
		if err == nil {
			page = int(page_)
		}
	}

	if len(strings.TrimSpace(c.Query("order_by"))) > 0 {
		order_by = c.Query("order_by")
	}

	if len(strings.TrimSpace(c.Query("order_dir"))) > 0 {
		order_dir = c.Query("order_dir")
	}

	db = db.Where("user_id = ?", auth["id"])

	if len(strings.TrimSpace(c.Query("search"))) > 0 {
		db = db.Where("subject LIKE ? OR message LIKE ?", c.Query("search"), c.Query("search"))
	}

	db = db.Limit(limit).Offset(offset).Order(order_by + " " + order_dir).Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func NotificationRead(c *gin.Context) {

	var notification models.Notification

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), auth["id"]).First(&notification).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notification})

}

func NotificationRemove(c *gin.Context) {

	var notification models.Notification
	var user models.User

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), auth["id"]).First(&notification).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	db.Delete(&notification)

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Delete notification",
		Description: "The user delete notification with subject " + notification.Subject,
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"data": notification})

}
