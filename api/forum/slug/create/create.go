package create

import (
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/thread"
	"log"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/VictoriaOtm/forum-api/database"
)

// создание ветки

func Create(ctx *fasthttp.RequestCtx) {
	thr := model_pool.ThreadPool.Get().(*thread.Thread)
	defer model_pool.ThreadPool.Put(thr)

	err := thr.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		log.Fatalln(err)
	}
	thr.Forum = ctx.UserValue("slug").(string)

	err = thr.Create()

	var resp []byte
	switch err {
	case nil:
		resp, err = thr.MarshalJSON()
		ctx.SetStatusCode(201)

	case database.ErrorUserNotExists:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)

	case database.ErrorForumNotExists:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)

	case database.ErrorThreadConflict:
		resp, err = thr.MarshalJSON()
		ctx.SetStatusCode(409)

	default:
		log.Fatalln(err)
	}

	thr.Slug = nil
	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
