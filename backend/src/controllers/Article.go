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
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func ArticleList(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	page := 1
	limit := 10
	order_by := "id"
	order_dir := "desc"

	var data []schema.ArticleListSchema

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

	offset := ((page - 1) * limit)
	db = db.Table("articles").Select(`
		articles.id,
		articles.title,
		articles.image,
		articles.description,
		articles.categories,
		articles.tags,
		articles.status,
		articles.created_at,
		articles.updated_at,
		users.first_name,
		users.last_name,
		users.gender,
		users.email
	`).Where("articles.status = 1")

	if len(strings.TrimSpace(c.Query("search"))) > 0 {
		db = db.Where("articles.title LIKE ? OR articles.description LIKE ? OR articles.content LIKE ? OR articles.categories LIKE ? OR articles.tags LIKE ?", c.Query("search"), c.Query("search"))
	}

	db.Limit(limit).Offset(offset).Order(order_by + " " + order_dir).Joins("inner join users on articles.user_id = users.id").Scan(&data)

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
		Slug:        slug.Make(input.Title),
		Description: input.Description,
		Categories:  strings.Join(input.Categories, ","),
		Tags:        strings.Join(input.Tags, ","),
		Status:      input.Status,
		Content:     input.Content,
	}
	db.Create(&NewArticle)

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Create New Article",
		Description: "A new article with title `" + input.Title + "` has been created. ",
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": NewArticle})
}

func ArticleRead(c *gin.Context) {

	var article models.Article
	var user models.User
	var totalViewer int64
	var _article models.Article

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("slug = ?", c.Param("slug")).First(&article).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	auth := c.MustGet("claims").(jwt.MapClaims)
	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	if user.Id != article.UserId {
		db.Model(&models.Viewer{}).Where("article_id = ? AND user_id = ?", article.Id, auth["id"]).Count(&totalViewer)
		if totalViewer == 0 {
			Viewer := models.Viewer{
				UserId:    uint64(user.Id),
				ArticleId: uint64(article.Id),
				Status:    0,
			}
			db.Create(&Viewer)
			Activity := models.Activity{
				UserId:      uint64(user.Id),
				Event:       "Read Article",
				Description: "An article with title `" + article.Title + "` has been viewed. ",
			}
			db.Create(&Activity)
			db.Model(&models.Viewer{}).Where("article_id = ? AND user_id = ?", article.Id, auth["id"]).Count(&totalViewer)
			_article.TotalViewer = uint16(totalViewer)
			db.Model(&article).Updates(_article)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": article})

}

func ArticleUpdate(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	var article models.Article
	var articleModel models.Article

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

	if err := db.Where("title = ? AND id <> ?", input.Title, c.Param("id")).First(&article).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The title of article already exists !!"})
		return
	}

	if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), auth["id"]).First(&articleModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var _article models.Article
	_article.Title = input.Title
	_article.Slug = slug.Make(input.Title)
	_article.Description = input.Description
	_article.Categories = strings.Join(input.Categories, ",")
	_article.Tags = strings.Join(input.Categories, ",")
	_article.Status = input.Status
	_article.Content = input.Title
	db.Model(&articleModel).Updates(_article)

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Update Article",
		Description: "The user editing article with title " + input.Title,
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": _article})
}

func ArticleRemove(c *gin.Context) {

	var article models.Article
	var user models.User

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), auth["id"]).First(&article).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	db.Delete(&article)

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Delete Article",
		Description: "The user delete article with title " + article.Title,
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"data": article})
}

func ArticleUpload(c *gin.Context) {

	var user models.User
	var article models.Article

	authUser := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id = ?", authUser["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), authUser["id"]).First(&article).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	godotenv.Load(".env")
	file, err := c.FormFile("file")
	currentTime := time.Now()
	datePath := currentTime.Format("2006-01-02")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	path := os.Getenv("UPLOAD_PATH")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	realPath := path + "/" + datePath
	if _, err := os.Stat(realPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(realPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, realPath+"/"+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	result := datePath + "/" + newFileName

	if result != "" {

		if user.Image.Valid {
			image := sql.NullString{user.Image.String, true}
			e := os.Remove(path + "/" + image.String)
			if e != nil {
				log.Fatal(e)
			}
		}

		var _article models.Article
		_article.Image = sql.NullString{result, true}
		db.Model(&article).Updates(_article)

		Activity := models.Activity{
			UserId:      uint64(user.Id),
			Event:       "Upload Article Image",
			Description: "Your upload file has been successfully !!",
		}
		db.Create(&Activity)

	}

	// File saved successfully. Return proper result
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func ArticleWords(c *gin.Context) {
	max := 10
	if len(strings.TrimSpace(c.Query("max"))) > 0 {
		max_, err := strconv.ParseInt(c.Query("max"), 0, 32)
		if err == nil {
			max = int(max_)
		}
	}
	var data []string
	for i := 0; i < max; i++ {
		word := faker.Word() + " " + faker.Word()
		data = append(data, strings.Title(strings.ToLower(word)))
	}
	sort.Sort(sort.StringSlice(data))
	c.JSON(http.StatusOK, gin.H{"data": data})
}
