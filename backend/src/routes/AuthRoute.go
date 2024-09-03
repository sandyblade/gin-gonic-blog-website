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

package routes

import (
	controllers "api/backend/src/controllers"
)

func AuthRoutes() []RouteSource {
	routes := []RouteSource{
		{
			"/api/auth/login",
			"POST",
			false,
			controllers.Login,
		},
		{
			"/api/auth/register",
			"POST",
			false,
			controllers.Register,
		},
		{
			"/api/auth/email/forgot",
			"POST",
			false,
			controllers.EmailForgot,
		},
		{
			"/api/auth/email/reset/:token",
			"POST",
			false,
			controllers.EmailReset,
		},
		{
			"/api/auth/confirm/:token",
			"GET",
			false,
			controllers.Confirm,
		},
	}
	return routes
}
