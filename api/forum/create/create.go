package create

import (
	"github.com/VictoriaOtm/forum-api/database/stores/forumstore"
	"github.com/VictoriaOtm/forum-api/database/stores/userstore"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/valyala/fasthttp"
)

var responseErrorUserNotExists = e.MakeError("user not exists")

// Insert - Создание форума
func Create(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	frm := forumstore.Pool.Acquire()
	defer forumstore.Pool.Utilize(frm)

	frm.MustUnmarshalJSON(ctx.Request.Body())

	err := frm.Get(frm.Slug)
	if err == nil {
		ctx.SetStatusCode(409)
		ctx.Write(frm.MustMarshalJSON())
		return
	}

	usr := userstore.Pool.Acquire()
	defer userstore.Pool.Utilize(usr)

	err = usr.Get(frm.User)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(responseErrorUserNotExists)
		return
	}

	frm.User = usr.Nickname
	err = frm.Insert()
	if err != nil {
		ctx.SetStatusCode(500)
		return
	}

	ctx.SetStatusCode(201)
	ctx.Write(frm.MustMarshalJSON())
}
