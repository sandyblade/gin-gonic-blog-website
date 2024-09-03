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

package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id        uint64        `json:"id" gorm:"primary_key"`
	ParentId  sql.NullInt64 `json:"parent_id" gorm:"index;default:null;"`
	ArticleId uint64        `json:"article_id" gorm:"index;not null"`
	UserId    uint64        `json:"user_id" gorm:"index;not null"`
	Comment   string        `json:"comment"  gorm:"type:text;not null"`
	CreatedAt time.Time     `gorm:"index;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time     `gorm:"index;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Comment) TableName() string {
	return "comments"
}
