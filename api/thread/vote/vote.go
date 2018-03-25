package vote

import (
	"github.com/VictoriaOtm/forum-api/helpers"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/VictoriaOtm/forum-api/models/thread"
	"github.com/valyala/fasthttp"
	"log"
)

// Проголосовать за веть обсуждения

func Vote(ctx *fasthttp.RequestCtx) {
	thr := model_pool.ThreadPool.Get().(*thread.Thread)
	defer func() {
		thr.Votes = 0
		model_pool.ThreadPool.Put(thr)
	}()

	var err error
	slugOrID := ctx.UserValue("slug_or_id").(string)
	if helpers.IsNumber(slugOrID) {
		err = thr.GetById(slugOrID)
	} else {
		err = thr.GetBySlug(slugOrID)
	}

	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Response.Header.Set("Content-type", "application/json")
		ctx.Write(error_m.CommonError)
		return
	}

	v := thread.Vote{}
	err = v.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		log.Fatalln(err)
	}

	err = thr.PutVote(v)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Response.Header.Set("Content-type", "application/json")
		ctx.Write(error_m.CommonError)
		return
	}

	resp, err := thr.MarshalJSON()
	if err != nil {
		log.Fatalln(err)
	}

	ctx.Write(resp)
	ctx.Response.Header.Set("Content-type", "application/json")
}
