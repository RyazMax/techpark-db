package routers

import (
	"techpark-db/controllers"
	"techpark-db/database"

	"github.com/astaxie/beego"
)

func InitRouter(db *database.DB) {
	ns :=
		beego.NewNamespace("/api",
			beego.NSNamespace("/user",
				beego.NSRouter("/:nickname/profile", &controllers.UserProfileController{DB: db}),
				beego.NSRouter("/:nickname/create", &controllers.UserController{DB: db})),
			beego.NSNamespace("/thread/:id",
				beego.NSRouter("/create", &controllers.ThreadCreateController{DB: db}),
				beego.NSRouter("/details", &controllers.ThreadDetailsController{DB: db}),
				beego.NSRouter("/posts", &controllers.ThreadPostsController{DB: db}),
				beego.NSRouter("/vote", &controllers.ThreadVoteController{DB: db})),
			beego.NSRouter("/service/*", &controllers.ServiceController{DB: db}),
			beego.NSRouter("/post/:id/details", &controllers.PostController{DB: db}),
			beego.NSNamespace("/forum",
				beego.NSRouter("/create", &controllers.ForumController{DB: db}),
				beego.NSRouter("/:slug/*", &controllers.ForumSlugController{DB: db})))
	beego.AddNamespace(ns)
}
