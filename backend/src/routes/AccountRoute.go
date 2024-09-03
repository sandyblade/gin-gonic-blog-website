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

func AccountRoutes() []RouteSource {
	routes := []RouteSource{
		{
			"/api/account/detail",
			"GET",
			true,
			controllers.AccountDetail,
		},
		{
			"/api/account/update",
			"POST",
			true,
			controllers.AccountUpdate,
		},
		{
			"/api/account/password",
			"POST",
			true,
			controllers.AccountPassword,
		},
		{
			"/api/account/upload",
			"POST",
			true,
			controllers.AccountUpload,
		},
		{
			"/api/account/token",
			"POST",
			true,
			controllers.AccountRefresh,
		},
		{
			"/api/account/activity",
			"GET",
			true,
			controllers.AccountActivity,
		},
	}
	return routes
}
