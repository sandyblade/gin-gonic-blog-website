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

func CommentRoutes() []RouteSource {
	routes := []RouteSource{
		{
			"/api/comment/list/:id",
			"GET",
			false,
			controllers.CommentList,
		},
		{
			"/api/comment/create/:id",
			"POST",
			true,
			controllers.CommentCreate,
		},
		{
			"/api/comment/remove/:id",
			"DELETE",
			true,
			controllers.CommentRemove,
		},
	}
	return routes
}
