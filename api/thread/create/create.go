package create

import (
	"log"

	"github.com/VictoriaOtm/forum-api/database"
	"github.com/VictoriaOtm/forum-api/model_pool"
	"github.com/VictoriaOtm/forum-api/models/error_m"
	"github.com/VictoriaOtm/forum-api/models/post"
	"github.com/valyala/fasthttp"
)

func Create(ctx *fasthttp.RequestCtx) {
	postArr := model_pool.PostArrPool.Get().(post.PostsArr)
	defer model_pool.PostArrPool.Put(postArr)

	err := postArr.UnmarshalJSON(ctx.Request.Body())
	if err != nil {
		log.Fatalln(err)
	}

	err = postArr.Create(ctx.UserValue("slug_or_id").(string))

	var resp []byte
	switch err {
	case nil:
		resp, err = postArr.MarshalJSON()
		if err != nil {
			log.Fatalln(err)
		}
		ctx.SetStatusCode(201)

	case database.ErrorPostParentConflict:
		ctx.SetStatusCode(409)
		resp = error_m.CommonError

	case database.ErrorUserNotExists:
		ctx.SetStatusCode(404)
		resp = error_m.CommonError

	case database.ErrorThreadNotExists:
		ctx.SetStatusCode(404)
		resp = error_m.CommonError

	default:
		log.Fatalln(err)
	}

	postArr = postArr[:0]
	ctx.Response.Header.Set("Content-type", "application/json")
	ctx.Write(resp)
}
