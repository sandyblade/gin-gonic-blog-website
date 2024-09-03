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
	schema "api/backend/src/schema"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

func ArticleList(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	page := 1
	limit := 10
	order_by := "id"
	order_dir := "desc"
	offset := ((page - 1) * limit)
	var data []models.Article

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

	db = db.Where("status = 1")

	if len(strings.TrimSpace(c.Query("search"))) > 0 {
		db = db.Where("title LIKE ? OR description LIKE ? OR content LIKE ? OR categories LIKE ? OR tags LIKE ?", c.Query("search"), c.Query("search"))
	}

	db = db.Limit(limit).Offset(offset).Order(order_by + " " + order_dir).Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func ArticleListUser(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)
	page := 1
	limit := 10
	order_by := "id"
	order_dir := "desc"
	offset := ((page - 1) * limit)
	var data []models.Article

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
		db = db.Where("title LIKE ? OR description LIKE ? OR content LIKE ? OR categories LIKE ? OR tags LIKE ?", c.Query("search"), c.Query("search"))
	}

	db = db.Limit(limit).Offset(offset).Order(order_by + " " + order_dir).Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func ArticleCreate(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	var article models.Article

	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	var input schema.ArticleSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Title)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The title field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Description)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The description field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Content)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The content field is required.!"})
		return
	}

	if err := db.Where("title = ?", input.Title).First(&article).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The title of article already exists !!"})
		return
	}

	NewArticle := models.Article{
		UserId:      uint64(user.Id),
		Title:       input.Title,
		Description: input.Description,
		Content:     input.Content,
	}
	db.Create(&NewArticle)

}

func ArticleRead(c *gin.Context) {

}

func ArticleUpdate(c *gin.Context) {

}

func ArticleRemove(c *gin.Context) {

}

func ArticleUpload(c *gin.Context) {

}

func ArticleWords(c *gin.Context) {

}
