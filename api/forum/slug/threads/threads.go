package threads

import (
	"bytes"

	"github.com/VictoriaOtm/forum-api/database/stores/forumstore"
	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/valyala/fasthttp"
)

var responseErrorForumNotExists = e.MakeError("forum not exists")

// Список ветвей форума
func Threads(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")
	frmSlug := ctx.UserValue("slug")

	ts := threadstore.PoolThreadSlice.Acquire()
	defer threadstore.PoolThreadSlice.Utilize(ts)

	err := ts.Get(
		frmSlug,
		ctx.QueryArgs().Peek("limit"),
		ctx.QueryArgs().Peek("since"),
		bytes.Equal(ctx.QueryArgs().Peek("desc"), []byte("true")),
	)
	if err != nil {
		panic(err)
	}

	if len(ts) == 0 {
		frm := forumstore.Pool.Acquire()
		defer forumstore.Pool.Utilize(frm)

		err := frm.Get(frmSlug)
		if err != nil {
			ctx.SetStatusCode(404)
			ctx.Write(responseErrorForumNotExists)
			return
		}
	}

	ctx.Write(ts.MustMarshalJSON())
}
