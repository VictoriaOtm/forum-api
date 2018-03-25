package threads

import (
	"log"

	"github.com/VictoriaOtm/forum-api/database"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/VictoriaOtm/forum-api/models/thread"
	"github.com/valyala/fasthttp"
)

// Список ветвей форума

func Threads(ctx *fasthttp.RequestCtx) {
	tArr := model_pool.ThreadArrPool.Get().(thread.Arr)
	defer model_pool.ThreadArrPool.Put(tArr)

	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")

	err := tArr.Get(ctx.UserValue("slug").(string), limit, since, desc)

	var resp []byte
	switch err {
	case nil:
		resp, err = tArr.MarshalJSON()
		if err != nil {
			log.Fatalln(err)
		}

	case database.ErrorForumNotExists:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)

	default:
		log.Fatalln(err)
	}

	tArr = tArr[:0]
	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
