/**
 * This file is part of the Sandy Andryanto Company Profile Website.
 *
 * @author     Sandy Andryanto <sandy.andryanto404@gmail.com>
 * @copyright  2024
 *
 * For the full copyright and license information,
 * please view the LICENSE.md file that was distributed
 * with this source code.
 */

package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kristijorgji/goseeder"
	"log"
	"net/url"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
	utilities.Config()
	goseeder.WithSeeder(connectToDbOrDie, func() {})
	db := utilities.SetupDB()
	db.LogMode(true)
	r := routes.SetupRoutes(db)
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
