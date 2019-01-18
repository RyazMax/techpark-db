package routers

import (
	"techpark-db/controllers"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

/*func InitRouter(db *database.DB) {
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
}*/

// Handler fasthttp Handler for api
var Handler fasthttp.RequestHandler

func init() {
	router := fasthttprouter.New()

	//TODO: forum/create
	router.POST("/api/forum/:slug/create", controllers.ForumSlugCreate)
	router.GET("/api/forum/:slug/details", controllers.ForumSlugDetails)
	router.GET("/api/forum/:slug/threads", controllers.ForumSlugThreads)
	router.GET("/api/forum/:slug/users", controllers.ForumSlugUsers)

	router.GET("/api/post/:id/details", controllers.PostDetailsGet)
	router.POST("/api/post/:id/details", controllers.PostDetailsPost)

	router.POST("/api/service/clear", controllers.ServiceClear)
	router.GET("/api/service/status", controllers.ServiceStatus)

	router.POST("/api/thread/:slug_or_id/create", controllers.ThreadCreatePosts)
	router.GET("/api/thread/:slug_or_id/details", controllers.ThreadDetailsGet)
	router.POST("/api/thread/:slug_or_id/details", controllers.ThreadDetailsPost)
	router.GET("/api/thread/:slug_or_id/posts", controllers.ThreadGetPosts)
	router.POST("/api/thread/:slug_or_id/vote", controllers.ThreadVote)

	router.POST("/api/user/:nickname/create", controllers.UserCreate)
	router.GET("/api/user/:nickname/profile", controllers.UserGetProfile)
	router.POST("/api/user/:nickname/profile", controllers.UserUpdate)

	Handler = func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Path()) == "/api/forum/create" {
			controllers.ForumCreate(ctx)
		} else {
			router.Handler(ctx)
		}
	}
}
