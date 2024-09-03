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
	helpers "api/backend/src/helpers"
	models "api/backend/src/models"
	schema "api/backend/src/schema"
	services "api/backend/src/services"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func AccountActivity(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)
	page := 1
	limit := 10
	order_by := "id"
	order_dir := "desc"
	offset := ((page - 1) * limit)
	var data []models.Activity

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
		db = db.Where("event LIKE ? OR description LIKE ?", c.Query("search"), c.Query("search"))
	}

	db = db.Limit(limit).Offset(offset).Order(order_by + " " + order_dir).Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func AccountRefresh(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": services.JWTAuthService().GenerateToken(int(user.Id), user.Email, true)})
}

func AccountDetail(c *gin.Context) {

	auth := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.Where("id = ?", auth["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	var payload schema.UserProfileSchema
	payload.Email = user.Email
	payload.Phone = user.Phone
	payload.FirstName = user.FirstName.String
	payload.LastName = user.LastName.String
	payload.Gender = user.Gender.String
	payload.Country = user.Country.String
	payload.JobTitle = user.JobTitle.String
	payload.Facebook = user.Facebook.String
	payload.Twitter = user.Twitter.String
	payload.LinkedIn = user.LinkedIn.String
	payload.Instagram = user.Instagram.String
	payload.AboutMe = user.AboutMe.String
	payload.Address = user.Address.String

	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": true, "data": payload})
}

func AccountUpdate(c *gin.Context) {

	authUser := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input schema.UserProfileSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if err := db.Where("email = ? AND id != ?", input.Email, authUser["id"]).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email address already exists !!"})
		return
	}

	if len(strings.TrimSpace(input.Phone)) > 0 {
		if err := db.Where("phone = ? AND id != ?", input.Phone, authUser["id"]).First(&user).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists !!"})
			return
		}
	}

	if err := db.Where("id = ?", authUser["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var _user models.User
	_user.Email = input.Email
	_user.Phone = input.Phone
	_user.FirstName = sql.NullString{String: input.FirstName, Valid: true}
	_user.LastName = sql.NullString{String: input.LastName, Valid: true}
	_user.Gender = sql.NullString{String: input.Gender, Valid: true}
	_user.Country = sql.NullString{String: input.Country, Valid: true}
	_user.JobTitle = sql.NullString{String: input.JobTitle, Valid: true}
	_user.Facebook = sql.NullString{String: input.Facebook, Valid: true}
	_user.Instagram = sql.NullString{String: input.Instagram, Valid: true}
	_user.LinkedIn = sql.NullString{String: input.LinkedIn, Valid: true}
	_user.Twitter = sql.NullString{String: input.Twitter, Valid: true}
	_user.Address = sql.NullString{String: input.Address, Valid: true}
	_user.AboutMe = sql.NullString{String: input.AboutMe, Valid: true}
	db.Model(&user).Updates(_user)

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Update Profile",
		Description: "Edit user profile account",
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"message": "Your profile has been changed!", "status": true, "data": nil})

}

func AccountPassword(c *gin.Context) {

	authUser := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input schema.UserPasswordSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.OldPassword)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The old_password field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The password field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.ConfirmPassword)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The password_confirm field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you have to enter at least 8 digit!"})
		return
	}

	if strings.TrimSpace(input.Password) != strings.TrimSpace(input.ConfirmPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "These passwords don't match!"})
		return
	}

	if err := db.Where("id = ?", authUser["id"]).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	decrypt := helpers.Decrypt(user.Password, user.Salt)

	if input.OldPassword != decrypt {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect current password!"})
		return
	}

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	encrypted := helpers.Encrypt(input.Password, key)

	var _user models.User
	_user.Password = encrypted
	_user.Salt = key
	db.Model(&user).Updates(_user)

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Change Password",
		Description: "Change new password account",
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"message": "Your password has been changed!"})

}

func AccountUpload(c *gin.Context) {

	authUser := c.MustGet("claims").(jwt.MapClaims)
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", authUser["id"]).First(&user).Error; err != nil {
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

		var _user models.User
		_user.Image = sql.NullString{result, true}
		db.Model(&user).Updates(_user)

		Activity := models.Activity{
			UserId:      uint64(user.Id),
			Event:       "Upload Profile Image",
			Description: "Upload new user profile image",
		}
		db.Create(&Activity)

	}

	// File saved successfully. Return proper result
	c.JSON(http.StatusOK, gin.H{"data": result})

}
