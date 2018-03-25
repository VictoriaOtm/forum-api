package post_details

import (
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/thread"
	"github.com/VictoriaOtm/forum-api/helpers"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"log"
)

// Обновление ветки

func Details(ctx *fasthttp.RequestCtx) {
	thr := model_pool.ThreadPool.Get().(*thread.Thread)
	defer func() {
		thr.Votes = 0
		model_pool.ThreadPool.Put(thr)
	}()

	var err error
	slugOrId := ctx.UserValue("slug_or_id").(string)
	if helpers.IsNumber(slugOrId) {
		err = thr.GetById(slugOrId)
	} else {
		err = thr.GetBySlug(slugOrId)
	}

	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Response.Header.Set("Content-type", "application/json")
		ctx.Write(error_m.CommonError)
		return
	}

	upd := thread.Update{}
	err = upd.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		log.Fatalln(err)
	}

	err = thr.Update(upd)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := thr.MarshalJSON()
	if err != nil {
		log.Fatalln(err)
	}

	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
