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

func ArticleRoutes() []RouteSource {
	routes := []RouteSource{
		{
			"/api/article/list",
			"GET",
			false,
			controllers.ArticleList,
		},
		{
			"/api/article/create",
			"POST",
			true,
			controllers.ArticleCreate,
		},
		{
			"/api/article/read/:slug",
			"GET",
			false,
			controllers.ArticleRead,
		},
		{
			"/api/article/update/:id",
			"PATCH",
			true,
			controllers.ArticleUpdate,
		},
		{
			"/api/article/remove/:id",
			"DELETE",
			false,
			controllers.ArticleRemove,
		},
		{
			"/api/article/user",
			"GET",
			true,
			controllers.ArticleListUser,
		},
		{
			"/api/article/words",
			"GET",
			false,
			controllers.ArticleWords,
		},
		{
			"/api/article/upload",
			"POST",
			true,
			controllers.ArticleUpload,
		},
	}
	return routes
}
