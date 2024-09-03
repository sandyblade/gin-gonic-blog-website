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

type UserProfileSchema struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Country   string `json:"country"`
	JobTitle  string `json:"job_title"`
	Facebook  string `json:"facebook"`
	Instagram string `json:"instagram"`
	LinkedIn  string `json:"linked_in"`
	Twitter   string `json:"twitter"`
	Address   string `json:"address"`
	AboutMe   string `json:"about_me"`
}

type UserPasswordSchema struct {
	OldPassword     string `json:"old_password"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"password_confirm"`
}
