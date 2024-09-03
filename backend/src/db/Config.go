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

package db

import (
	models "api/backend/src/models"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func SetupDB() *gorm.DB {
	godotenv.Load(".env")
	USER := os.Getenv("DB_USERNAME")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	PORT := os.Getenv("DB_PORT")
	DBNAME := os.Getenv("DB_DATABASE")
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	db, err := gorm.Open(os.Getenv("DB_CONNECTION"), URL)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Config() {
	db := SetupDB()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Activity{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Article{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Comment{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Notification{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Viewer{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Article{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Comment{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Comment{}).AddForeignKey("article_id", "articles(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Comment{}).AddForeignKey("parent_id", "comments(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Viewer{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Viewer{}).AddForeignKey("article_id", "articles(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Notification{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models.Activity{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
}
