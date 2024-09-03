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
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func AuthLogin(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input schema.UserLoginSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The password field is required.!"})
		return
	}

	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user " + input.Email + " not found!"})
		return
	}

	decrypt := helpers.Decrypt(user.Password, user.Salt)

	if user.Confirmed == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to confirm your account. We have sent you an activation code, please check your email.!"})
		return
	}

	if input.Password != decrypt {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password!"})
		return
	}

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Sign In",
		Description: "Sign in to application",
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"token": services.JWTAuthService().GenerateToken(int(user.Id), user.Email, true)})
}

func AuthRegister(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input schema.UserRegisterSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The passwword field is required.!"})
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

	if err := db.Where("email = ?", input.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email already exists"})
		return
	}

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	encrypted := helpers.Encrypt(input.Password, key)
	token := (uuid.New()).String()

	User := models.User{
		Email:        input.Email,
		Password:     encrypted,
		Confirmed:    0,
		ConfirmToken: sql.NullString{String: token, Valid: true},
		Salt:         key,
	}
	db.Create(&User)

	Activity := models.Activity{
		UserId:      uint64(User.Id),
		Event:       "Sign Up",
		Description: "Register new user account",
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"message": "Your account has been created. Please check your email for the confirmation message we just sent you."})
}

func AuthConfirm(c *gin.Context) {

	var user models.User

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("confirm_token = ? AND confirmed = ? ", c.Param("token"), 0).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Exec("UPDATE users SET confirmed = 1, confirm_token = NULL, updated_at = NOW() WHERE confirm_token = ? ", c.Param("token"))

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Email Verification",
		Description: "Confirm new member registration account",
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"message": "Your registration is complete. Now you can login."})
}

func AuthEmailForgot(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input schema.UserForgotSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "We can't find a user with that e-mail address."})
		return
	}

	token := (uuid.New()).String()

	var _user models.User
	_user.ResetToken = sql.NullString{String: token, Valid: true}
	db.Model(&user).Updates(_user)

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Forgot Password",
		Description: "Request reset password link",
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"message": "We have e-mailed your password reset link!", "token": token})
}

func AuthEmailReset(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	var input schema.UserResetSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(input.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The email field is required.!"})
		return
	}

	if len(strings.TrimSpace(input.Password)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The passwword field is required.!"})
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

	if err := db.Where("email = ? AND reset_token = ?", input.Email, c.Param("token")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This password and email reset token is invalid."})
		return
	}

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	encrypted := helpers.Encrypt(input.Password, key)

	db.Exec("UPDATE users SET confirmed = 1, password = ?, salt = ?,  reset_token = NULL, updated_at = NOW() WHERE reset_token = ? ", encrypted, key, c.Param("token"))

	Activity := models.Activity{
		UserId:      uint64(user.Id),
		Event:       "Reset Password",
		Description: "Reset account password",
	}
	db.Create(&Activity)

	c.JSON(http.StatusOK, gin.H{"message": "Your password has been reset!"})
}
