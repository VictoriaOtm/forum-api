package vote

import (
	"log"

	"github.com/VictoriaOtm/forum-api/database/stores/threadstore"
	e "github.com/VictoriaOtm/forum-api/helpers/error"
	"github.com/valyala/fasthttp"
)

var errorThreadNotFound = e.MakeError("error: not found")

// Vote - Проголосовать за веть обсуждения
func Vote(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-type", "application/json")

	thr := threadstore.Pool.Acquire()
	defer threadstore.Pool.Utilize(thr)

	v := threadstore.Vote{}
	if err := v.UnmarshalJSON(ctx.Request.Body()); err != nil {
		log.Fatalln(err)
	}

	if err := thr.PutVote(ctx.UserValue("slug_or_id").(string), v); err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(errorThreadNotFound)
		return
	}

	ctx.Write(thr.MustMarshalJSON())
}
