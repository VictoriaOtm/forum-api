package get_details

import (
	"log"

	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/VictoriaOtm/forum-api/models/thread"
	"github.com/valyala/fasthttp"
	"github.com/VictoriaOtm/forum-api/helpers"
)

// Получение информации о ветке

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

	var resp []byte
	switch err {
	case nil:
		resp, err = thr.MarshalJSON()
		if err != nil {
			log.Fatalln(err)
		}

	default:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)
	}

	thr.Slug = nil
	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
