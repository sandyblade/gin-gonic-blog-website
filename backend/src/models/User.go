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

type User struct {
	Id           uint64         `json:"id" gorm:"primary_key"`
	Email        string         `json:"email" gorm:"index;size:191;not null"`
	Phone        string         `json:"phone" gorm:"index;size:191;default:null"`
	Password     string         `json:"password" gorm:"index;size:255;not null"`
	Salt         string         `json:"salt" gorm:"index;size:255;"`
	Image        sql.NullString `json:"image" gorm:"index;size:191;default:null;"`
	FirstName    sql.NullString `json:"first_name" gorm:"index;size:191;default:null;"`
	LastName     sql.NullString `json:"last_name" gorm:"index;size:191;default:null;"`
	Gender       sql.NullString `json:"gender" gorm:"index;size:2;default:null;"`
	Country      sql.NullString `json:"country" gorm:"index;size:191;default:null;"`
	JobTitle     sql.NullString `json:"job_title" gorm:"index;size:191;default:null;"`
	Facebook     sql.NullString `json:"facebook" gorm:"index;size:191;default:null;"`
	Instagram    sql.NullString `json:"instagram" gorm:"index;size:191;default:null;"`
	Twitter      sql.NullString `json:"twitter" gorm:"index;size:191;default:null;"`
	LinkedIn     sql.NullString `json:"linked_in" gorm:"index;size:191;default:null;"`
	Address      sql.NullString `json:"address"  gorm:"type:text;default:null;"`
	AboutMe      sql.NullString `json:"about_me"  gorm:"type:text;default:null;"`
	ResetToken   sql.NullString `json:"reset_token" gorm:"index;size:36;default:null;"`
	ConfirmToken sql.NullString `json:"confirm_token" gorm:"index;size:36;default:null;"`
	Confirmed    uint8          `json:"confirmed" gorm:"index;default:0"`
	CreatedAt    time.Time      `gorm:"index;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"index;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
