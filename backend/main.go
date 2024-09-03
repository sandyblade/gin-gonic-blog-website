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

package main

import (
	_db "api/backend/src/db"
	_routes "api/backend/src/routes"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/kristijorgji/goseeder"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
	_db.Config()
	goseeder.WithSeeder(connectToDbOrDie, func() {})
	db := _db.SetupDB()
	db.LogMode(true)
	r := _routes.SetupRoutes(db)
	r.Run("0.0.0.0:" + os.Getenv("APP_PORT"))
}

func connectToDbOrDie() *sql.DB {
	dbDriver := os.Getenv("DB_CONNECTION")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		url.QueryEscape(dbPassword),
		dbHost,
		dbPort,
		dbName,
	)
	con, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	return con
}
