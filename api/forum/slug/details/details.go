package details

import (
	"github.com/VictoriaOtm/forum-api/database/stores/forumstore"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/valyala/fasthttp"
)

var responseErrorForumNotExists = e.MakeError("forum doesn't exist")

// Получение информации о форуме
func Details(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	frm := forumstore.Pool.Acquire()
	defer forumstore.Pool.Utilize(frm)

	err := frm.Get(ctx.UserValue("slug"))
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(responseErrorForumNotExists)
		return
	}

	ctx.Write(frm.MustMarshalJSON())
}
