package main

import (
	FCreate "github.com/VictoriaOtm/forum-api/api/forum/create"
	FSCreate "github.com/VictoriaOtm/forum-api/api/forum/slug/create"
	FSDetails "github.com/VictoriaOtm/forum-api/api/forum/slug/details"
	FSThreads "github.com/VictoriaOtm/forum-api/api/forum/slug/threads"
	FSUsers "github.com/VictoriaOtm/forum-api/api/forum/slug/users"

	PIGetDetails "github.com/VictoriaOtm/forum-api/api/post/get_details"
	PIPostDetails "github.com/VictoriaOtm/forum-api/api/post/post_details"

	SClear "github.com/VictoriaOtm/forum-api/api/service/clear"
	SStatus "github.com/VictoriaOtm/forum-api/api/service/status"

	TSCreate "github.com/VictoriaOtm/forum-api/api/thread/create"
	TSGetDetails "github.com/VictoriaOtm/forum-api/api/thread/get_details"
	TSPostDetails "github.com/VictoriaOtm/forum-api/api/thread/post_details"
	TSPosts "github.com/VictoriaOtm/forum-api/api/thread/posts"
	TSVote "github.com/VictoriaOtm/forum-api/api/thread/vote"

	UCreate "github.com/VictoriaOtm/forum-api/api/user/create"
	UGetProfile "github.com/VictoriaOtm/forum-api/api/user/get_profile"
	UPostProfile "github.com/VictoriaOtm/forum-api/api/user/post_profile"

	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	_ "net/http/pprof"

	"github.com/VictoriaOtm/forum-api/database"
	"github.com/VictoriaOtm/forum-api/database/stores/forumstore"
	"github.com/VictoriaOtm/forum-api/database/stores/poststore"
	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	"github.com/VictoriaOtm/forum-api/database/stores/userstore"
)

func initRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/api/forum/:slug", FCreate.Create)

	router.POST("/api/forum/:slug/create", FSCreate.Create)
	router.GET("/api/forum/:slug/details", FSDetails.Details)
	router.GET("/api/forum/:slug/threads", FSThreads.Threads)
	router.GET("/api/forum/:slug/users", FSUsers.Users)

	router.GET("/api/post/:id/details", PIGetDetails.Details)
	router.POST("/api/post/:id/details", PIPostDetails.Details)

	router.POST("/api/service/clear", SClear.Clear)
	router.GET("/api/service/status", SStatus.Status)

	router.POST("/api/thread/:slug_or_id/create", TSCreate.Create)
	router.GET("/api/thread/:slug_or_id/details", TSGetDetails.Details)
	router.POST("/api/thread/:slug_or_id/details", TSPostDetails.Details)
	router.GET("/api/thread/:slug_or_id/posts", TSPosts.Posts)
	router.POST("/api/thread/:slug_or_id/vote", TSVote.Vote)

	router.POST("/api/user/:nickname/create", UCreate.Create)
	router.GET("/api/user/:nickname/profile", UGetProfile.Profile)
	router.POST("/api/user/:nickname/profile", UPostProfile.Profile)

	return router
}

func main() {
	log.SetFlags(log.Llongfile)
	database.InitSchema("res/init.sql")
	//go func() {
	//	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	//}()

	forumstore.PrepareStatements()
	poststore.PrepareStatements()
	threadstore.PrepareStatements()
	userstore.PrepareStatements()

	router := initRouter()

	log.Println(fasthttp.ListenAndServe("0.0.0.0:5000", router.Handler))
}
