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
	"time"
)

type Notification struct {
	Id        uint64    `json:"id" gorm:"primary_key"`
	UserId    uint64    `json:"user_id" gorm:"index;not null"`
	Subject   string    `json:"subject" gorm:"index;size:255;not null"`
	Message   string    `json:"message"  gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
