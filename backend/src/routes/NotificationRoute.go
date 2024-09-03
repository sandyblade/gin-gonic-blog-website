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

func NotificationRoutes() []RouteSource {
	routes := []RouteSource{
		{
			"/api/notification/list",
			"GET",
			true,
			controllers.NotificationList,
		},
		{
			"/api/notification/read/:id",
			"GET",
			true,
			controllers.NotificationRead,
		},
		{
			"/api/notification/remove/:id",
			"DELETE",
			true,
			controllers.NotificationRemove,
		},
	}
	return routes
}
