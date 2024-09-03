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
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CommentList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var result []schema.CommentListSchema
	db.Raw(`
		SELECT 
			c.id,
			IFNULL(c.parent_id, 0) parent_id,
			c.comment,
			c.created_at,
			u.first_name,
			u.last_name,
			u.gender,
			u.email
		FROM comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.article_id = ?
		ORDER BY c.id DESC
	`, c.Param("id")).Scan(&result)
	commentTree := ArticleTree(result, 0)
	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": commentTree})
}

func CommentCreate(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	var article models.Article
	var input schema.CommentSchema
	var _article models.Article
	var totalComment int64

	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	if err := db.Where("id = ?", c.Param("id")).First(&article).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article with id " + c.Param("id") + " not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Comment)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The comment field is required.!"})
		return
	}

	NewComment := models.Comment{
		UserId:    uint64(user.Id),
		ArticleId: article.Id,
		ParentId:  input.ParentId,
		Comment:   input.Comment,
	}
	db.Create(&NewComment)

	Event := "Reply Comment"
	Description := "The user " + user.Email + " replied to your article with title `" + article.Title + "`."

	if input.ParentId.Valid {
		Activity := models.Activity{
			UserId:      uint64(user.Id),
			Event:       Event,
			Description: Description,
		}
		db.Create(&Activity)
	} else {
		Activity := models.Activity{
			UserId:      uint64(user.Id),
			Event:       Event,
			Description: Description,
		}
		db.Create(&Activity)
	}

	if user.Id != article.UserId {
		Notification := models.Notification{
			UserId:  uint64(user.Id),
			Subject: Event,
			Message: Description,
		}
		db.Create(&Notification)
	}

	db.Model(&models.Comment{}).Where("article_id = ?", article.Id).Count(&totalComment)
	_article.TotalComment = uint16(totalComment)
	db.Model(&article).Updates(_article)

	c.JSON(http.StatusOK, gin.H{"data": NewComment})
}

func CommentRemove(c *gin.Context) {

	var comment models.Comment
	var user models.User
	var article models.Article
	var totalComment int64
	var _article models.Article

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), auth["id"]).First(&comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	if err := db.Where("id = ? AND user_id = ? ", comment.ArticleId, auth["id"]).First(&article).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&comment)

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Delete Comment",
		Description: "The user " + user.Email + " deleted comment of article with title `" + article.Title + "`.",
	}
	db.Create(&Activity)

	db.Model(&models.Comment{}).Where("article_id = ?", article.Id).Count(&totalComment)
	_article.TotalComment = uint16(totalComment)
	db.Model(&article).Updates(_article)

	c.JSON(http.StatusOK, gin.H{"data": comment})

}

func ArticleTree(elements []schema.CommentListSchema, ParentId uint64) []schema.CommentListSchema {
	var branch []schema.CommentListSchema
	for _, element := range elements {

		if element.ParentId == ParentId {

			childern := []schema.CommentListSchema{}

			getChildren := ArticleTree(elements, element.Id)
			if len(getChildren) > 0 {
				childern = getChildren
			}

			var obj = schema.CommentListSchema{
				Id:        element.Id,
				ParentId:  element.ParentId,
				Comment:   element.Comment,
				CreatedAt: element.CreatedAt,
				FirstName: element.FirstName,
				LastName:  element.LastName,
				Gender:    element.Gender,
				Childern:  childern,
			}

			branch = append(branch, obj)
		}
	}
	return branch
}
