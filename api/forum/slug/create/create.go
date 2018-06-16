package create

import (
	"strings"

	"github.com/VictoriaOtm/forum-api/database/stores/forumstore"
	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	"github.com/VictoriaOtm/forum-api/database/stores/userstore"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/VictoriaOtm/forum-api/helpers/unsafe_map"
	"github.com/valyala/fasthttp"
)

var (
	responseErrorForumNotExists = e.MakeError("forum not exists")
	responseErrorUserNotExists  = e.MakeError("user not exists")
)

// Insert - создание ветки
func Create(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	thr := threadstore.Pool.Acquire()
	defer threadstore.Pool.Utilize(thr)

	thr.MustUnmarshalJSON(ctx.PostBody())

	if thr.Slug != nil {
		err := thr.GetBySlug(*thr.Slug)
		if err == nil {
			ctx.Write(thr.MustMarshalJSON())
			ctx.SetStatusCode(409)
			return
		}
	}

	frm := forumstore.Pool.Acquire()
	defer forumstore.Pool.Utilize(frm)

	forumSlug := ctx.UserValue("slug")
	err := frm.Get(forumSlug)
	if err != nil {
		ctx.Write(responseErrorForumNotExists)
		ctx.SetStatusCode(404)
		return
	}

	usr := userstore.Pool.Acquire()
	defer userstore.Pool.Utilize(usr)

	err = usr.Get(thr.Author)
	if err != nil {
		ctx.Write(responseErrorUserNotExists)
		ctx.SetStatusCode(404)
		return
	}

	thr.Forum = frm.Slug
	thr.Author = usr.Nickname

	err = thr.Insert()
	if err != nil {
		ctx.SetStatusCode(500)
		return
	}
	userstore.StoreInForumUserTable(thr.Forum, []interface{}{*usr})

	if thr.Slug != nil {
		unsafe_map.SlugIDMap.StoreSlug(strings.ToLower(*thr.Slug), thr.Id)
	}
	ctx.SetStatusCode(201)
	ctx.Write(thr.MustMarshalJSON())
}
