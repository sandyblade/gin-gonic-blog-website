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

package schema

import (
	"database/sql"
	"time"
)

type ArticleSchema struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
	Status      uint8    `json:"status"`
}

type ArticleListSchema struct {
	Id          uint64         `json:"id"`
	Title       string         `json:"title"`
	Image       sql.NullString `json:"image"`
	Description string         `json:"description"`
	Categories  string         `json:"categories"`
	Tags        string         `json:"tags"`
	Status      uint8          `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	FirstName   sql.NullString `json:"first_name"`
	LastName    sql.NullString `json:"last_name"`
	Gender      sql.NullString `json:"gender"`
	Email       string         `json:"email"`
}

type CommentSchema struct {
	ParentId sql.NullInt64 `json:"parent_id"`
	Comment  string        `json:"comment"`
}

type CommentListSchema struct {
	Id        uint64              `json:"id"`
	ParentId  uint64              `json:"parent_id"`
	Comment   string              `json:"comment"`
	FirstName sql.NullString      `json:"first_name"`
	LastName  sql.NullString      `json:"last_name"`
	Gender    sql.NullString      `json:"gender"`
	Email     string              `json:"email"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Childern  []CommentListSchema `json:"childern"`
}
