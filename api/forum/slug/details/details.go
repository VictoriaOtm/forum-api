package details

import (
	"log"

	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/forum"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/valyala/fasthttp"
)

// Получение информации о форуме

func Details(ctx *fasthttp.RequestCtx) {
	frm := model_pool.ForumPool.Get().(*forum.Forum)
	frm.Posts = 0
	frm.Threads = 0
	defer model_pool.ForumPool.Put(frm)

	err := frm.Get(ctx.UserValue("slug").(string))

	var resp []byte
	switch err {
	case nil:
		resp, err = frm.MarshalJSON()
		if err != nil {
			log.Fatalln(err)
		}
	default:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)
	}

	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
