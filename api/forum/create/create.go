package create

import (
	"log"

	"github.com/VictoriaOtm/forum-api/database"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/VictoriaOtm/forum-api/models/forum"
	"github.com/valyala/fasthttp"
)

// Создание форума

func Create(ctx *fasthttp.RequestCtx) {
	frm := model_pool.ForumPool.Get().(*forum.Forum)
	frm.Posts = 0
	frm.Threads = 0
	defer model_pool.ForumPool.Put(frm)

	err := frm.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		log.Fatalln(err)
	}

	err = frm.Create()

	var resp []byte
	switch err {
	case nil:
		resp, err = frm.MarshalJSON()
		ctx.SetStatusCode(201)

	case database.ErrorUserNotExists:
		resp = error_m.CommonError
		ctx.SetStatusCode(404)

	case database.ErrorForumConflict:
		resp, err = frm.MarshalJSON()
		ctx.SetStatusCode(409)

	default:
		log.Fatalln(err)
	}

	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
