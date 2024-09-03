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

type Article struct {
	Id           uint64         `json:"id" gorm:"primary_key"`
	UserId       uint64         `json:"user_id" gorm:"index;not null"`
	Image        sql.NullString `json:"image" gorm:"index;size:191;default:null;"`
	Title        string         `json:"title" gorm:"index;size:255;not null"`
	Slug         string         `json:"slug" gorm:"index;size:255;not null"`
	Description  string         `json:"description"  gorm:"type:text;not null"`
	Content      string         `json:"content"  gorm:"type:longtext;default:null;"`
	Categories   string         `json:"categories"  gorm:"type:longtext;default:null;"`
	Tags         string         `json:"tags"  gorm:"type:longtext;default:null;"`
	TotalComment uint16         `json:"total_comment" gorm:"index;default:0"`
	TotalViewer  uint16         `json:"total_viewer" gorm:"index;default:0"`
	Status       uint8          `json:"status" gorm:"index;default:0"`
	CreatedAt    time.Time      `gorm:"index;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"index;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Article) TableName() string {
	return "articles"
}
